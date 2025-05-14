package database

import (
	"api/model"
	"api/repository"
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// SeedData limpa e insere dados iniciais em cars e stations
func SeedData(stationRepo repository.StationRepository) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Dados iniciais para estações
	stations := []model.Station{
		{StationID: 1, InUseBy: -1, Company: "A"},
		{StationID: 5, InUseBy: -1, Company: "A"},
		{StationID: 11, InUseBy: -1, Company: "B"},
		{StationID: 2, InUseBy: -1, Company: "B"},
		{StationID: 3, InUseBy: -1, Company: "C"},
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

func SeedRoutes(db *mongo.Database) {
	routesCollection := db.Collection("routes")

	// Rotas pré-configuradas
	routes := []interface{}{
		model.Route{
			StartCity:  "A",
			EndCity:    "C",
			Waypoints:  []string{"B"}, // Rota completa A -> B -> C
			Company:    "A",
			DistanceKM: 300,
		},
		model.Route{
			StartCity:  "A",
			EndCity:    "B",
			Waypoints:  []string{"A", "B"}, // Rota direta A -> B
			Company:    "A",
			DistanceKM: 150,
		},

		// Rotas de A para D
		model.Route{
			StartCity:  "A",
			EndCity:    "D",
			Waypoints:  []string{"B", "C"}, // Rota completa A -> B -> C -> D
			Company:    "A",
			DistanceKM: 500,
		},
		model.Route{
			StartCity:  "A",
			EndCity:    "D",
			Waypoints:  []string{"C"}, // Rota alternativa A -> C -> D
			Company:    "A",
			DistanceKM: 400,
		},
		model.Route{
			StartCity:  "A",
			EndCity:    "D",
			Waypoints:  []string{"B"}, // Rota alternativa A -> B -> D
			Company:    "A",
			DistanceKM: 450,
		},
	}

	// Insere as rotas no banco de dados
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := routesCollection.InsertMany(ctx, routes)
	if err != nil {
		log.Fatalf("Erro ao inserir rotas: %v", err)
	}

	log.Println("Rotas inseridas com sucesso!")
}
