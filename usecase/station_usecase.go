package usecase

import (
	"context"
	"main/model"
	"main/repository"
)

type StationUsecase struct {
	repository repository.StationRepository
}

func NewStationUseCase(repo repository.StationRepository) StationUsecase {
	return StationUsecase{
		repository: repo,
	}
}

func (su *StationUsecase) CreateStation(station model.Station) (model.Station, error) {

	id, err := su.repository.CreateStation(station)
	if err != nil {
		return model.Station{}, err
	}
	station.StationID = id

	return station, nil
}
func (su *StationUsecase) GetAllStations(ctx context.Context, company string) ([]model.Station, error) {
	stations, err := su.repository.GetAllStations(ctx, company)
	if err != nil {
		return nil, err
	}
	return stations, nil
}
