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

func (rc *RouteController) CreateRoute(ctx *gin.Context) {
    
    route := model.Route{}
    
    err := ctx.BindJSON(&route);
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


