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

// StationUsecase é responsável por interagir com o repositório de estações
func NewStationUseCase(repo repository.StationRepository) StationUsecase {
	return StationUsecase{
		repository: repo,
	}
}

// CreateStation cria uma nova estação no repositório
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

// CommitStation confirma a reserva de uma estação
func (su *StationUsecase) CommitStation(ctx context.Context, stationID int, carID int) error {

	err := su.ReserveStation(ctx, stationID, carID)
	if err != nil {
		return fmt.Errorf("erro ao confirmar a reserva da estação: %w", err)
	}

	return nil
}

// RemoveStation desativa uma estação no banco de dados
func (su *StationUsecase) RemoveStation(ctx context.Context, stationID int) error {
	err := su.repository.RemoveStation(ctx, stationID)
	if err != nil {
		return fmt.Errorf("erro ao remover estação: %w", err)
	}
	return nil
}

// GetAllStations busca todas as estações no repositório de um servidor
func (su *StationUsecase) GetAllStations(ctx context.Context) ([]model.Station, error) {
	stations, err := su.repository.GetAllStations(ctx)
	if err != nil {
		return nil, err
	}
	return stations, nil
}

// GetStationByID busca uma estação pelo ID
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

// PrepareStation prepara uma estação para uso (2PC/prepare)
func (su *StationUsecase) PrepareStation(ctx context.Context, stationID int, carID int) error {
	// Busca a estação pelo ID
	station, err := su.GetStationByID(ctx, stationID)
	if err != nil {
		return fmt.Errorf("erro ao buscar estação: %w", err)
	}

	if !station.IsActive {
		return fmt.Errorf("estação com ID %d não está ativa", stationID)
	}

	// Verifica se a estação já está em uso
	if station.InUseBy != -1 {
		return fmt.Errorf("estação com ID %d já está em uso", stationID)
	}

	return nil
}

// ReserveStation reserva uma estação para um carro específico
func (su *StationUsecase) ReserveStation(ctx context.Context, stationID int, carID int) error {
	// Busca a estação pelo ID
	station, err := su.GetStationByID(ctx, stationID)
	if err != nil {
		return fmt.Errorf("erro ao buscar estação: %w", err)
	}

	station.InUseBy = carID

	// Confirma a reserva
	err = su.repository.UpdateStation(ctx, station)
	if err != nil {
		return fmt.Errorf("erro ao confirmar reserva da estação: %w", err)
	}

	return nil
}

// ReleaseStation libera uma estação que estava reservada
func (su *StationUsecase) ReleaseStation(ctx context.Context, stationID int) error {
	station, err := su.GetStationByID(ctx, stationID)
	if err != nil {
		return fmt.Errorf("erro ao buscar estação %d: %w", stationID, err)
	}
	station.InUseBy = -1 // Libera a estação
	if err := su.repository.UpdateStation(ctx, station); err != nil {
		return fmt.Errorf("erro ao liberar estação %d: %w", stationID, err)
	}
	return nil
}
