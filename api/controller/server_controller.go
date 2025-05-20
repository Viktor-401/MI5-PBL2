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
		Company    string `json:"company"`
		ServerIP   string `json:"server_ip"`
		ServerPort string `json:"server_port"`
	}

	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := sc.serverUsecase.RegisterOrUpdateServer(request.Company, request.ServerIP, request.ServerPort)
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

	url := "http://" + server.ServerIP + ":" + server.ServerPort + "/stations"

	stations, err := sc.serverUsecase.GetStationsFromServer(url)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, stations)
}

func (sc *ServerController) PrepareStationOnServer(ctx *gin.Context) {
	// Captura o ID da estação da URL
	stationID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Captura o ID do servidor da URL
	serverID := ctx.Param("sid")
	if serverID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "server_id é obrigatório"})
		return
	}

	// Obtém o servidor pelo ID da empresa
	server, err := sc.serverUsecase.GetServerByCompany(serverID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Construa a URL do servidor para preparar a estação
	url := "http://" + server.ServerIP + ":" + server.ServerPort + "/stations/" + strconv.Itoa(stationID) + "/prepare/"

	// Captura o payload da requisição
	var request struct {
		CarID int `json:"car_id"`
	}
	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Chama o caso de uso para preparar a estação no servidor remoto
	err = sc.serverUsecase.PrepareStationOnServer(url, request.CarID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Retorna uma resposta de sucesso
	ctx.JSON(http.StatusOK, gin.H{"message": "Estação preparada com sucesso no servidor"})
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
func (sc *ServerController) CommitStationOnServer(ctx *gin.Context) {
	// Captura o ID da estação e o ID do servidor da URL
	stationID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "station_id inválido"})
		return
	}

	serverID := ctx.Param("sid")
	if serverID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "server_id é obrigatório"})
		return
	}

	// Obtém o servidor pelo ID da empresa
	server, err := sc.serverUsecase.GetServerByCompany(serverID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	url := "http://" + server.ServerIP + ":" + server.ServerPort + "/stations/" + strconv.Itoa(stationID) + "/commit/"

	// Captura o payload da requisição para obter o carID
	var request struct {
		CarID int `json:"car_id"`
	}
	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Chama o caso de uso para realizar o commit da estação no servidor
	err = sc.serverUsecase.CommitStationOnServer(url, request.CarID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Retorna uma resposta de sucesso
	ctx.JSON(http.StatusOK, gin.H{"message": "Estação comitada com sucesso no servidor"})
}

