package main

import (
	"main/controller"
	"main/repository"
	"main/usecase"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	//Realizar a conexão com o banco de dados e conectar com o repositório

	//Camada de repositórios
	StationRepository := repository.NewStationRepository(nil) //
	//Camada de usecases
	StationUsecase := usecase.NewStationUseCase(StationRepository)
	//Camada de controllers
	StationController := controller.NewStationController(StationUsecase)

	server.POST("/stations", StationController.CreateStation)
	server.GET("/stations", StationController.GetAllStations)

	server.Run(":8000")

}
