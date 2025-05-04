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

    // cria servidor HTTP Gin
    r := gin.Default()
    r.GET("/cars", GetCars)
    r.GET("/stations", GetStations)
    // (aqui você adiciona /reservations/prepare, /commit, etc.)
    r.Run(":8080") // ajusta porta por empresa (8080, 8081, 8082)
}
