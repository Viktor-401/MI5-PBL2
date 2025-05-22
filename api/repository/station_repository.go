package repository

import (
	"api/model"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type StationRepository struct {
	collection *mongo.Collection
}

// mudar a collection para o mongo
func NewStationRepository(db *mongo.Database) StationRepository {
	return StationRepository{
		collection: db.Collection("stations"),
	}
}


func (sr *StationRepository) CreateStation(station model.Station) (int, error) {
    // Verifica se a estação já existe
    filter := bson.M{"station_id": station.StationID}
    var existingStation model.Station
    err := sr.collection.FindOne(context.TODO(), filter).Decode(&existingStation)

    // Se a estação não existir, cria uma nova estação com todos os campos
    if err == mongo.ErrNoDocuments {
        _, err := sr.collection.InsertOne(context.TODO(), station)
        if err != nil {
            return 0, fmt.Errorf("erro ao criar nova estação: %w", err)
        }
        return station.StationID, nil
    }

    // Se a estação já existe, atualiza apenas o campo is_active
    if err != nil {
        return 0, fmt.Errorf("erro ao verificar a estação: %w", err)
    }

    update := bson.M{
        "$set": bson.M{"is_active": station.IsActive}, // Atualiza apenas o campo is_active
    }
    opts := options.Update().SetUpsert(false) // Não faz upsert porque a estação já existe

    _, err = sr.collection.UpdateOne(context.TODO(), filter, update, opts)
    if err != nil {
        return 0, fmt.Errorf("erro ao atualizar estação: %w", err)
    }
    
    return station.StationID, nil
}


// func (sr *StationRepository) CreateStation(station model.Station) (int, error) {
// 	filter := bson.M{"station_id": station.StationID}
// 	update := bson.M{
// 		"$set": station, // Atualiza todos os campos, inclusive is_active
// 	}
// 	opts := options.Update().SetUpsert(true)

// 	_, err := sr.collection.UpdateOne(context.TODO(), filter, update, opts)
// 	if err != nil {
// 		return 0, fmt.Errorf("erro ao criar/atualizar estação: %w", err)
// 	}
// 	return station.StationID, nil
// }

// func (sr *StationRepository) CreateStation(station model.Station) (int, error) {
// 	// Insere a estação na coleção
// 	_, err := sr.collection.InsertOne(context.TODO(), station)
// 	if err != nil {
// 		return 0, fmt.Errorf("erro ao criar estação: %w", err)
// 	}
// 	// Retorna o ID da estação
// 	return station.StationID, nil
// }

func (sr *StationRepository) RemoveStation(ctx context.Context, stationID int) error {
	// Define o filtro para encontrar a estação pelo ID
	filter := bson.M{"station_id": stationID}

	// Define a atualização para alterar o campo IsActive para false
	update := bson.M{
		"$set": bson.M{
			"is_active": false,
		},
	}

	// Atualiza o documento correspondente
	_, err := sr.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("erro ao desativar estação com ID %d: %w", stationID, err)
	}

	return nil
}

func (sr *StationRepository) GetAllStations(ctx context.Context) ([]model.Station, error) {
	// Define o filtro para a consulta
	filter := bson.M{}
	// Realiza a consulta no MongoDB
	cursor, err := sr.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar estações: %w", err)
	}
	defer cursor.Close(ctx)

	// Decodifica todos os documentos encontrados
	var stations []model.Station
	if err := cursor.All(ctx, &stations); err != nil {
		return nil, fmt.Errorf("erro ao decodificar estações: %w", err)
	}

	return stations, nil
}
func (sr *StationRepository) ClearStations(ctx context.Context) error {
	err := sr.collection.Drop(ctx)
	if err != nil {
		return fmt.Errorf("erro ao limpar a coleção de estações: %w", err)
	}
	return nil
}
func (sr *StationRepository) UpdateStation(ctx context.Context, station model.Station) error {
	filter := bson.M{"station_id": station.StationID}
	update := bson.M{"$set": station}
	_, err := sr.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("erro ao atualizar estação: %w", err)
	}

	return nil
}