package controller

import (
    "main/model"
    "main/usecase"
    "net/http"

    "github.com/gin-gonic/gin"
)

type RouteController struct {
    routeUsecase usecase.RouteUsecase
}

func NewRouteController(usecase usecase.RouteUsecase) RouteController {
    return RouteController{
        routeUsecase: usecase,
    }
}

// Endpoint para criar uma nova rota
func (rc *RouteController) CreateRoute(ctx *gin.Context) {
    route := model.Route{}

    err := ctx.BindJSON(&route)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    err = rc.routeUsecase.CreateRoute(&route)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, route)
}

// Endpoint para buscar todas as rotas com base na cidade de origem e destino final
func (rc *RouteController) GetRoutes(ctx *gin.Context) {
    startCity := ctx.Query("start_city") // Obtém a cidade de origem da query string
    endCity := ctx.Query("end_city")    // Obtém a cidade de destino da query string

    if startCity == "" || endCity == "" {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "start_city e end_city são obrigatórios"})
        return
    }

    routes, err := rc.routeUsecase.GetRoutesBetweenCities(startCity, endCity)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"routes": routes})
}



