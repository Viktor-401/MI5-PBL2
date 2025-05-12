package controller

import (
	"api/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ServerController struct {
	serverUsecase usecase.ServerUsecase
}

func NewServerController(usecase usecase.ServerUsecase) ServerController {
	return ServerController{
		serverUsecase: usecase,
	}
}

func (sc *ServerController) GetStationsFromServer(ctx *gin.Context) {
	serverURL := ctx.Query("server_url")
	if serverURL == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "server_url is required"})
		return
	}

	stations, err := sc.serverUsecase.GetStationsFromServer(serverURL)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"stations": stations})
}

func (sc *ServerController) ReserveStationOnServer(ctx *gin.Context) {
	// Captura os parâmetros da query string
	serverURL := ctx.Query("server_url")
	stationID := ctx.Query("station_id")
	carID := ctx.Query("car_id")

	// Valida os parâmetros
	if serverURL == "" || stationID == "" || carID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "server_url, station_id, and car_id are required"})
		return
	}

	// Converte stationID e carID para inteiros
	stationIDInt, err := strconv.Atoi(stationID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "station_id must be an integer"})
		return
	}

	carIDInt, err := strconv.Atoi(carID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "car_id must be an integer"})
		return
	}

	// Chama o usecase para realizar a reserva
	err = sc.serverUsecase.ReserveStationOnServer(serverURL, stationIDInt, carIDInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Station reserved successfully"})
}
