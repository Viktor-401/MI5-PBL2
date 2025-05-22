package main

import (
	mqtt "clients/mqtt"
	types "clients/types"
	"encoding/json"
	"fmt"
	"math/rand"
	"slices"
	"strconv"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
)

type CarState struct {
	Car      types.Car
	ServerIP string
	Mqtt     *mqtt.MQTT
}

var carState CarState = CarState{}
var waitChan chan bool = make(chan bool, 1)
var autoClient bool = false

func main() {

	// Cria o cliente MQTT
	mqttClient, err := mqtt.NewMQTTClient(types.PORT, types.BROKER)
	if err != nil {
		fmt.Println("Error creating MQTT client:", err)
		return
	}

	autoClientInput := ""
	fmt.Println("Client carro automatizado?(y/n)")
	fmt.Scanln(&autoClientInput)
	if autoClientInput == "y" {
		autoClient = true
	}

	serverIP := ""
	fmt.Println("Insira o IP do server a qual esse carro vai se conectar:")
	fmt.Scanln(&serverIP)
	car := types.GetNewRandomCar()

	carState.Car = car
	carState.ServerIP = serverIP
	carState.Mqtt = mqttClient

	// Mensagem de nascimento do carro, que informa o servidor que o carro está online
	err = carState.BirthMessage()
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
	carState.Mqtt.Subscribe(topic, ResponseCarReserveCallback)

	// Inscrição no tópico de liberação de rotas
	topic = types.ResponseFinishRouteTopic(serverIP, carState.Car.GetCarID())
	carState.Mqtt.Subscribe(topic, ResponseFinishRouteCallback)

	exit := false
	city1, city2 := "", ""
	for !exit {
		action := 0
		// Menu de ações
		for {
			fmt.Println(` Escolha uma ação:
        1- Consultar Rotas
        2- Reservar Postos
        3- Liberar Rota`)
			// Cliente Automatico
			if autoClient {
				action = rand.Intn(3) + 1
				fmt.Printf("Ação automática: %d\n", action)
				time.Sleep(2 * time.Second)
				break
				// Cliente Normal
			} else {
				fmt.Scanln(&action)
				if action == 1 || action == 2 || action == 3 {
					break
				}
				fmt.Println("Ação inválida. Tente novamente.")
			}
		}

		// Consultar Estações
		if action == 1 {
			city1, city2 = CityInput()

			err = carState.ConsultRouteMessage(city1, city2)
			if err != nil {
				fmt.Println("Erro em ConsultRoutMessage:", err)
				continue
			}

			// Reservar Estações
		} else if action == 2 {
			city1, city2 = CityInput()

			// Envia a mensagem de reserva de rotas
			err = carState.ReserveRouteMessage(city1, city2)
			if err != nil {
				fmt.Println("Erro em ConsultRoutMessage:", err)
				continue
			}

			// Liberar Rota
		} else if action == 3 {
			if len(carState.Car.ReservedStations) == 0 {
				fmt.Println("Nenhuma rota reservada para liberar.")
				continue
			}
			stations := carState.Car.ReservedStations
			// Envia mensagem para liberação de postos
			err = carState.FinishRouteMessage(stations)
			if err != nil {
				fmt.Println("Erro em ConsultRoutMessage:", err)
				continue
			}
		} else {
			fmt.Println("Ação inválida. Tente novamente.")
		}
		// Espera até que as respostas do servidor sejam tratadas pelo cliente
		<-waitChan
	}

	// Mantem o cliente MQTT ativo até o usuário encerrar
	fmt.Println("Enter para encerra o posto")
	fmt.Scanln()
	// Mensagem de morte do posto, que informa o servidor que o posto está offline
	err = carState.DeathMessage()
	if err != nil {
		fmt.Println("Error publishing death message:", err)
		return
	}
}

// Retorna a mensagem de nascimento do carro, que informa o servidor que o carro está online
func (s *CarState) BirthMessage() error {
	topic := types.CarBirthTopic(s.ServerIP)

	payload, err := json.Marshal(s.Car)
	if err != nil {
		return err
	}

	message := types.MQTT_Message{
		Topic:   topic,
		Message: payload,
	}

	err = carState.Mqtt.Publish(message)
	return err
}

