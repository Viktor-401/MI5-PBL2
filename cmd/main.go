package main

import "github.com/gin-gonic/gin"

func main() {
    // conecta no Mongo e garante desconexão ao sair
    ConnectDB()
    defer DisconnectDB()

    // popula dados iniciais
    SeedData()

    // exemplo de uso da lógica de negócio
    UpdateCarBattery(3, 90)
    UpdateCarBattery(7, 50)

    // cria servidor HTTP Ginserver := gin.Default()

	//Realizar a conexão com o banco de dados e conectar com o repositório

	//Camada de repositórios
	StationRepository := repository.NewStationRepository(nil) //
	//Camada de usecases
	StationUsecase := usecase.NewStationUseCase(StationRepository)
	//Camada de controllers
	StationController := controller.NewStationController(StationUsecase)

	server.POST("/stations", StationController.CreateStation)
	server.GET("/stations", StationController.GetAllStations)

    // (aqui você adiciona /reservations/prepare, /commit, etc.)
    r.Run(":8080") // ajusta porta por empresa (8080, 8081, 8082)
}

