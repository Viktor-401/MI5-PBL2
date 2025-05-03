package controller

import (
	"main/model"
	"main/usecase"
	"net/http"

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

func (sc *StationController) GetAllStations(ctx *gin.Context) {
	stations, err := sc.stationUsecase.GetAllStations()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, stations)
}
