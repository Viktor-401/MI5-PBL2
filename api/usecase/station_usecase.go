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
// Garante que o campo IsActive seja true ao criar a estação
station.IsActive = true

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

func (su *StationUsecase) PrepareStation(ctx context.Context, stationID int, carID int) error {
	// Busca a estação pelo ID
	station, err := su.GetStationByID(ctx, stationID)
	if err != nil {
		return fmt.Errorf("erro ao buscar estação: %w", err)
	}

	// Verifica se a estação já está em uso
	if station.InUseBy != -1 {
		return fmt.Errorf("estação com ID %d já está em uso", stationID)
	}

	// Atualiza a estação para marcar como "preparada" (reservada temporariamente)
	station.InUseBy = carID
	err = su.repository.UpdateStation(ctx, station)
	if err != nil {
		return fmt.Errorf("erro ao preparar estação: %w", err)
	}

	return nil
}
func (su *StationUsecase) ReserveStation(ctx context.Context, stationID int, carID int) error {
	// Busca a estação pelo ID
	station, err := su.GetStationByID(ctx, stationID)
	if err != nil {
		return fmt.Errorf("erro ao buscar estação: %w", err)
	}

	// Verifica se a estação já está reservada por outro carro
	if station.InUseBy != carID {
		return fmt.Errorf("estação com ID %d não está reservada pelo carro %d", stationID, carID)
	}

	// Confirma a reserva (pode incluir lógica adicional, se necessário)
	err = su.repository.UpdateStation(ctx, station)
	if err != nil {
		return fmt.Errorf("erro ao confirmar reserva da estação: %w", err)
	}

	return nil
}
