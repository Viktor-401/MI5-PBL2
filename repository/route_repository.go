package repository

import (
	"database/sql"
	//"fmt"
	//"main/model"
)

type RouteRepository struct {
	connection *sql.DB
}

//mudar a connection para o mongo

func NewRouteRepository(connection *sql.DB) RouteRepository {
	return RouteRepository{
		connection: connection,
	}
}
