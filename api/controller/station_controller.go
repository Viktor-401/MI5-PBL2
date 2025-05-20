package controller

import (
	"api/model"
	"api/usecase"
	"net/http"
	"strconv"
	"fmt"
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

	// Chama o caso de uso para desativar a estação
	err = sc.stationUsecase.RemoveStation(ctx, stationID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Estação desativada com sucesso"})
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
func (sc *StationController) PrepareStation(ctx *gin.Context) {
    // Captura o ID da estação da URL
    stationID, err := strconv.Atoi(ctx.Param("id"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID da estação inválido"})
        return
    }

    // Captura o payload da requisição, apenas com o CarID
    var request struct {
        CarID int `json:"CarID"`  // Apenas CarID no payload JSON
    }

    if err := ctx.BindJSON(&request); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Debug: imprime o CarID para verificar se foi deserializado corretamente
    fmt.Printf("STATION CONTROLLER PREPARESTATION - Car ID = %d\n", request.CarID)

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

func (sc *StationController) CommitStation(ctx *gin.Context) {
    // Captura o ID da estação da URL
    stationID, err := strconv.Atoi(ctx.Param("id"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
        return
    }

    // Captura o payload da requisição, que agora contém apenas o CarID
    var request struct {
        CarID int `json:"car_id"`  // Apenas CarID no payload JSON
    }

    // Faz o binding do payload JSON na estrutura
    if err := ctx.BindJSON(&request); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Debug: imprime o CarID para verificar se foi deserializado corretamente
    fmt.Printf("STATION CONTROLLER COMMIT STATION - Car ID = %d\n", request.CarID)

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



func (sc *StationController) ReserveStation(ctx *gin.Context) {
	// Captura o ID da estação da URL
	stationID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Captura o ID do carro do corpo da requisição
	var request struct {
		Car model.Car `json:"car_id"`
	}

	if err := ctx.BindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Chama o caso de uso para reservar a estação
	err = sc.stationUsecase.ReserveStation(ctx, stationID, request.Car.CarID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Estação reservada com sucesso"})
}
