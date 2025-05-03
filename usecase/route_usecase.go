package usecase

import (
	"main/model"
	"main/repository"
)

type RouteUsecase struct {
	repository repository.RouteRepository
}

func NewRouteUseCase(repo repository.RouteRepository) RouteUsecase {
	return RouteUsecase{
		repository: repo,
	}
}

func (ru *RouteUsecase) CreateRoute(route *model.Route) error {

	return nil
}
