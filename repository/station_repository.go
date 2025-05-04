package repository

import (
    "context"
	"fmt"
	"main/model"

    "go.mongodb.org/mongo-driver/bson"
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
	fmt.Println("Creating station in repository")

	return 0, nil
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