// Retorna a mensagem de morte do carro, que informa o servidor que o carro está offline
func (s *CarState) DeathMessage() error {
	topic := types.CarDeathTopic(s.ServerIP)

	payload, err := json.Marshal(s.Car)
	if err != nil {
		return err
	}

	message := types.MQTT_Message{
		Topic:   topic,
		Message: payload,
	}

	err = carState.Mqtt.Publish(message)
	return err
}

// Retorna a mensagem de consulta de rotas para ser enviada ao servidor via MQTT
func (s *CarState) ConsultRouteMessage(city1 string, city2 string) error {
	topic := types.CarConsultTopic(s.ServerIP, s.Car.GetCarID())

	consultRoute := types.RoutesMessage{
		City1: city1,
		City2: city2,
	}

	payload, err := json.Marshal(consultRoute)
	if err != nil {
		return err
	}

	message := types.MQTT_Message{
		Topic:   topic,
		Message: payload,
	}

	err = carState.Mqtt.Publish(message)
	return err
}

// Retorna a mensagem de reserva de rotas para ser enviada ao servidor via MQTT
func (s *CarState) ReserveRouteMessage(city1 string, city2 string) error {
	topic := types.CarReserveTopic(s.ServerIP, s.Car.GetCarID())

	reserveRoute := types.RoutesMessage{
		City1: city1,
		City2: city2,
	}

	payload, err := json.Marshal(reserveRoute)
	if err != nil {
		return err
	}

	message := types.MQTT_Message{
		Topic:   topic,
		Message: payload,
	}

	err = carState.Mqtt.Publish(message)
	return err
}

// Retorna a mensagem de reserva de rotas para ser enviada ao servidor via MQTT
func (s *CarState) SelectRouteMessage(car types.Car, selectedStations []types.Station) error {
	topic := types.CarSelectRouteTopic(s.ServerIP, s.Car.GetCarID())

	message := types.SelectRouteMessage{
		Car:          car,
		StationsList: selectedStations,
	}

	payload, err := json.Marshal(message)
	if err != nil {
		return err
	}

	mqttMessage := types.MQTT_Message{
		Topic:   topic,
		Message: payload,
	}

	err = carState.Mqtt.Publish(mqttMessage)
	return err
}

// Retorna a mensagem de finalização de rota para liberar as estações reservadas
func (s *CarState) FinishRouteMessage(stations []types.Station) error {
	topic := types.FinishRouteTopic(s.ServerIP, s.Car.GetCarID())

	payloadStruct := struct {
		Car          types.Car       `json:"car"`
		StationsList []types.Station `json:"route"`
	}{
		Car:          s.Car,
		StationsList: stations,
	}

	payload, err := json.Marshal(payloadStruct)
	if err != nil {
		return err
	}

	message := types.MQTT_Message{
		Topic:   topic,
		Message: payload,
	}

	err = carState.Mqtt.Publish(message)
	return err
}

// Recebe as cidades de origem e destino através do terminal
func CityInput() (string, string) {
	city1, city2 := "", ""
	if autoClient {
		cityList := []string{"A", "B", "C", "D"}

		randInt := rand.Intn(3)
		city1 = cityList[randInt]

		cityList = slices.Delete(cityList, randInt, randInt+1)

		randInt = rand.Intn(2)
		city2 = cityList[randInt]

		fmt.Println("Cidades escolhidas automaticamente:")
		fmt.Printf("Cidade 1: %s. Cidade 2: %s\n", city1, city2)
	} else {
		for {
			fmt.Println("Insira a primeira cidade :")
			fmt.Scanln(&city1)

			fmt.Println("Insira a segunda cidade :")
			fmt.Scanln(&city2)

			if city1 == city2 {
				fmt.Println("As cidades devem ser diferentes.")
				continue
			}
			break
		}
	}

	return city1, city2
}

