package mqtt_server

import (
	model "api/model"
	mqtt "api/mqtt"
	types "api/types"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	paho "github.com/eclipse/paho.mqtt.golang"
)

type ServerState struct {
	ServerIP string
	Mqtt     *mqtt.MQTT
}

func MqttMain() {
	// Cria o cliente MQTT
	mqttClient, err := mqtt.NewMQTTClient(types.PORT, types.BROKER)
	if err != nil {
		fmt.Println("Error creating MQTT client:", err)
		return
	}

	serverIP := ""
	fmt.Println("Insira o IP do server/empresa:")
	fmt.Scanln(&serverIP)

	serverState := ServerState{
		ServerIP: serverIP,
		Mqtt:     mqttClient,
	}

	// Inscrição no tópico de nascimento do carro
	topic := model.CarBirthTopic(serverState.ServerIP)
	serverState.Mqtt.Subscribe(topic, func(client paho.Client, msg paho.Message) {
		// Funcao de callback
		// adiciona o carro na database e se inscreve no tópico de consulta e reserva de rotas

		mqttMessage := &model.MQTT_Message{}
		json.Unmarshal(msg.Payload(), mqttMessage)

		car := &model.Car{}
		json.Unmarshal(mqttMessage.Message, car)

		// TODO adicionar o carro (car) na database

		// Inscrição no tópico de consulta de rotas
		topic = model.CarConsultTopic(serverIP, car.GetCarID())
		serverState.Mqtt.Subscribe(topic, func(client paho.Client, msg paho.Message) {
			// Função de callback
			mqttMessage := &model.MQTT_Message{}
			if err := json.Unmarshal(msg.Payload(), mqttMessage); err != nil {
				log.Printf("Erro ao decodificar MQTT_Message: %v", err)
				return
			}

			routesMessage := &model.RoutesMessage{}
			if err := json.Unmarshal(mqttMessage.Message, routesMessage); err != nil {
				log.Printf("Erro ao decodificar RoutesMessage: %v", err)
				return
			}

			city1, city2 := routesMessage.City1, routesMessage.City2
			url := "http://172.16.103.10:8081/routes?start_city=" + city1 + "&end_city=" + city2

			// Realiza a requisição HTTP
			resp, err := http.Get(url)
			if err != nil {
				log.Printf("Erro na requisição GET para %s: %v", url, err)
				return
			}
			defer resp.Body.Close() // Certifique-se de fechar o corpo da resposta

			// Verifique se o status HTTP é 200 OK
			if resp.StatusCode != http.StatusOK {
				log.Printf("Erro: status de resposta %d para %s", resp.StatusCode, url)
				return
			}

			// Lê o corpo da resposta
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("Erro ao ler corpo da resposta: %v", err)
				return
			}

			// Deserializa a resposta no formato esperado
			var unmarshal []model.Route
			if err := json.Unmarshal(body, &unmarshal); err != nil {
				log.Printf("Erro ao deserializar o corpo da resposta: %v", err)
				return
			}

			// Imprime a resposta (útil para debugging)
			fmt.Println(unmarshal)

			// Cria a mensagem MQTT de resposta
			mqttMessage = &model.MQTT_Message{
				Topic:   model.ResponseCarConsultTopic(serverIP, car.GetCarID()),
				Message: body,
			}

			// Publica a mensagem MQTT de volta
			if err := serverState.Mqtt.Publish(*mqttMessage); err != nil {
				log.Printf("Erro ao publicar a mensagem MQTT: %v", err)
			}
		})

		// Inscrição no tópico de reserva de rotas
		topic = model.CarReserveTopic(serverIP, car.GetCarID())
		serverState.Mqtt.Subscribe(topic, func(client paho.Client, msg paho.Message) {
			// Funcao de callback
			// Deve retornar uma mensagem com payload ListRoutes
			mqttMessage := &model.MQTT_Message{}
			json.Unmarshal(msg.Payload(), mqttMessage)

			routesMessage := &model.RoutesMessage{}
			json.Unmarshal(mqttMessage.Message, routesMessage)

			// city1, city2 := routesMessage.City1, routesMessage.City2
			// TODO requisitar as rotas pela API e retornar na variavel routesList

			routesList := model.RoutesList{
				Routes: []model.Route{},
			}

			payload, _ := json.Marshal(routesList)

			mqttMessage = &model.MQTT_Message{
				Topic:   model.CarReserveTopic(serverIP, car.GetCarID()),
				Message: payload,
			}
			serverState.Mqtt.Publish(*mqttMessage)
		})

		topic = model.CarSelectRouteTopic(serverIP, car.GetCarID())
		serverState.Mqtt.Subscribe(topic, func(client paho.Client, msg paho.Message) {
			// Funcao de callback
			// Deve retornar uma mensagem com payload ListRoutes
			mqttMessage := &model.MQTT_Message{}
			json.Unmarshal(msg.Payload(), mqttMessage)

			selectRouteMessage := &model.SelectRouteMessage{}
			json.Unmarshal(mqttMessage.Message, selectRouteMessage)
			// car := selectRouteMessage.Car
			// route := selectRouteMessage.Route
			// for _, waypoint := range route.Waypoints {
			// 	topic = model.StationReserveTopic(serverIP, waypoint)

			// 	carInfo := &model.CarInfo{
			// 		CarId: car.GetCarID(),
			// 	}
			// 	payload, _ := json.Marshal(carInfo)

			// 	mqttMessage = &model.MQTT_Message{
			// 		Topic:   topic,
			// 		Message: payload,
			// 	}

			// 	serverState.Mqtt.Publish(*mqttMessage)
			// }

			// TODO reservar a rota pela API
		})

		topic = model.CarDeathTopic(serverIP)
		serverState.Mqtt.Subscribe(topic, func(client paho.Client, msg paho.Message) {
			// Funcao de callback
			// Retira o carro da database
			mqttMessage := &model.MQTT_Message{}
			json.Unmarshal(msg.Payload(), mqttMessage)

			car := &model.Car{}
			json.Unmarshal(mqttMessage.Message, car)

			serverState.Mqtt.Client.Unsubscribe(
				model.CarConsultTopic(serverIP, car.GetCarID()),
				model.CarReserveTopic(serverIP, car.GetCarID()),
				model.CarSelectRouteTopic(serverIP, car.GetCarID()),
			)

			// TODO retirar o carro da database
		})
	})

	// Inscrição no tópico de nascimento de um posto
	topic = model.StationBirthTopic(serverIP)
	serverState.Mqtt.Subscribe(topic, func(client paho.Client, msg paho.Message) {
		// Funcao de callback
		// Adiciona o posto na database
		mqttMessage := &model.MQTT_Message{}
		json.Unmarshal(msg.Payload(), mqttMessage)
		station := &model.Station{}
		json.Unmarshal(mqttMessage.Message, station)

		// TODO adicionar o posto (station) na database

		// Inscrição no tópico de nascimento de um posto
		topic = model.StationDeathTopic(serverIP)
		serverState.Mqtt.Subscribe(topic, func(client paho.Client, msg paho.Message) {
			// Funcao de callback
			// Retira o posto da database

			mqttMessage := &model.MQTT_Message{}
			json.Unmarshal(msg.Payload(), mqttMessage)
			station := &model.Station{}
			json.Unmarshal(mqttMessage.Message, station)

			// TODO retirar o posto (station) da database
		})
	})

	// Mantem o cliente MQTT ativo até o usuário encerrar
	fmt.Println("Enter para encerra o server")
	fmt.Scanln()
}
