package main

import (
	"encoding/json"
	"fmt"
	mqtt "main/mqtt"
	types "main/types"
	"strconv"

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

	// Inscrição no tópico de consulta de rotas
	topic := types.CarConsultTopic(serverIP, carState.Car.GetCarID())
	carState.Mqtt.Subscribe(topic, func(client paho.Client, msg paho.Message) {
		// Funcao de callback para quando uma mensagem é recebida
		// Deve retornar uma mensagem com payload ListRoutes
	})

	// Inscrição no tópico de reserva de rotas
	topic = types.CarReserveTopic(serverIP, carState.Car.GetCarID())
	carState.Mqtt.Subscribe(topic, func(client paho.Client, msg paho.Message) {
		// Funcao de callback para quando uma mensagem é recebida
		// Deve retornar uma mensagem com payload ListRoutes
	})

	// Mantem o cliente MQTT ativo até o usuário encerrar
	fmt.Println("Enter para encerra o posto")
	fmt.Scanln()
	// Mensagem de morte do servidor, que informa o seu subscribers que o servidor está offline
	message, err := serverState.DeathMessage()
	if err != nil {
		fmt.Println("Error creating death message:", err)
		return
	}
	serverState.Mqtt.Publish(message)
}

// Retorna a mensagem de morte do carro, que informa o servidor que o carro está offline
func (s *ServerState) DeathMessage() (types.MQTT_Message, error) {
	topic := types.StationDeathTopic(s.ServerIP)

	payload, err := json.Marshal(s)
	if err != nil {
		return types.MQTT_Message{}, err
	}

	return types.MQTT_Message{
		Topic:   topic,
		Message: payload,
	}, nil
}

// Retorna a mensagem de consulta de rotas para ser enviada ao servidor via MQTT
func (s *ServerState) ConsultRouteMessage(city1 string, city2 string) (types.MQTT_Message, error) {
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
func (s *ServerState) ReserveRouteMessage(city1 string, city2 string) (types.MQTT_Message, error) {
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
func (s *ServerState) SelectRouteMessage(route types.Route) (types.MQTT_Message, error) {
	topic := types.CarReserveTopic(s.ServerIP, s.Car.GetCarID())

	payload, err := json.Marshal(route)
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

func UnmarshalListRoutes(msg paho.Message) types.RoutesList {
	// Deserializa a mensagem recebida
	routesMessage := &types.RoutesList{}
	err := json.Unmarshal(msg.Payload(), &routesMessage)
	if err != nil {
		fmt.Println("Error unmarshalling message:", err)
		return types.RoutesList{}
	}

	// Lista as rotas no terminal
	for i, route := range routesMessage.Routes {
		fmt.Printf("%d: %s -> %s\n", i, route.StartCity, route.EndCity)
	}

	return *routesMessage
}

func UnmarshalListRoutesSelect(msg paho.Message) (types.Route, error) {
	routesMessage := UnmarshalListRoutes(msg)
	fmt.Println("Escolha uma rota para reservar:")
	selectedRoute := ""
	fmt.Scanln(&selectedRoute)

	selectedRouteInt, err := strconv.Atoi(selectedRoute)
	if err != nil || selectedRouteInt < 0 || selectedRouteInt >= len(routesMessage.Routes) {
		return types.Route{}, fmt.Errorf("invalid route selection")
	}

	return routesMessage.Routes[selectedRouteInt], nil
}
