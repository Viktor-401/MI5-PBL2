package main

import (
	mqtt "clients/mqtt"
	types "clients/types"
	"encoding/json"
	"fmt"
	"strconv"

	paho "github.com/eclipse/paho.mqtt.golang"
)

type CarState struct {
	Car      types.Car
	ServerIP string
	Mqtt     *mqtt.MQTT
}

func main() {
	waitChan := make(chan bool, 1)
	car := types.GetNewRandomCar()

	// Cria o cliente MQTT
	mqttClient, err := mqtt.NewMQTTClient(types.PORT, types.BROKER)
	if err != nil {
		fmt.Println("Error creating MQTT client:", err)
		return
	}

	serverIP := ""
	fmt.Println("Insira o IP do server/empresa a qual esse carro pertence:")
	fmt.Scanln(&serverIP)

	carState := CarState{
		Car:      car,
		ServerIP: serverIP,
		Mqtt:     mqttClient,
	}

	// Mensagem de nascimento do carro, que informa o servidor que o carro está online
	birthMessage, err := carState.BirthMessage()
	if err != nil {
		fmt.Println("Error creating birth message:", err)
		return
	}
	err = carState.Mqtt.Publish(birthMessage)
	if err != nil {
		fmt.Println("Error publishing birth message:", err)
		return
	}

	// Inscrição no tópico de consulta de rotas
	topic := types.ResponseCarConsultTopic(serverIP, carState.Car.GetCarID())
	carState.Mqtt.Subscribe(topic, func(client paho.Client, msg paho.Message) {
		// Funcao de callback para quando uma mensagem é recebida
		UnmarshalListRoutes(msg)
		waitChan <- true
	})

	// Inscrição no tópico de reserva de rotas
	topic = types.ResponseCarReserveTopic(serverIP, carState.Car.GetCarID())
	carState.Mqtt.Subscribe(topic, func(client paho.Client, msg paho.Message) {
		// Funcao de callback para quando uma mensagem é recebida
		selectedStations, err := UnmarshalListRoutesSelect(msg)
		if err != nil {
			fmt.Println("Error unmarshalling route selection:", err)
			return
		}

		message, err := carState.SelectRouteMessage(carState.Car, selectedStations)
		if err != nil {
			fmt.Println("Error creating reserve route message:", err)
			return
		}
		err = carState.Mqtt.Publish(message)
		if err != nil {
			fmt.Println("Error publishing reserve route message:", err)
			return
		}
		waitChan <- true
	})

	exit := false
	for !exit {
		action := 0
        for {
            fmt.Println(` Escolha uma ação:
        1- Consultar Rotas
        2- Reservar Postos
        `)
            fmt.Scanln(&action)
            if action == 1 || action == 2 {
                break
            }
            fmt.Println("Ação inválida. Tente novamente.")
        }
		if action == 1 {
			// Consultar rotas
			city1, city2 := CityInput()

			// Envia a mensagem de consulta de rotas
			message, err := carState.ConsultRouteMessage(city1, city2)
			if err != nil {
				fmt.Println("Error creating consult route message:", err)
				return
			}
			err = carState.Mqtt.Publish(message)
			if err != nil {
				fmt.Println("Error publishing consult route message:", err)
				return
			}
		} else if action == 2 {
			city1, city2 := CityInput()

			// Envia a mensagem de reserva de rotas
			message, err := carState.ReserveRouteMessage(city1, city2)
			if err != nil {
				fmt.Println("Error creating reserve route message:", err)
				return
			}
			err = carState.Mqtt.Publish(message)
			if err != nil {
				fmt.Println("Error publishing reserve route message:", err)
				return
			}
		} else {
			fmt.Println("Ação inválida. Tente novamente.")
		}
		<-waitChan
	}

	// Mantem o cliente MQTT ativo até o usuário encerrar
	fmt.Println("Enter para encerra o posto")
	fmt.Scanln()
	// Mensagem de morte do posto, que informa o servidor que o posto está offline
	message, err := carState.DeathMessage()
	if err != nil {
		fmt.Println("Error creating death message:", err)
		return
	}
	carState.Mqtt.Publish(message)
}

// Retorna a mensagem de nascimento do carro, que informa o servidor que o carro está online
func (s *CarState) BirthMessage() (types.MQTT_Message, error) {
	topic := types.CarBirthTopic(s.ServerIP)

	payload, err := json.Marshal(s.Car)
	if err != nil {
		return types.MQTT_Message{}, err
	}

	return types.MQTT_Message{
		Topic:   topic,
		Message: payload,
	}, nil
}

