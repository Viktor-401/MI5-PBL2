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
	// Rotas relacionadas às estações
	server.POST("/stations", stationController.CreateStation)
	server.GET("/stations", stationController.GetAllStations)
	server.GET("/stations/:id", stationController.GetStationByID)
	server.POST("/stations/:id/remove", stationController.RemoveStation)
	server.POST("/stations/:id/reserve", stationController.ReserveStation)
	//server.POST("/stations/:id/prepare", stationController.PrepareStation)

	// Rotas relacionadas ao 2PC
	server.POST("/stations/:id/prepare", stationController.PrepareStation)
	server.POST("/stations/:id/commit", stationController.CommitStation)
	//server.POST("/stations/:id/abort", stationController.AbortStation)

	// Rotas relacionadas à comunicação entre servidores
	server.POST("/servers/register", serverController.RegisterServer)
	server.GET("/servers/:id", serverController.GetServerByCompany)

	// Rotas relacionadas ao servidor remoto
	server.GET("/server/:sid/stations", serverController.GetStationsFromServer)

	server.POST("/server/:sid/stations/:id/reserve/", serverController.ReserveStationOnServer)
	server.POST("/server/:sid/stations/:id/prepare/", serverController.PrepareStationOnServer)
	server.POST("/server/:sid/stations/:id/commit/", serverController.CommitStationOnServer)

	// Rotas relacionadas às rotas
	server.POST("/routes", routeController.CreateRoute)
	server.GET("/routes", routeController.GetRoutes)

	// Inicia o servidor MQTT

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Fallback para a porta 8080 se a variável de ambiente não estiver configurada
	}
	var company = ""
	fmt.Println("Insira a empresa do Servidor:")
	fmt.Scanln(&company)

	go mqtt_server.MqttMain(company, port)
	server.Run(fmt.Sprintf(":%s", port))

}
