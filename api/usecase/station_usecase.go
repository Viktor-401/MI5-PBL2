package usecase

import (
	"api/model"
	"api/repository"
	"context"
	"fmt"
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

func (su *StationUsecase) RemoveStation(ctx context.Context, stationID int) error {
	err := su.repository.RemoveStation(ctx, stationID)
	if err != nil {
		return fmt.Errorf("erro ao remover estação: %w", err)
	}
	return nil
}

func (su *StationUsecase) GetAllStations(ctx context.Context) ([]model.Station, error) {
	stations, err := su.repository.GetAllStations(ctx)
	if err != nil {
		return nil, err
	}
	return stations, nil
}

func (su *StationUsecase) GetStationByID(ctx context.Context, stationID int) (model.Station, error) {
	stations, err := su.repository.GetAllStations(ctx)
	if err != nil {
		return model.Station{}, fmt.Errorf("erro ao buscar estações: %w", err)
	}

	for _, station := range stations {
		if station.StationID == stationID {
			return station, nil
		}
	}

	return model.Station{}, fmt.Errorf("estação com ID %d não encontrada", stationID)
}

func (su *StationUsecase) ReserveStation(ctx context.Context, stationID int, carID int) error {
	// Busca a estação pelo ID
	stations, err := su.repository.GetAllStations(ctx)
	if err != nil {
		return fmt.Errorf("erro ao buscar estações: %w", err)
	}

	var stationToReserve *model.Station
	for _, station := range stations {
		if station.StationID == stationID {
			stationToReserve = &station
			break
		}
	}

	if stationToReserve == nil {
		return fmt.Errorf("estação com ID %d não encontrada", stationID)
	}

	// Verifica se a estação já está em uso
	if stationToReserve.InUseBy != -1 {
		return fmt.Errorf("estação com ID %d já está em uso", stationID)
	}

	// Atualiza a estação para marcar como reservada
	stationToReserve.InUseBy = carID
	err = su.repository.UpdateStation(ctx, *stationToReserve)
	if err != nil {
		return fmt.Errorf("erro ao reservar estação: %w", err)
	}

	return nil
}
