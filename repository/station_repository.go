package repository

import (
	"database/sql"
	"fmt"
	"main/model"
)

type StationRepository struct {
	connection *sql.DB
}

func NewStationRepository(connection *sql.DB) StationRepository {
	return StationRepository{
		connection: connection,
	}
}

func (sr *StationRepository) CreateStation(station model.Station) (int, error) {
	fmt.Println("Creating station in repository")

	return 0, nil
}
