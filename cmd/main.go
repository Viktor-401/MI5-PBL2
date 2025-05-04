package main

import (
	"main/controller"
	"main/repository"
	"main/usecase"
	"main/database"

	"github.com/gin-gonic/gin"
)
func main() {
    // Conecta ao MongoDB
    database.ConnectDB()
    defer database.DisconnectDB()

    db := database.MongoClient.Database("mydb")

    // Configura o repositório
    stationRepo := repository.NewStationRepository(db)

    // Popula o banco de dados
    database.SeedData(stationRepo)

    // Configurações do servidor
    server := gin.Default()
    stationUsecase := usecase.NewStationUseCase(stationRepo)
    stationController := controller.NewStationController(stationUsecase)

    server.POST("/stations", stationController.CreateStation)
    server.GET("/stations", stationController.GetAllStations)

    server.Run(":8080")
}