func UnmarshalListRoutes(msg paho.Message) (map[string][]types.Station, bool) {
	// Deserializa a mensagem recebida
	mqttMessage := &types.MQTT_Message{}
	err := json.Unmarshal(msg.Payload(), mqttMessage)
	if err != nil {
		fmt.Println("Error unmarshalling message:", err)
		return nil, false
	}

	availableCompanys := make(map[string][]types.Station)
	err = json.Unmarshal(mqttMessage.Message, &availableCompanys)
	if err != nil {
		fmt.Println("Error unmarshalling message:", err)
		return nil, false
	}

	// Verifica se alguma companhia tem a lista de estações vazia
	hasAvailableRoutes := false

	for company, stationList := range availableCompanys {
		if len(stationList) == 0 {
			fmt.Printf("A companhia '%s' não tem estações disponíveis.\n", company)
			return availableCompanys, false
		} else {
			hasAvailableRoutes = true
			fmt.Print("Companhia: ", company, " - Estações: \n")
			for _, station := range stationList {
				fmt.Printf("StationID: %d\n", station.StationID)
			}
		}
	}

	// Se nenhuma rota disponível for encontrada
	if !hasAvailableRoutes {
		fmt.Println("Não há rotas disponíveis entre as cidades escolhidas.")
	}

	return availableCompanys, hasAvailableRoutes
}

// Recebe input do usuário para selecionar os postos que deseja reservar entre as cidades
func UnmarshalListRoutesSelect(msg paho.Message) ([]types.Station, error) {
	selectedStations := []types.Station{}

	// Recebe o mapa das estações disponiveis
	availableStations, hasAvailableRoutes := UnmarshalListRoutes(msg)
	// Retorna se não houverem rotas
	if !hasAvailableRoutes {
		fmt.Printf("Não há rotas entre as cidades escolhidas")
		return selectedStations, nil
	}
	// Intera entre as companias e armazena as estações escolhida pelo usuário em selectedStations
	for company, stationList := range availableStations {
		fmt.Printf("Escolha uma estação da companhia %s para reservar(insira o número à esquerda):\n", company)
		// Mostra cada posto de cada compania
		for i, station := range stationList {
			fmt.Printf(" %d - StationID: %d\n", i, station.StationID)
		}
		// Auto cliente seleciona um posto da compania aleatoriamente
		if autoClient {
			fmt.Printf("Tamanho da lista: %d", len(stationList))
			randInt := rand.Intn(len(stationList))
			// Guarda estação escolhida aleatoriamente
			selectedStations = append(selectedStations, stationList[randInt])
			// Cĺiente normal
		} else {
			// Recebe input
			selectedStation := ""
			fmt.Scanln(&selectedStation)
			// Trata o input do usuário
			selectedStationInt, err := strconv.Atoi(selectedStation)
			if err != nil || selectedStationInt < 0 {
				return []types.Station{}, fmt.Errorf("invalid route selection")
			} else {
				// Armazena estação escolhida pelo usuário
				selectedStations = append(selectedStations, stationList[selectedStationInt])
			}
		}
	}

	// Finaliza mostrando as estações selecionadas pelo usuário
	fmt.Println("Estações selecionadas:")
	for i, station := range selectedStations {
		fmt.Printf("%d: %d\n", i, station.StationID)
	}

	return selectedStations, nil
}

func ResponseCarReserveCallback(client paho.Client, msg paho.Message) {
	// Funcao de callback para quando uma mensagem é recebida
	selectedStations, err := UnmarshalListRoutesSelect(msg)
	if err != nil {
		fmt.Println("Error unmarshalling route selection:", err)
		return
	}
	carState.Car.ReservedStations = selectedStations
	fmt.Println("Estações reservadas com sucesso!")
	err = carState.SelectRouteMessage(carState.Car, selectedStations)
	if err != nil {
		fmt.Println("Error publishing reserve route message:", err)
		return
	}

	waitChan <- true
}

func ResponseFinishRouteCallback(client paho.Client, msg paho.Message) {
	// Função de callback para quando uma mensagem de resposta de finalização de rota é recebida
	fmt.Println("Resposta de finalização de rota recebida!")
	mqttMessage := &types.MQTT_Message{}
	err := json.Unmarshal(msg.Payload(), mqttMessage)
	if err != nil {
		fmt.Println("Erro ao decodificar MQTT_Message:", err)
		return
	}
	fmt.Println(string(mqttMessage.Message))
	waitChan <- true
}
