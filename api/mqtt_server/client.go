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
		// adiciona o carro na database e se inscreve no tópico de consulta e reserva de rotas

		mqttMessage := &model.MQTT_Message{}
		json.Unmarshal(msg.Payload(), mqttMessage)

		car := &model.Car{}
		json.Unmarshal(mqttMessage.Message, car)

		// TODO adicionar o carro (car) na database

		// Inscrição no tópico de consulta de rotas
		topic = model.CarConsultTopic(serverState.ServerIP, car.GetCarID())
		serverState.Mqtt.Subscribe(topic, func(client paho.Client, msg paho.Message) {
			// Deserializa a mensagem MQTT
			mqttMessage := &model.MQTT_Message{}
			if err := json.Unmarshal(msg.Payload(), mqttMessage); err != nil {
				log.Printf("Erro ao decodificar MQTT_Message: %v", err)
			}
			// Deserializa o contéudo da mensagemMQTT
			routesMessage := &model.RoutesMessage{}
			if err := json.Unmarshal(mqttMessage.Message, routesMessage); err != nil {
				log.Printf("Erro ao decodificar RoutesMessage: %v", err)
			}

			// Extrai as cidades de origem e destino
			city1, city2 := routesMessage.City1, routesMessage.City2
			// Monta a URL para a requisição HTTP das rotas entre as cidades
			url := fmt.Sprintf("http://%s:%s/routes?start_city=%s&end_city=%s", serverState.ServerIP, serverState.Port, city1, city2)
			body := SendHttpGetRequest(url)

			// recebe a lista de rotas
			routesList := []model.Route{}
			json.Unmarshal(body, &routesList)

			// Mapa para armazenar as estações disponíveis de cada empresa
			availableStations := make(map[string][]model.Station)
			// para cada rota na lista de rotas
			for _, route := range routesList {
				// para cada empresa na rota
				for _, company := range route.Waypoints {
					// Monta a URL para obter as informações do servidor da empresa
					url = fmt.Sprintf("http://%s:%s/servers/%s", serverState.ServerIP, serverState.Port, company)
					body = SendHttpGetRequest(url)

					// Recebe as informações do servidor
					server := &model.Server{}
					json.Unmarshal(body, server)

					// Monta a URL para obter as estações disponíveis do servidor
					url = fmt.Sprintf("http://%s:%s/stations", server.ServerIP, server.ServerPort)
					body = SendHttpGetRequest(url)
					// Recebe a lista de estações
					stationList := []model.Station{}
					json.Unmarshal(body, &stationList)
					// Adiciona a lista de estações ao mapa de estações disponíveis
					availableStations[company] = stationList
				}
			}
			// Serializa o mapa de estações disponíveis para JSON
			body, err := json.Marshal(availableStations)
			if err != nil {
				log.Printf("Erro ao serializar availableStations: %v", err)
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
			// O procedimento é o mesmo do tópico de consulta de rotas, esse tópico serve para marcar um fluxo diferente
			mqttMessage := &model.MQTT_Message{}
			if err := json.Unmarshal(msg.Payload(), mqttMessage); err != nil {
				log.Printf("Erro ao decodificar MQTT_Message: %v", err)
			}

			routesMessage := &model.RoutesMessage{}
			if err := json.Unmarshal(mqttMessage.Message, routesMessage); err != nil {
				log.Printf("Erro ao decodificar RoutesMessage: %v", err)
			}

			city1, city2 := routesMessage.City1, routesMessage.City2
			url := fmt.Sprintf("http://%s:%s/routes?start_city=%s&end_city=%s", serverState.ServerIP, serverState.Port, city1, city2)

			// Realiza a requisição HTTP
			body := SendHttpGetRequest(url)
			if body == nil {
			}

			availableStations := make(map[string][]model.Station)

			routesList := []model.Route{}
			json.Unmarshal(body, &routesList)
			for _, route := range routesList {
				for _, company := range route.Waypoints {

					url = fmt.Sprintf("http://%s:%s/servers/%s", serverState.ServerIP, serverState.Port, company)
					body = SendHttpGetRequest(url)

					server := &model.Server{}
					json.Unmarshal(body, server)

					url = fmt.Sprintf("http://%s:%s/stations", server.ServerIP, server.ServerPort)
					body = SendHttpGetRequest(url)
					stationList := []model.Station{}
					json.Unmarshal(body, &stationList)
					availableStations[company] = stationList
				}
			}

			body, err := json.Marshal(availableStations)
			if err != nil {
				log.Printf("Erro ao serializar availableStations: %v", err)
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
			/*
				Recebe uma mensagem do topico "CarSelectRouteTopic", com um payload
				de uma estrutura SelectRouteMessage, com essas informações, uma mensagem
				para reserva no posto desejado é enviada ao broker, que será respondida pelo
				posto requisitado
			*/
			mqttMessage := &model.MQTT_Message{}
			json.Unmarshal(msg.Payload(), mqttMessage)

			selectedRouteMessage := model.SelectRouteMessage{}
			json.Unmarshal(mqttMessage.Message, &selectedRouteMessage)

			// Extrai as informações do carro e da rota selecionada da mensagem
			car := selectedRouteMessage.Car
			selectedStations := selectedRouteMessage.StationsList

			// Fase 1 do 2PC
			// Envia a requisição de preparação para cada station
			url := ""
			allPrepared := true
			for _, station := range selectedStations {
				// Caso o posto seja de outro servidor, a requisição deve ser enviada para o servidor
				if serverState.ServerIP != station.ServerIP {
					url = fmt.Sprintf("http://%s:%s/server/%s/stations/%d",
						serverState.ServerIP, serverState.Port, station.Company, station.StationID)
					// Caso contrário, a requisição deve ser enviada para o próprio servidor
				} else {
					url = fmt.Sprintf("http://%s:%s/stations/%d",
						serverState.ServerIP, serverState.Port, station.StationID)
				}
				// Envia a requisição de preparação
				prepared, err := SendPrepareRequest(url, car.CarID)
				// Verifica se a requisição foi bem sucedida
				if err != nil || !prepared {
					allPrepared = false
					break
				}
			}

			// Se alguma estação não estiver preparada, envia a requisição de abortar
			if !allPrepared {
				for _, station := range selectedStations {
					// caso o posto seja de outro servidor
					if serverState.ServerIP != station.ServerIP {
						url = fmt.Sprintf("http://%s:%s/server/%s/stations/%d",
							serverState.ServerIP, serverState.Port, station.Company, station.StationID)
					} else {
						url = fmt.Sprintf("http://%s:%s/stations/%d",
							serverState.ServerIP, serverState.Port, station.StationID)
					}

					SendAbortRequest(url, car)
				}
				fmt.Printf("Prepare failed")
			}

			// Phase 2: Commit
			if allPrepared {
				// Envia a requisição de commit para cada estação
				for _, station := range selectedStations {
					// caso o posto seja de outro servidor
					if serverState.ServerIP != station.ServerIP {
						url = fmt.Sprintf("http://%s:%s/server/%s/stations/%d",
							serverState.ServerIP, serverState.Port, station.Company, station.StationID)
					} else {
						url = fmt.Sprintf("http://%s:%s/stations/%d",
							serverState.ServerIP, serverState.Port, station.StationID)
					}
					// Envia a requisição de commit
					if err := SendCommitRequest(url, car.CarID); err != nil {
						fmt.Printf("Commit failed for %d: %v\n", station.StationID, err)
					}
					// Envia resposta para o carro
					topic = model.ResponseStationReserveTopic(station.ServerIP, fmt.Sprintf("%d", station.StationID))

					carInfo := &model.CarInfo{
						CarId: car.GetCarID(),
					}
					payload, _ := json.Marshal(carInfo)

					mqttMessage = &model.MQTT_Message{
						Topic:   topic,
						Message: payload,
					}

					serverState.Mqtt.Publish(*mqttMessage)
				}
			}
		})
		topic = model.FinishRouteTopic(serverState.ServerIP, car.CarID)
		serverState.Mqtt.Subscribe(topic, func(client paho.Client, msg paho.Message) {

			var mqttMsg model.MQTT_Message
			if err := json.Unmarshal(msg.Payload(), &mqttMsg); err != nil {
				log.Printf("Erro ao decodificar MQTT_Message: %v", err)
				return
			}

			var finishMsg model.FinishRouteMessage
			if err := json.Unmarshal(mqttMsg.Message, &finishMsg); err != nil {
				log.Printf("Erro ao decodificar FinishRouteMessage: %v", err)
				return
			}

			// Libera cada estação informada, enviando também o car_id no payload
			for _, station := range finishMsg.StationsList {
				var url string
				// Caso o posto seja de outro servidor
				if serverState.ServerIP != station.ServerIP {
					url = fmt.Sprintf("http://%s:%s/server/%s/stations/%d/release",
						serverState.ServerIP, serverState.Port, station.Company, station.StationID)
				} else {
					url = fmt.Sprintf("http://%s:%s/stations/%d/release",
						serverState.ServerIP, serverState.Port, station.StationID)
				}
				// Serializa o car ID
				payload := struct {
					CarID int `json:"car_id"`
				}{CarID: finishMsg.Car.CarID}
				payloadBytes, _ := json.Marshal(payload)

				// Envia a requisição HTTP para liberar a estação
				_, err := SendHttpPutRequest(url, payloadBytes)
				if err != nil {
					log.Printf("Erro ao liberar estação %d: %v", station.StationID, err)
				} else {
					log.Printf("Estação %d liberada com sucesso!", station.StationID)
				}
			}

			// Envia mensagem de resposta para o carro
			responsePayload := struct {
				Message string `json:"message"`
			}{
				Message: "Estações liberadas com sucesso",
			}
			responseBytes, _ := json.Marshal(responsePayload)
			responseMsg := &model.MQTT_Message{
				Topic:   model.ResponseFinishRouteTopic(serverState.ServerIP, finishMsg.Car.CarID),
				Message: responseBytes,
			}
			if err := serverState.Mqtt.Publish(*responseMsg); err != nil {
				log.Printf("Erro ao publicar mensagem de resposta de finalização de rota: %v", err)
			}
		})
		// Inscrição no tópico de morte do carro
		topic = model.CarDeathTopic(serverState.ServerIP)
		serverState.Mqtt.Subscribe(topic, func(client paho.Client, msg paho.Message) {
			mqttMessage := &model.MQTT_Message{}
			json.Unmarshal(msg.Payload(), mqttMessage)

			car := &model.Car{}
			json.Unmarshal(mqttMessage.Message, car)

			// Desinscreve dos tópicos relacionados ao carro
			serverState.Mqtt.Client.Unsubscribe(
				model.CarConsultTopic(serverState.ServerIP, car.GetCarID()),
				model.CarReserveTopic(serverState.ServerIP, car.GetCarID()),
				model.CarSelectRouteTopic(serverState.ServerIP, car.GetCarID()),
			)
		})
	})

	// Inscrição no tópico de nascimento de um posto
	topic = model.StationBirthTopic(serverState.ServerIP)
	serverState.Mqtt.Subscribe(topic, func(client paho.Client, msg paho.Message) {
		// Adiciona o posto na database
		mqttMessage := &model.MQTT_Message{}
		if err := json.Unmarshal(msg.Payload(), mqttMessage); err != nil {
			log.Printf("Erro ao decodificar MQTT_Message: %v", err)
		}

		station := &model.Station{}
		if err := json.Unmarshal(mqttMessage.Message, station); err != nil {
			log.Printf("Erro ao decodificar Station: %v", err)
		}

		// Altera o campo Company
		station.Company = serverState.ServerCompany

		// Codifica novamente o objeto Station para JSON
		updatedPayload, err := json.Marshal(station)
		if err != nil {
			log.Printf("Erro ao serializar Station: %v", err)
		}

		// Envia a requisição HTTP com o payload atualizado
		url := fmt.Sprintf("http://%s:%s/stations", serverState.ServerIP, serverState.Port)
		SendHttpPostRequest(url, updatedPayload)

		// Inscrição no tópico de nascimento de um posto
		topic = model.StationDeathTopic(serverState.ServerIP)
		serverState.Mqtt.Subscribe(topic, func(client paho.Client, msg paho.Message) {
			// Função de callback para processar mensagens de morte de postos

			// Decodifica a mensagem MQTT
			mqttMessage := &model.MQTT_Message{}
			if err := json.Unmarshal(msg.Payload(), mqttMessage); err != nil {
				log.Printf("Erro ao decodificar MQTT_Message: %v", err)
			}

			// Decodifica o ID da estação a partir do campo Message
			var stationID int
			if err := json.Unmarshal(mqttMessage.Message, &stationID); err != nil {
				log.Printf("Erro ao decodificar StationID: %v", err)
			}

			// Define a URL do endpoint para remover a estação com o ID na URL
			url := fmt.Sprintf("http://%s:%s/stations/%d/remove", serverState.ServerIP, serverState.Port, stationID)

			// Envia a requisição HTTP PUT
			SendHttpPutRequest(url, nil) // Nenhum payload é necessário, pois o ID está na URL
		})

	})

	// Inscrição no tópico de nascimento do servidor, utiliazado para registrar novos servidores
	topic = model.ServerBirthTopic(serverState.ServerIP)
	serverState.Mqtt.Subscribe(topic, func(client paho.Client, msg paho.Message) {
		// Função de callback para processar mensagens de nascimento de servidores
		mqttMessage := &model.MQTT_Message{}
		if err := json.Unmarshal(msg.Payload(), mqttMessage); err != nil {
			log.Printf("Erro ao decodificar MQTT_Message: %v", err)
		}

		var serverInfo struct {
			Company    string `json:"company"`
			ServerIP   string `json:"server_ip"`
			ServerPort string `json:"server_port"`
		}
		if err := json.Unmarshal(mqttMessage.Message, &serverInfo); err != nil {
			log.Printf("Erro ao decodificar informações do servidor: %v", err)
		}

		// Registra ou atualiza o servidor no banco de dados via POST HTTP
		url := fmt.Sprintf("http://%s:%s/servers/register", serverIP, port)
		payload := map[string]string{
			"company":     serverInfo.Company,
			"server_ip":   serverInfo.ServerIP,
			"server_port": serverInfo.ServerPort,
		}

		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			log.Printf("Erro ao serializar payload: %v", err)
		}

		SendHttpPostRequest(url, jsonPayload)

		// Envia uma mensagem de resposta para que o servidor que enviou a mensagem de nascimento
		// atualize sua data base com os servidores já cadastrados
		responsePayload := map[string]string{
			"company":     serverState.ServerCompany,
			"server_ip":   serverState.ServerIP,
			"server_port": serverState.Port,
		}
		responseJson, err := json.Marshal(responsePayload)
		if err != nil {
			log.Printf("Erro ao serializar payload de resposta: %v", err)
		}
		responseMessage := model.MQTT_Message{
			Topic:   model.ResponseServerBirthTopic(serverInfo.ServerIP),
			Message: responseJson,
		}

		if err := serverState.Mqtt.Publish(responseMessage); err != nil {
			log.Printf("Erro ao publicar mensagem de resposta: %v", err)
		}
	})

	// Envia uma mensagem de nascimento do servidor para o broker
	// Alertando os outros servidores que esse servidor está ativo
	topic = model.ServerBirthTopic(serverState.ServerIP)
	// Cria o payload da mensagem
	payload := map[string]string{
		"company":     serverState.ServerCompany, // Substitua pelo nome real da empresa
		"server_ip":   serverState.ServerIP,
		"server_port": serverState.Port,
	}
	// Serializa o payload para JSON
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Erro ao serializar payload: %v", err)
	}
	// Publica a mensagem no tópico
	err = serverState.Mqtt.Publish(model.MQTT_Message{
		Topic:   topic,
		Message: jsonPayload,
	})
	if err != nil {
		log.Printf("Erro ao publicar mensagem de nascimento do servidor: %v", err)
	}

	// Inscrição no tópico de resposta de nascimento do servidor
	// Atualiza a base de dados com os servidores já cadastrados
	topic = model.ResponseServerBirthTopic(serverState.ServerIP)
	serverState.Mqtt.Subscribe(topic, func(client paho.Client, msg paho.Message) {
		// Função de callback para processar mensagens de nascimento de servidores
		mqttMessage := &model.MQTT_Message{}
		if err := json.Unmarshal(msg.Payload(), mqttMessage); err != nil {
			log.Printf("Erro ao decodificar MQTT_Message: %v", err)
		}

		var serverInfo struct {
			Company    string `json:"company"`
			ServerIP   string `json:"server_ip"`
			ServerPort string `json:"server_port"`
		}
		if err := json.Unmarshal(mqttMessage.Message, &serverInfo); err != nil {
			log.Printf("Erro ao decodificar informações do servidor: %v", err)
		}

		// Registra ou atualiza o servidor no banco de dados via POST HTTP
		url := fmt.Sprintf("http://%s:%s/servers/register", serverIP, port)
		payload := map[string]string{
			"company":     serverInfo.Company,
			"server_ip":   serverInfo.ServerIP,
			"server_port": serverInfo.ServerPort,
		}

		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			log.Printf("Erro ao serializar payload: %v", err)
		}

		SendHttpPostRequest(url, jsonPayload)
	})
}

