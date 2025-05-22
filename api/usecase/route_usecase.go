package usecase

import (
	"api/model"
	"api/repository"
)

type RouteUsecase struct {
	routeRepo repository.RouteRepository
}

// RouteUsecase é responsável por interagir com o repositório de rotas
func NewRouteUsecase(routeRepo *repository.RouteRepository) RouteUsecase {
	return RouteUsecase{
		routeRepo: *routeRepo, // Desreferencia o ponteiro
	}
}

// CreateRoute cria uma nova rota no repositório
func (ru *RouteUsecase) CreateRoute(route *model.Route) error {
	return ru.routeRepo.CreateRoute(route)
}

// Busca todas as rotas entre duas cidades usando o repositório
func (ru *RouteUsecase) GetRoutesBetweenCities(startCity, endCity string) ([]model.Route, error) {
	return ru.routeRepo.GetRoutesBetweenCities(startCity, endCity)
}
