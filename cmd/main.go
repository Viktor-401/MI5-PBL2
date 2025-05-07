package main

import (
	"fmt"
	"main/controller"
	"main/database"
	"main/repository"
	"main/usecase"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Conecta ao MongoDB
	database.ConnectDB()
	defer database.DisconnectDB()

	// Usa o banco de dados dinâmico
	db := database.Database

	// Configura o repositório
	stationRepo := repository.NewStationRepository(db)

	// Popula o banco de dados
	database.SeedData(stationRepo)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Fallback para a porta 8080 se a variável de ambiente não estiver configurada
	}

	// Configurações do servidor
	server := gin.Default()
	stationUsecase := usecase.NewStationUseCase(stationRepo)
	serverUsecase := usecase.NewServerUsecase()
	stationController := controller.NewStationController(stationUsecase)
	serverController := controller.NewServerController(serverUsecase)

	// Rotas relacionadas às estações
	server.POST("/stations", stationController.CreateStation)
	server.GET("/stations", stationController.GetAllStations)

	// Rotas relacionadas à comunicação entre servidores
	server.GET("/server/stations", serverController.GetStationsFromServer)
	server.POST("/server/reserve", serverController.ReserveStationOnServer)

	server.POST("/stations/reserve", stationController.ReserveStation)

	server.Run(fmt.Sprintf(":%s", port))
}