// Função para obter o IP local da máquina
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

// Envia uma requisição HTTP GET e retorna o corpo da resposta
func SendHttpGetRequest(url string) []byte {
	// Realiza a requisição HTTP
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Erro na requisição GET para %s: %v", url, err)
		return nil
	}
	defer resp.Body.Close() // Certifique-se de fechar o corpo da resposta

	// Verifique se o status HTTP é 200 OK
	if resp.StatusCode != http.StatusOK {
		log.Printf("Erro: status de resposta %d para %s", resp.StatusCode, url)
		return nil
	}

	// Lê o corpo da resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Erro ao ler corpo da resposta: %v", err)
		return nil
	}

	return body
}

// Envia uma requisição HTTP PUT e retorna o corpo da resposta
// Se o status da resposta não for 200 OK, retorna um erro
func SendHttpPutRequest(url string, payload []byte) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição PUT para %s: %v", url, err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro na requisição PUT para %s: %v", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler corpo da resposta: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return body, fmt.Errorf("erro: status %d, resposta: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// Envia uma requisição HTTP POST e não espera resposta
func SendHttpPostRequest(url string, payload []byte) {
	// Realiza a requisição HTTP
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Printf("Erro na requisição POST para %s: %v", url, err)
		return
	}
	defer resp.Body.Close() // Certifique-se de fechar o corpo da resposta

	// Verifique se o status HTTP é 200 OK
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Erro: status %d, resposta: %s", resp.StatusCode, string(body))
		return
	}
}

