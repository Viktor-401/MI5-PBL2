package controller

import (
	"api/model"
	"api/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type StationController struct {
	stationUsecase usecase.StationUsecase
}

func NewStationController(usecase usecase.StationUsecase) StationController {
	return StationController{
		stationUsecase: usecase,
	}
}

// Handler para criar uma nova estação
func (sc *StationController) CreateStation(ctx *gin.Context) {

	station := model.Station{}

	err := ctx.BindJSON(&station)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Chama o caso de uso para criar a estação
	station, err = sc.stationUsecase.CreateStation(station)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, station)
}

// Handler para desativar uma estação
func (sc *StationController) RemoveStation(ctx *gin.Context) {
	// Captura o ID da estação da URL
	stationID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Chama o caso de uso para desativar a estação
	err = sc.stationUsecase.RemoveStation(ctx, stationID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Estação desativada com sucesso"})
}

// Handler para obter todas as estações de um servidor
func (sc *StationController) GetAllStations(ctx *gin.Context) {
	stations, err := sc.stationUsecase.GetAllStations(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, stations)
}

// Handler para obter uma estação específica pelo ID
func (sc *StationController) GetStationByID(ctx *gin.Context) {
	stationID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	station, err := sc.stationUsecase.GetStationByID(ctx, stationID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, station)
}

// Handler para preparar uma estação (2PC/prepare)
func (sc *StationController) PrepareStation(ctx *gin.Context) {
	// Captura o ID da estação da URL
	stationID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID da estação inválido"})
		return
	}

	// Captura o payload da requisição, apenas com o CarID
	var request struct {
		CarID int `json:"car_id"` // Apenas CarID no payload JSON
	}

	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verifica se CarID foi corretamente deserializado
	if request.CarID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "car_id inválido"})
		return
	}

	// Chama o caso de uso para preparar a estação
	err = sc.stationUsecase.PrepareStation(ctx, stationID, request.CarID)
	if err != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Estação preparada com sucesso"})
}

// Handler para confirmar uma estação (2PC/commit)
func (sc *StationController) CommitStation(ctx *gin.Context) {
	// Captura o ID da estação da URL
	stationID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Captura o payload da requisição, que agora contém apenas o CarID
	var request struct {
		CarID int `json:"car_id"` // Apenas CarID no payload JSON
	}

	// Faz o binding do payload JSON na estrutura
	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verifica se o CarID foi corretamente deserializado
	if request.CarID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "car_id inválido"})
		return
	}

	// Chama a função de "commit" que efetiva a reserva da estação
	err = sc.stationUsecase.CommitStation(ctx, stationID, request.CarID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Estação confirmada e reservada com sucesso"})
}

// Handler para liberar uma estação apos o fim da viagem
func (sc *StationController) ReleaseStation(ctx *gin.Context) {
	stationID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := sc.stationUsecase.ReleaseStation(ctx, stationID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Estação liberada com sucesso"})
}
