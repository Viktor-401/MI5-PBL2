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

func (sc *StationController) CreateStation(ctx *gin.Context) {

	station := model.Station{}

	err := ctx.BindJSON(&station)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	station, err = sc.stationUsecase.CreateStation(station)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, station)
}

func (sc *StationController) RemoveStation(ctx *gin.Context) {
	// Captura o ID da estação da URL
	stationID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Chama o caso de uso para remover a estação
	err = sc.stationUsecase.RemoveStation(ctx, stationID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Estação removida com sucesso"})
}

func (sc *StationController) GetAllStations(ctx *gin.Context) {
	stations, err := sc.stationUsecase.GetAllStations(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, stations)
}

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

func (sc *StationController) ReserveStation(ctx *gin.Context) {
	// Captura o ID da estação da URL
	stationID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Captura o ID do carro do corpo da requisição
	var request struct {
		CarID int `json:"car_id"`
	}

	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Chama o caso de uso para reservar a estação
	err = sc.stationUsecase.ReserveStation(ctx, stationID, request.CarID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Estação reservada com sucesso"})
}
