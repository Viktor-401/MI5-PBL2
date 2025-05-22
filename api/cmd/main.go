package main

import (
	"api/controller"
	"api/database"
	"api/mqtt_server"
	"api/repository"
	"api/usecase"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Conecta ao MongoDB
	database.ConnectDB()
	defer database.DisconnectDB()

	// Usa o banco de dados dinâmico
	db := database.Database

	// Popula o banco de dados com rotas
	database.SeedRoutes(db)

	// Configura o repositório
	stationRepo := repository.NewStationRepository(db)
	//database.SeedData(stationRepo)
	serverRepo := repository.NewServerRepository(db)
	routeRepo := repository.NewRouteRepository(db)

	// Configura o usecase
	routeUsecase := usecase.NewRouteUsecase(routeRepo)
	stationUsecase := usecase.NewStationUseCase(stationRepo)
	serverUsecase := usecase.NewServerUsecase(serverRepo)

	// Configura o controlador
	routeController := controller.NewRouteController(routeUsecase)
	stationController := controller.NewStationController(stationUsecase)
	serverController := controller.NewServerController(serverUsecase)

	server := gin.Default()
	// Rotas relacionadas às estações locais
	server.POST("/stations", stationController.CreateStation)
	server.GET("/stations", stationController.GetAllStations)
	server.GET("/stations/:id", stationController.GetStationByID)
	server.PUT("/stations/:id/remove", stationController.RemoveStation)
	server.PUT("/stations/:id/release", stationController.ReleaseStation)

	// Rotas relacionadas às rotas predefinidas
	server.POST("/routes", routeController.CreateRoute)
	server.GET("/routes", routeController.GetRoutes)

	// Rotas relacionadas ao 2PC no servidor local
	server.PUT("/stations/:id/prepare", stationController.PrepareStation)
	server.PUT("/stations/:id/commit", stationController.CommitStation)

	// Rotas relacionadas à comunicação entre servidores (database dos IPs dos servidores)
	server.POST("/servers/register", serverController.RegisterServer)
	server.GET("/servers/:id", serverController.GetServerByCompany)

	// Rotas relacionadas ao servidor remoto
	server.GET("/server/:sid/stations", serverController.GetStationsFromServer)
	server.PUT("/server/:sid/stations/:id/prepare", serverController.PrepareStationOnServer)
	server.PUT("/server/:sid/stations/:id/commit", serverController.CommitStationOnServer)
	server.PUT("/server/:sid/stations/:id/release", serverController.ReleaseStationOnServer)

	// Inicia o servidor MQTT
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Fallback para a porta 8080 se a variável de ambiente não estiver configurada
	}
	var company = ""
	fmt.Println("Insira a empresa do Servidor:")
	fmt.Scanln(&company)

	// Gorutina para o servidor MQTT
	go mqtt_server.MqttMain(company, port)
	// Inicia o Servidor HTTP
	server.Run(fmt.Sprintf(":%s", port))

}
