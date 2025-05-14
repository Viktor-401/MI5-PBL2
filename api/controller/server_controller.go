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

func (sc *ServerController) GetServerByCompany(ctx *gin.Context) {
	// Obtém o parâmetro "id" da URL
	companyID := ctx.Param("id")
	if companyID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "company ID is required"})
		return
	}

	// Chama o usecase para buscar os servidores pela companhia
	server, err := sc.serverUsecase.GetServerByCompany(companyID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Retorna os servidor encontrado
	ctx.JSON(http.StatusOK, server)
}

func (sc *ServerController) GetStationsFromServer(ctx *gin.Context) {
	serverID := ctx.Param("sid")
	if serverID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "server_id is required"})
		return
	}
	server, err := sc.serverUsecase.GetServerByCompany(serverID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	url := "http://" + server.ServerIP + "/stations"

	stations, err := sc.serverUsecase.GetStationsFromServer(url)
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
