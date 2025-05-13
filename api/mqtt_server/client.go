package mqtt_server

import (
	model "api/model"
	mqtt "api/mqtt"
	types "api/types"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
)

type ServerState struct {
	ServerIP      string
	Mqtt          *mqtt.MQTT
	Port          string
	ServerCompany string
}

func MqttMain(serverCompany string, port string) {
	// Aguarda iniciação da API
	time.Sleep(2 * time.Second)
	// Cria o cliente MQTT
	mqttClient, err := mqtt.NewMQTTClient(types.PORT, types.BROKER)
	if err != nil {
		fmt.Println("Error creating MQTT client:", err)
		return
	}

	serverIP := ""
	// fmt.Println("Insira o IP do server:")
	// fmt.Scanln(&serverIP)
	serverIP, _ = GetLocalIP()
	fmt.Printf("IP local detectado: %s\n", serverIP)

	serverState := ServerState{
		ServerIP:      serverIP,
		Mqtt:          mqttClient,
		Port:          port,
		ServerCompany: serverCompany,
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
		topic = model.CarConsultTopic(serverState.ServerIP, car.GetCarID())
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
			url := fmt.Sprintf("http://%s:%s/routes?start_city=%s&end_city=%s", serverState.ServerIP, serverState.Port, city1, city2)

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

			// Cria a mensagem MQTT de resposta
			mqttMessage = &model.MQTT_Message{
				Topic:   model.ResponseCarConsultTopic(serverState.ServerIP, car.GetCarID()),
				Message: body,
			}

			// Publica a mensagem MQTT de volta
			if err := serverState.Mqtt.Publish(*mqttMessage); err != nil {
				log.Printf("Erro ao publicar a mensagem MQTT: %v", err)
			}
		})

		// Inscrição no tópico de reserva de rotas
		topic = model.CarReserveTopic(serverState.ServerIP, car.GetCarID())
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
			url := fmt.Sprintf("http://%s:%s/routes?start_city=%s&end_city=%s", serverState.ServerIP, serverState.Port, city1, city2)

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

			// Cria a mensagem MQTT de resposta
			mqttMessage = &model.MQTT_Message{
				Topic:   model.ResponseCarReserveTopic(serverState.ServerIP, car.GetCarID()),
				Message: body,
			}

			// Publica a mensagem MQTT de volta
			if err := serverState.Mqtt.Publish(*mqttMessage); err != nil {
				log.Printf("Erro ao publicar a mensagem MQTT: %v", err)
			}
		})

		topic = model.CarSelectRouteTopic(serverState.ServerIP, car.GetCarID())
		serverState.Mqtt.Subscribe(topic, func(client paho.Client, msg paho.Message) {
			// Funcao de callback
			// Recebe uma mensagem do topico 
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

		topic = model.CarDeathTopic(serverState.ServerIP)
		serverState.Mqtt.Subscribe(topic, func(client paho.Client, msg paho.Message) {
			// Funcao de callback
			// Retira o carro da database
			mqttMessage := &model.MQTT_Message{}
			json.Unmarshal(msg.Payload(), mqttMessage)

			car := &model.Car{}
			json.Unmarshal(mqttMessage.Message, car)

			serverState.Mqtt.Client.Unsubscribe(
				model.CarConsultTopic(serverState.ServerIP, car.GetCarID()),
				model.CarReserveTopic(serverState.ServerIP, car.GetCarID()),
				model.CarSelectRouteTopic(serverState.ServerIP, car.GetCarID()),
			)

			// TODO retirar o carro da database
		})
	})

	// Inscrição no tópico de nascimento de um posto
	topic = model.StationBirthTopic(serverState.ServerIP)
	serverState.Mqtt.Subscribe(topic, func(client paho.Client, msg paho.Message) {
		// Funcao de callback
		// Adiciona o posto na database
		mqttMessage := &model.MQTT_Message{}
		json.Unmarshal(msg.Payload(), mqttMessage)
		station := &model.Station{}
		json.Unmarshal(mqttMessage.Message, station)

		// TODO adicionar o posto (station) na database

		// Inscrição no tópico de nascimento de um posto
		topic = model.StationDeathTopic(serverState.ServerIP)
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

	topic = model.ServerBirthTopic(serverState.ServerIP)
	serverState.Mqtt.Subscribe(topic, func(client paho.Client, msg paho.Message) {
		// Função de callback para processar mensagens de nascimento de servidores
		log.Printf("Mensagem recebida no tópico %s: %s", msg.Topic(), string(msg.Payload()))
		mqttMessage := &model.MQTT_Message{}
		if err := json.Unmarshal(msg.Payload(), mqttMessage); err != nil {
			log.Printf("Erro ao decodificar MQTT_Message: %v", err)
			return
		}

		var serverInfo struct {
			Company  string `json:"company"`
			ServerIP string `json:"server_ip"`
		}
		if err := json.Unmarshal(mqttMessage.Message, &serverInfo); err != nil {
			log.Printf("Erro ao decodificar informações do servidor: %v", err)
			return
		}

		// Registra ou atualiza o servidor no banco de dados via POST HTTP
		url := fmt.Sprintf("http://%s:%s/servers/register", serverIP, port)
		payload := map[string]string{
			"company":   serverInfo.Company,
			"server_ip": serverInfo.ServerIP,
		}

		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			log.Printf("Erro ao serializar payload: %v", err)
			return
		}

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
		if err != nil {
			log.Printf("Erro ao fazer POST HTTP para %s: %v", url, err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			log.Printf("Erro ao registrar servidor: status %d, resposta: %s", resp.StatusCode, string(body))
			return
		}

		responsePayload := map[string]string{
			"company":   serverState.ServerCompany,
			"server_ip": serverState.ServerIP,
		}
		responseJson, err := json.Marshal(responsePayload)
		if err != nil {
			log.Printf("Erro ao serializar payload de resposta: %v", err)
			return
		}
		responseMessage := model.MQTT_Message{
			Topic:   model.ResponseServerBirthTopic(serverInfo.ServerIP),
			Message: responseJson,
		}

		if err := serverState.Mqtt.Publish(responseMessage); err != nil {
			log.Printf("Erro ao publicar mensagem de resposta: %v", err)
			return
		}
	})

	topic = model.ServerBirthTopic(serverState.ServerIP)
	// Cria o payload da mensagem
	payload := map[string]string{
		"company":   serverState.ServerCompany, // Substitua pelo nome real da empresa
		"server_ip": serverState.ServerIP,
	}

	// Serializa o payload para JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Erro ao serializar payload: %v", err)
		return
	}

	// Publica a mensagem no tópico
	err = serverState.Mqtt.Publish(model.MQTT_Message{
		Topic:   topic,
		Message: jsonPayload,
	})
	if err != nil {
		log.Printf("Erro ao publicar mensagem de nascimento do servidor: %v", err)
		return
	}

	topic = model.ResponseServerBirthTopic(serverState.ServerIP)
	serverState.Mqtt.Subscribe(topic, func(client paho.Client, msg paho.Message) {
		// Função de callback para processar mensagens de nascimento de servidores
		log.Printf("Mensagem recebida no tópico %s: %s", msg.Topic(), string(msg.Payload()))
		mqttMessage := &model.MQTT_Message{}
		if err := json.Unmarshal(msg.Payload(), mqttMessage); err != nil {
			log.Printf("Erro ao decodificar MQTT_Message: %v", err)
			return
		}

		var serverInfo struct {
			Company  string `json:"company"`
			ServerIP string `json:"server_ip"`
		}
		if err := json.Unmarshal(mqttMessage.Message, &serverInfo); err != nil {
			log.Printf("Erro ao decodificar informações do servidor: %v", err)
			return
		}

		// Registra ou atualiza o servidor no banco de dados via POST HTTP
		url := fmt.Sprintf("http://%s:%s/servers/register", serverIP, port)
		payload := map[string]string{
			"company":   serverInfo.Company,
			"server_ip": serverInfo.ServerIP,
		}

		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			log.Printf("Erro ao serializar payload: %v", err)
			return
		}

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
		if err != nil {
			log.Printf("Erro ao fazer POST HTTP para %s: %v", url, err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			log.Printf("Erro ao registrar servidor: status %d, resposta: %s", resp.StatusCode, string(body))
			return
		}
	})
}

func GetLocalIP() (string, error) {
	// Obtém todas as interfaces de rede do sistema
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("erro ao obter interfaces de rede: %v", err)
	}

	for _, iface := range interfaces {
		// Ignora interfaces que estão desativadas ou não suportam multicast
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		// Obtém os endereços associados à interface
		addrs, err := iface.Addrs()
		if err != nil {
			return "", fmt.Errorf("erro ao obter endereços da interface %s: %v", iface.Name, err)
		}

		for _, addr := range addrs {
			// Verifica se o endereço é do tipo IP
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			// Ignora endereços IPv6 ou endereços de loopback
			if ip == nil || ip.IsLoopback() || ip.To4() == nil {
				continue
			}

			// Retorna o primeiro endereço IPv4 encontrado
			return ip.String(), nil
		}
	}

	return "", fmt.Errorf("nenhum endereço IP válido encontrado")
}
