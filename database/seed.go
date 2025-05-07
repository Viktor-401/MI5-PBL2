package database

import (
	"context"
	"fmt"
	"main/model"
	"main/repository"
	"time"
)

// SeedData limpa e insere dados iniciais em cars e stations
func SeedData(stationRepo repository.StationRepository) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Dados iniciais para estações
	stations := []model.Station{
		// {StationID: 1, CoordX: 2, CoordY: 2, InUseBy: 0, Company: "A"},
		// {StationID: 10, CoordX: 50, CoordY: 100, InUseBy: 0, Company: "A"},
		// {StationID: 11, CoordX: 25, CoordY: 50, InUseBy: 3, Company: "B"},
		// {StationID: 2, CoordX: 3, CoordY: 3, InUseBy: 0, Company: "B"},
		// {StationID: 3, CoordX: 4, CoordY: 4, InUseBy: 0, Company: "C"},
	}

	// Limpa a coleção de estações
	err := stationRepo.ClearStations(ctx)
	if err != nil {
		fmt.Printf("Erro ao limpar a coleção de estações: %v\n", err)
		return
	}

	// Insere as estações
	for _, station := range stations {
		_, err := stationRepo.CreateStation(station)
		if err != nil {
			fmt.Printf("Erro ao inserir estação %d: %v\n", station.StationID, err)
			return
		}
	}

	fmt.Println("✅ Banco populado com sucesso")
}