// Retorna a mensagem de morte do carro, que informa o servidor que o carro está offline
func (s *CarState) DeathMessage() (types.MQTT_Message, error) {
	topic := types.CarDeathTopic(s.ServerIP)

	payload, err := json.Marshal(s.Car)
	if err != nil {
		return types.MQTT_Message{}, err
	}

	return types.MQTT_Message{
		Topic:   topic,
		Message: payload,
	}, nil
}

// Retorna a mensagem de consulta de rotas para ser enviada ao servidor via MQTT
func (s *CarState) ConsultRouteMessage(city1 string, city2 string) (types.MQTT_Message, error) {
	topic := types.CarConsultTopic(s.ServerIP, s.Car.GetCarID())

	consultRoute := types.RoutesMessage{
		City1: city1,
		City2: city2,
	}

	payload, err := json.Marshal(consultRoute)
	if err != nil {
		return types.MQTT_Message{}, err
	}

	return types.MQTT_Message{
		Topic:   topic,
		Message: payload,
	}, nil
}

// Retorna a mensagem de reserva de rotas para ser enviada ao servidor via MQTT
func (s *CarState) ReserveRouteMessage(city1 string, city2 string) (types.MQTT_Message, error) {
	topic := types.CarReserveTopic(s.ServerIP, s.Car.GetCarID())

	reserveRoute := types.RoutesMessage{
		City1: city1,
		City2: city2,
	}

	payload, err := json.Marshal(reserveRoute)
	if err != nil {
		return types.MQTT_Message{}, err
	}

	return types.MQTT_Message{
		Topic:   topic,
		Message: payload,
	}, nil
}

// Retorna a mensagem de reserva de rotas para ser enviada ao servidor via MQTT
func (s *CarState) SelectRouteMessage(car types.Car, selectedStations []types.Station) (types.MQTT_Message, error) {
	topic := types.CarSelectRouteTopic(s.ServerIP, s.Car.GetCarID())

	message := types.SelectRouteMessage{
		Car:          car,
		StationsList: selectedStations,
	}

	payload, err := json.Marshal(message)
	if err != nil {
		return types.MQTT_Message{}, err
	}

	return types.MQTT_Message{
		Topic:   topic,
		Message: payload,
	}, nil
}

// Recebe as cidades de origem e destino através do terminal
func CityInput() (string, string) {
	city1, city2 := "", ""
	for {
		fmt.Println("Insira a primeira cidade (A, B, C, D, E ou F):")
		fmt.Scanln(&city1)

		fmt.Println("Insira a segunda cidade (A, B, C, D, E ou F):")
		fmt.Scanln(&city2)

		if city1 == city2 {
			fmt.Println("As cidades devem ser diferentes.")
			continue
		}
		break
	}

	return city1, city2
}

func UnmarshalListRoutes(msg paho.Message) map[string][]types.Station {
	// Deserializa a mensagem recebida
	mqttMessage := &types.MQTT_Message{}
	err := json.Unmarshal(msg.Payload(), mqttMessage)
	if err != nil {
		fmt.Println("Error unmarshalling message:", err)
		return nil
	}

	availableStations := make(map[string][]types.Station)
	err = json.Unmarshal(mqttMessage.Message, &availableStations)
	if err != nil {
		fmt.Println("Error unmarshalling message:", err)
		return nil
	}

	// Lista as rotas no terminal
	if len(availableStations) == 0 {
		fmt.Println("Não há rotas disponiveis entre as cidaddes escolhidas")
	} else {
		for company, stationList := range availableStations {
			fmt.Print("Companhia: ", company, " - Estações: \n")
			for _, station := range stationList {
				fmt.Printf("StationID: %d\n", station.StationID)
			}
		}
	}
	return availableStations
}

func UnmarshalListRoutesSelect(msg paho.Message) ([]types.Station, error) {
	selectedStations := []types.Station{}

	availableStations := UnmarshalListRoutes(msg)
	for company, stationList := range availableStations {
		fmt.Printf("Escolha uma estação da companhia %s para reservar(insira o número à esquerda):\n", company)

		for i, station := range stationList {
			fmt.Printf(" %d - StationID: %d\n", i, station.StationID)
		}
		selectedStation := ""
		fmt.Scanln(&selectedStation)
		selectedStationInt, err := strconv.Atoi(selectedStation)
		if err != nil || selectedStationInt < 0 {
			return []types.Station{}, fmt.Errorf("invalid route selection")
		} else {
			selectedStations = append(selectedStations, stationList[selectedStationInt])
		}
	}

	fmt.Println("Estações selecionadas:")
	for i, station := range selectedStations {
		fmt.Printf("%d: %d\n", i, station.StationID)
	}

	return selectedStations, nil
}
