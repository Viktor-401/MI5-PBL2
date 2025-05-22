package main

import (
	mqtt "clients/mqtt"
	types "clients/types"
	"encoding/json"
	"fmt"
	"strconv"

	paho "github.com/eclipse/paho.mqtt.golang"
)

// Estrutura do estado do posto
type Station struct {
	StationData types.Station
	Mqtt        *mqtt.MQTT
}

// Instância do estado do posto
var station Station = Station{
	StationData: types.Station{
		InUseBy:  -1,
		IsActive: true,
	},
}

// Fluxo principal
func main() {
	// Input das informações do posto
	serverIP, stationID := "", 0
	fmt.Println("Insira o IP do server/empresa a qual esse posto pertence:")
	fmt.Scanln(&serverIP)
	for {

		fmt.Printf("Insira o ID do posto:")
		var input string
		fmt.Scanln(&input)

		id, err := strconv.Atoi(input)
		if err == nil {
			stationID = id
			break
		} else {
			fmt.Println("Valor inválido! Por favor, insira um número inteiro.")
		}
	}
	fmt.Printf(`Informações do posto:
	Posto ID: %d
	IP do Servidor: %s`, stationID, serverIP)

	// Cria o cliente MQTT
	mqttClient, err := mqtt.NewMQTTClient(types.PORT, types.BROKER)
	if err != nil {
		fmt.Println("Error creating MQTT client:", err)
		return
	}

	// Atribuições dos dados inseridos pelo usuário e clíente mqtt
	station.StationData.StationID = stationID
	station.StationData.ServerIP = serverIP
	station.Mqtt = mqttClient

	// Mensagem de nascimento do posto, que informa o servidor que o posto está online
	birthMessage, err := station.BirthMessage()
	if err != nil {
		fmt.Println("Error creating birth message:", err)
		return
	}
	// Publicação da mensagem de nascimento
	err = station.Mqtt.Publish(birthMessage)
	if err != nil {
		fmt.Println("Error publishing birth message:", err)
		return
	}

	// Topico e atribuição da função de callback do posto para reservar o posto
	topic := types.ResponseStationReserveTopic(station.StationData.ServerIP, fmt.Sprintf("%d", station.StationData.StationID))
	station.Mqtt.Subscribe(topic, ResponseStationReserveCallback)

	// Mantem o cliente MQTT ativo até o usuário encerrar
	fmt.Println("\nEnter para encerra o posto")
	fmt.Scanln()

	// Mensagem de morte do posto, que informa o servidor que o posto está offline
	message, err := station.DeathMessage()
	if err != nil {
		fmt.Println("Error creating death message:", err)
		return
	}
	station.Mqtt.Publish(message)
}

// Constroi a mensagem de nascimento do posto
func (s *Station) BirthMessage() (types.MQTT_Message, error) {
	topic := types.StationBirthTopic(s.StationData.ServerIP)

	// Serializa as informações do posto para serem enviadas pelo mqtt
	payload, err := json.Marshal(s.StationData)
	if err != nil {
		return types.MQTT_Message{}, err
	}

	return types.MQTT_Message{
		Topic:   topic,
		Message: payload,
	}, nil
}

// Constroi a mensagem de morte do posto
func (s *Station) DeathMessage() (types.MQTT_Message, error) {
	topic := types.StationDeathTopic(s.StationData.ServerIP)

	// Serializa as informações do posto para serem enviadas pelo mqtt
	payload, err := json.Marshal(s.StationData.StationID)
	if err != nil {
		return types.MQTT_Message{}, err
	}

	return types.MQTT_Message{
		Topic:   topic,
		Message: payload,
	}, nil
}

// Salva o id do carro recebido de um servidor no estado do posto
func ResponseStationReserveCallback(client paho.Client, msg paho.Message) {
	// Desserializa a mensagem do mqtt
	message := types.MQTT_Message{}
	err := json.Unmarshal(msg.Payload(), &message)
	if err != nil {
		fmt.Println("Error unmarshalling message:", err)
		return
	}
	// Desserializa o campo de dados da mensagem
	carInfo := types.CarInfo{}
	err = json.Unmarshal(message.Message, &carInfo)
	if err != nil {
		fmt.Println("Error unmarshalling car info:", err)
		return
	}

	// Atualiza o ID do carro que reservou o posto
	station.StationData.InUseBy = carInfo.CarId
	fmt.Printf("Posto %d reservado pelo carro %d\n", station.StationData.StationID, carInfo.CarId)
}
