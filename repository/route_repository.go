package repository

import (
	"database/sql"
	//"fmt"
	//"main/model"
)

type RouteRepository struct {
	connection *sql.DB
}

func NewRouteRepository(connection *sql.DB) RouteRepository {
	return RouteRepository{
		connection: connection,
	}
}
