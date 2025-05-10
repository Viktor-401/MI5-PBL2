package main

import (
	"encoding/json"
	"fmt"
	mqtt "mqtt_config/mqtt"
	types "mqtt_config/types"

	paho "github.com/eclipse/paho.mqtt.golang"
)

type Station struct {
	StationID  int
	ServerIP   string
	ReservedBy int // ID do cliente carro que reservou o posto
	Mqtt       *mqtt.MQTT
}

func main() {
	serverIP, stationID := "", 0
	fmt.Println("Insira o IP do server/empresa a qual esse posto pertence:")
	fmt.Scanln(&serverIP)
	fmt.Printf("Insira o ID do posto:")
	fmt.Scanln(&stationID)

	fmt.Printf(`Informações do posto:
	Posto ID: %d
	IP do Servidor: %s`, stationID, serverIP)

	mqttClient, err := mqtt.NewMQTTClient(types.PORT, types.BROKER)
	if err != nil {
		fmt.Println("Error creating MQTT client:", err)
		return
	}

	station := Station{
		StationID:  stationID,
		ServerIP:   serverIP,
		ReservedBy: -1,
		Mqtt:       mqttClient,
	}

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

	// Subscribe to the topic
	topic := types.StationConsultTopic(station.ServerIP, station.StationID)
	station.Mqtt.Subscribe(topic, func(client paho.Client, msg paho.Message) {
		message := types.MQTT_Message{
			Topic: msg.Topic(),
		}

		station.Mqtt.Publish(message)

		mqttMessage := types.MQTT_Message{}
		json.Unmarshal(msg.Payload(), &mqttMessage)
		fmt.Printf("Topic: %s\n", mqttMessage.Topic)
		fmt.Printf("Message: %s\n", mqttMessage.Message)
	})
}

func (s *Station) BirthMessage() (types.MQTT_Message, error) {
	topic := types.StationBirthTopic(s.ServerIP)

	payload, err := json.Marshal(s)
	if err != nil {
		return types.MQTT_Message{}, err
	}

	return types.MQTT_Message{
		Topic:   topic,
		Message: payload,
	}, nil
}
