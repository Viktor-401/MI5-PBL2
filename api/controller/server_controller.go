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

func (sc *ServerController) RegisterServer(ctx *gin.Context) {
	var request struct {
		Company  string `json:"company"`
		ServerIP string `json:"server_ip"`
	}

	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := sc.serverUsecase.RegisterOrUpdateServer(request.Company, request.ServerIP)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Servidor registrado com sucesso"})
}

func (sc *ServerController) GetServersByCompany(ctx *gin.Context) {
	// Obtém o parâmetro "company" da query string
	company := ctx.Query("company")
	if company == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "company name is required"})
		return
	}

	// Chama o usecase para buscar os servidores pela companhia
	servers, err := sc.serverUsecase.GetServersByCompany(company)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Retorna os servidores encontrados
	ctx.JSON(http.StatusOK, servers)
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

	ctx.JSON(http.StatusOK, stations)
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
