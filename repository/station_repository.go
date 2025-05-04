package repository

import (
	"context"
	"fmt"
	"main/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
	// Insere a estação na coleção
	result, err := sr.collection.InsertOne(context.TODO(), station)
	if err != nil {
		return 0, fmt.Errorf("erro ao criar estação: %w", err)
	}

	// Converte o ID gerado para int, se possível
	id, ok := result.InsertedID.(int)
	if !ok {
		return 0, fmt.Errorf("erro ao converter o ID da estação para int")
	}

	return id, nil
}

func (sr *StationRepository) GetAllStations(ctx context.Context, company string) ([]model.Station, error) {
	filter := bson.M{}
	if company != "" {
		filter["company"] = company
	}

	cursor, err := sr.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var stations []model.Station

	err = cursor.All(ctx, &stations)
	if err != nil {
		return nil, err
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
