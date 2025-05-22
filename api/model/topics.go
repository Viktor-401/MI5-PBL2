package model

import (
	"fmt"
)

// Tópicos que identificam a ação a ser realizada
type Topics int

const (
	Consult Topics = iota
	Reserve
	Select
	Birth
	Death
	Finish
)

var TopicNames = map[Topics]string{
	Consult: "consult",
	Reserve: "reserve",
	Select:  "select",
	Birth:   "birth",
	Death:   "death",
	Finish:  "finish",
}

func (t Topics) String() string {
	return TopicNames[t]
}

// Tipos de clientes que se conectam no servidor MQTT
type MqttClientTypes int

const (
	StationClientType MqttClientTypes = iota
	CarClientType
	CompanyClientType
)

var MqttClientTypeNames = map[MqttClientTypes]string{
	StationClientType: "station",
	CarClientType:     "car",
	CompanyClientType: "company",
}

func (m MqttClientTypes) String() string {
	return MqttClientTypeNames[m]
}

// Mensagem enviada pelos clientes MQTT
type MQTT_Message struct {
	Topic   string `json:"topic"`
	Message []byte `json:"message"`
}

// Mensagem com as informações do carro que deseja reservar uma estação
type CarInfo struct {
	CarId int `json:"car_id"`
}

// Mensagem com as cidades de origem e destino
type RoutesMessage struct {
	City1 string `json:"city1"`
	City2 string `json:"city2"`
}

// Mensagem com uma lista de estações selecionadas pelo cliente carro
type SelectRouteMessage struct {
	Car          Car       `json:"car"`
	StationsList []Station `json:"route"`
}

// Mensagem com uma lista das estações que devem ser liberadas
type FinishRouteMessage struct {
	Car          Car       `json:"car"`
	StationsList []Station `json:"route"`
}

// Mensagem com uma lista de rotas entre duas cidades
type RoutesList struct {
	Routes []Route `json:"routes"`
}

// STATION TOPICS
func StationBirthTopic(serverIP string) string {
	// Birth of a station in serverIP
	return Birth.String() + StationClientType.String() + serverIP
}

func StationDeathTopic(serverIP string) string {
	// Death of a station in serverIP
	return Death.String() + StationClientType.String() + serverIP
}

func StationConsultTopic(serverIP string, stationID int) string {
	// Consult a station in serverIP with stationID
	return Consult.String() + StationClientType.String() + serverIP + fmt.Sprintf("%d", stationID)
}

func StationReserveTopic(serverIP string, stationID int) string {
	// Reserve a station in serverIP with stationID
	return Reserve.String() + StationClientType.String() + serverIP + fmt.Sprintf("%d", stationID)
}

// CAR TOPICS
func CarBirthTopic(serverIP string) string {
	// Birth of a Car in serverIP
	return Birth.String() + CarClientType.String() + serverIP
}

func CarDeathTopic(serverIP string) string {
	// Death of a Car in serverIP
	return Death.String() + CarClientType.String() + serverIP
}

func CarConsultTopic(serverIP string, CarID int) string {
	// Consult a Car in serverIP with CarID
	return Consult.String() + CarClientType.String() + serverIP + fmt.Sprintf("%d", CarID)
}

func CarReserveTopic(serverIP string, CarID int) string {
	// Reserve a Car in serverIP with CarID
	return Reserve.String() + CarClientType.String() + serverIP + fmt.Sprintf("%d", CarID)
}

func CarSelectRouteTopic(serverIP string, CarID int) string {
	// Select a route for a Car in serverIP with CarID
	return Select.String() + CarClientType.String() + serverIP + fmt.Sprintf("%d", CarID)
}

func FinishRouteTopic(serverIP string, carID int) string {
	return Finish.String() + CarClientType.String() + serverIP + fmt.Sprintf("%d", carID)
}

func ResponseFinishRouteTopic(serverIP string, carID int) string {
	return "response" + Finish.String() + CarClientType.String() + serverIP + fmt.Sprintf("%d", carID)
}

// SERVER TOPICS
func ServerBirthTopic(serverIP string) string {
	// Birth of a server in serverIP
	return Birth.String() + CompanyClientType.String()
}

func ResponseServerBirthTopic(serverIP string) string {
	// Birth of a server in serverIP
	return "response" + Birth.String() + CompanyClientType.String()
}

func ResponseCarConsultTopic(serverIP string, CarID int) string {
	// Consult a Car in serverIP with CarID
	return "response" + Consult.String() + CarClientType.String() + serverIP + fmt.Sprintf("%d", CarID)
}

func ResponseCarReserveTopic(serverIP string, CarID int) string {
	// Reserve a Car in serverIP with CarID
	return "response" + Reserve.String() + CarClientType.String() + serverIP + fmt.Sprintf("%d", CarID)
}

func ResponseStationReserveTopic(serverIP string, stationID string) string {
	// Reserve a station in serverIP with stationID
	return "response" + Reserve.String() + StationClientType.String() + serverIP + stationID
}
