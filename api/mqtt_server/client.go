package main

import (
	mqtt "api/mqtt"
	types "api/types"
	"encoding/json"
	"fmt"

	paho "github.com/eclipse/paho.mqtt.golang"
)

type ServerState struct {
	ServerIP string
	Mqtt     *mqtt.MQTT
}

func main() {
	// Cria o cliente MQTT
	mqttClient, err := mqtt.NewMQTTClient(types.PORT, types.BROKER)
	if err != nil {
		fmt.Println("Error creating MQTT client:", err)
		return
	}

	serverIP := ""
	fmt.Println("Insira o IP do server/empresa a qual esse carro pertence:")
	fmt.Scanln(&serverIP)

	serverState := ServerState{
		ServerIP: serverIP,
		Mqtt:     mqttClient,
	}

	// Inscrição no tópico de nascimento do carro
	topic := types.CarBirthTopic(serverState.ServerIP)
	serverState.Mqtt.Subscribe(topic, func(client paho.Client, msg paho.Message) {
		// Funcao de callback
		// adiciona o carro na database e se inscreve no tópico de consulta e reserva de rotas

		mqttMessage := &types.MQTT_Message{}
		json.Unmarshal(msg.Payload(), mqttMessage)

		car := &types.Car{}
		json.Unmarshal(mqttMessage.Message, car)

		// TODO adicionar o carro na database

		// Inscrição no tópico de consulta de rotas
		topic = types.CarConsultTopic(serverIP, car.GetCarID())
		serverState.Mqtt.Subscribe(topic, func(client paho.Client, msg paho.Message) {
			// Funcao de callback
			// Deve retornar uma mensagem com payload ListRoutes
			mqttMessage := &types.MQTT_Message{}
			json.Unmarshal(msg.Payload(), mqttMessage)

			routesMessage := &types.RoutesMessage{}
			json.Unmarshal(mqttMessage.Message, routesMessage)

			// city1, city2 := routesMessage.City1, routesMessage.City2
			// TODO requisitar as rotas pela API e retornar na variavel routesList
			routesList := types.RoutesList{
				Routes: []types.Route{},
			}

			payload, _ := json.Marshal(routesList)

			mqttMessage = &types.MQTT_Message{
				Topic:   topic,
				Message: payload,
			}
			serverState.Mqtt.Publish(*mqttMessage)
		})

		// Inscrição no tópico de reserva de rotas
		topic = types.CarReserveTopic(serverIP, car.GetCarID())
		serverState.Mqtt.Subscribe(topic, func(client paho.Client, msg paho.Message) {
			// Funcao de callback
			// Deve retornar uma mensagem com payload ListRoutes
			mqttMessage := &types.MQTT_Message{}
			json.Unmarshal(msg.Payload(), mqttMessage)

			routesMessage := &types.RoutesMessage{}
			json.Unmarshal(mqttMessage.Message, routesMessage)

			// city1, city2 := routesMessage.City1, routesMessage.City2
			// TODO requisitar as rotas pela API e retornar na variavel routesList
			routesList := types.RoutesList{
				Routes: []types.Route{},
			}

			payload, _ := json.Marshal(routesList)

			mqttMessage = &types.MQTT_Message{
				Topic:   topic,
				Message: payload,
			}
			serverState.Mqtt.Publish(*mqttMessage)
		})

		topic = types.CarSelectRouteTopic(serverIP, car.GetCarID())
		serverState.Mqtt.Subscribe(topic, func(client paho.Client, msg paho.Message) {
			// Funcao de callback
			// Deve retornar uma mensagem com payload ListRoutes
			mqttMessage := &types.MQTT_Message{}
			json.Unmarshal(msg.Payload(), mqttMessage)

			route := &types.Route{}
			json.Unmarshal(mqttMessage.Message, route)

			// TODO reservar a rota pela API
		})

		topic = types.CarDeathTopic(serverIP)
		serverState.Mqtt.Subscribe(topic, func(client paho.Client, msg paho.Message) {
			// Funcao de callback
			// Retira o carro da database
			mqttMessage := &types.MQTT_Message{}
			json.Unmarshal(msg.Payload(), mqttMessage)

			car := &types.Car{}
			json.Unmarshal(mqttMessage.Message, car)

			serverState.Mqtt.Client.Unsubscribe(
				types.CarConsultTopic(serverIP, car.GetCarID()),
				types.CarReserveTopic(serverIP, car.GetCarID()),
				types.CarSelectRouteTopic(serverIP, car.GetCarID()),
				types.CarDeathTopic(serverIP),
			)

			// TODO retirar o carro da database
		})
	})

	// Inscrição no tópico de nascimento de um posto
	topic = types.StationBirthTopic(serverIP)
	serverState.Mqtt.Subscribe(topic, func(client paho.Client, msg paho.Message) {
		// Funcao de callback
		// Adiciona o posto na database

	})

	// Inscrição no tópico de nascimento de um posto
	topic = types.StationDeathTopic(serverIP)
	serverState.Mqtt.Subscribe(topic, func(client paho.Client, msg paho.Message) {
		// Funcao de callback
		// Retira o posto da database
	})

	// Mantem o cliente MQTT ativo até o usuário encerrar
	fmt.Println("Enter para encerra o server")
	fmt.Scanln()
}
