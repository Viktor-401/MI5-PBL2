package main

import (
	mqtt "clients/mqtt"
	types "clients/types"
	"encoding/json"
	"fmt"

	paho "github.com/eclipse/paho.mqtt.golang"
)

type Station struct {
	StationData types.Station
	Mqtt        *mqtt.MQTT
}

func main() {
	// Input das informações do posto
	serverIP, stationID := "", 0
	fmt.Println("Insira o IP do server/empresa a qual esse posto pertence:")
	fmt.Scanln(&serverIP)
	fmt.Printf("Insira o ID do posto:")
	fmt.Scanln(&stationID)

	fmt.Printf(`Informações do posto:
	Posto ID: %d
	IP do Servidor: %s\n`, stationID, serverIP)

	// Cria o cliente MQTT
	mqttClient, err := mqtt.NewMQTTClient(types.PORT, types.BROKER)
	if err != nil {
		fmt.Println("Error creating MQTT client:", err)
		return
	}
	// Estado do posto
	station := Station{
		StationData: types.Station{
			StationID: stationID,
			ServerIP:  serverIP,
			Company:   "",
			InUseBy:   -1,
			IsActive:  true,
		},
		Mqtt: mqttClient,
	}

	// Mensagem de nascimento do posto, que informa o servidor que o posto está online
	birthMessage, err := station.BirthMessage()
	if err != nil {
		fmt.Println("Error creating birth message:", err)
		return
	}

	err = station.Mqtt.Publish(birthMessage)
	if err != nil {
		fmt.Println("Error publishing birth message:", err)
		return
	}

	// Topico para reservar o posto
	topic := types.ResponseStationReserveTopic(station.StationData.ServerIP, fmt.Sprintf("%d", station.StationData.StationID))
	// Inscrição no tópico de reserva, e atribui a função de callback
	station.Mqtt.Subscribe(topic, func(client paho.Client, msg paho.Message) {
		message := types.MQTT_Message{} 
		err := json.Unmarshal(msg.Payload(), &message)
		if err != nil {
			fmt.Println("Error unmarshalling message:", err)
			return
		}
		carInfo := types.CarInfo{}
		// Decodifica a mensagem recebida
		err = json.Unmarshal(message.Message, &carInfo)
		if err != nil {
			fmt.Println("Error unmarshalling car info:", err)
			return
		}
		// Atualiza o ID do carro que reservou o posto
		station.StationData.InUseBy = carInfo.CarId
		fmt.Printf("Posto %d reservado pelo carro %d\n", station.StationData.StationID, carInfo.CarId)
		
	})

	// Mantem o cliente MQTT ativo até o usuário encerrar
	fmt.Println("Enter para encerra o posto")
	fmt.Scanln()
	// Mensagem de morte do posto, que informa o servidor que o posto está offline
	message, err := station.DeathMessage()
	if err != nil {
		fmt.Println("Error creating death message:", err)
		return
	}
	station.Mqtt.Publish(message)
}

func (s *Station) BirthMessage() (types.MQTT_Message, error) {
	topic := types.StationBirthTopic(s.StationData.ServerIP)

	payload, err := json.Marshal(s.StationData)
	if err != nil {
		return types.MQTT_Message{}, err
	}

	return types.MQTT_Message{
		Topic:   topic,
		Message: payload,
	}, nil
}

func (s *Station) DeathMessage() (types.MQTT_Message, error) {
	topic := types.StationDeathTopic(s.StationData.ServerIP)

	payload, err := json.Marshal(s.StationData.StationID)
	if err != nil {
		return types.MQTT_Message{}, err
	}

	return types.MQTT_Message{
		Topic:   topic,
		Message: payload,
	}, nil
}
