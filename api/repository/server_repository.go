package repository

import (
	"api/model"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ServerRepository struct {
	collection *mongo.Collection
}

func NewServerRepository(db *mongo.Database) ServerRepository {
	return ServerRepository{
		collection: db.Collection("servers"),
	}
}

// Adiciona ou atualiza um servidor no banco de dados
func (sr *ServerRepository) RegisterOrUpdateServer(ctx context.Context, company string, serverIP string) error {
	filter := bson.M{"server_ip": serverIP}
	update := bson.M{
		"$set": bson.M{
			"company":   company,
			"timestamp": time.Now(),
		},
	}
	opts := options.Update().SetUpsert(true)

	_, err := sr.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("erro ao registrar ou atualizar servidor: %w", err)
	}
	return nil
}

func (sr *ServerRepository) GetServerByCompany(ctx context.Context, company string) (model.Server, error) {
	// Define o filtro para buscar o servidor pela companhia
	filter := bson.M{"company": company}

	// Busca um único servidor no MongoDB
	var server model.Server
	err := sr.collection.FindOne(ctx, filter).Decode(&server)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return model.Server{}, fmt.Errorf("nenhum servidor encontrado para a companhia: %s", company)
		}
		return model.Server{}, fmt.Errorf("erro ao buscar servidor: %w", err)
	}

	return server, nil
}

// Obtém a lista de servidores registrados
func (sr *ServerRepository) GetRegisteredServers(ctx context.Context) ([]string, error) {
	cursor, err := sr.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar servidores registrados: %w", err)
	}
	defer cursor.Close(ctx)

	var servers []struct {
		ServerIP string `bson:"server_ip"`
	}
	if err := cursor.All(ctx, &servers); err != nil {
		return nil, fmt.Errorf("erro ao decodificar servidores: %w", err)
	}

	var ips []string
	for _, server := range servers {
		ips = append(ips, server.ServerIP)
	}
	return ips, nil
}
