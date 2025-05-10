package global

import (
	"fmt"
)

type Topics int

const (
	Consult Topics = iota
	Reserve
	Birth
	Death
)

var TopicNames = map[Topics]string{
	Consult: "consult",
	Reserve: "reserve",
	Birth:   "birth",
	Death:   "death",
}

func (t Topics) String() string {
	return TopicNames[t]
}

type MqttClientTypes int

const (
	Station MqttClientTypes = iota
	Car
	Company
)

var MqttClientTypeNames = map[MqttClientTypes]string{
	Station: "station",
	Car:     "car",
	Company: "company",
}

func (m MqttClientTypes) String() string {
	return MqttClientTypeNames[m]
}

type MQTT_Message struct {
	Topic   string `json:"topic"`
	Message []byte `json:"message"`
}

func StationBirthTopic(serverIP string) string {
	return Birth.String() + Station.String() + serverIP
}

func StationDeathTopic(serverIP string) string {
	return Death.String() + Station.String() + serverIP
}

func StationConsultTopic(serverIP string, stationID int) string {
	return Consult.String() + Station.String() + serverIP + fmt.Sprintf("%d", stationID)
}
