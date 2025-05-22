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

// StationRepository é responsável por interagir com a coleção de estações no MongoDB
func NewStationRepository(db *mongo.Database) StationRepository {
	return StationRepository{
		collection: db.Collection("stations"),
	}
}

// CreateStation insere uma nova estação no banco de dados ou atualiza o campo is_active se a estação já existir
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

// RemoveStation desativa uma estação no banco de dados definindo o campo is_active como false
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

// Retorna todas as estações do banco de dados de um servidor
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

// Atualiza uma estação no banco de dados com o objeto station
func (sr *StationRepository) UpdateStation(ctx context.Context, station model.Station) error {
	filter := bson.M{"station_id": station.StationID}
	update := bson.M{"$set": station}
	_, err := sr.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("erro ao atualizar estação: %w", err)
	}

	return nil
}