// 2PC: Envia uma requisição HTTP PUT para preparar, commit ou abortar
func SendPrepareRequest(url string, carID int) (bool, error) {
	// Cria o payload com apenas o CarID
	payload := struct {
		CarID int `json:"car_id"`
	}{
		CarID: carID,
	}

	jsonPayload, _ := json.Marshal(payload)
	fmt.Println("Prepare Payload enviado:", string(jsonPayload)) // Imprime o payload enviado

	req, err := http.NewRequest(http.MethodPut, url+"/prepare", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK, nil
}

// 2PC: Envia uma requisição HTTP PUT para confirmar o commit
func SendCommitRequest(url string, carID int) error {
	// Cria o payload com apenas o CarID
	payload := struct {
		CarID int `json:"car_id"`
	}{
		CarID: carID,
	}

	jsonPayload, _ := json.Marshal(payload)
	fmt.Println("Commit Payload enviado:", string(jsonPayload)) // Imprime o payload enviado

	req, err := http.NewRequest(http.MethodPut, url+"/commit", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("commit failed for %s", url)
	}

	return nil
}

// 2PC: Envia uma requisição HTTP PUT para abortar
func SendAbortRequest(url string, payload interface{}) error {
	jsonPayload, _ := json.Marshal(payload)
	req, err := http.NewRequest(http.MethodPut, url+"/abort", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("abort failed for %s", url)
	}
	return nil
}
