package model

type Route struct {
    ID          string   `bson:"_id,omitempty"` // ID gerado pelo MongoDB
    StartCity   string   `bson:"start_city"`    // Cidade de origem
    EndCity     string   `bson:"end_city"`      // Cidade de destino
    Waypoints   []string `bson:"waypoints"`     // Cidades intermediárias
    Company     string   `bson:"company"`       // Empresa responsável
    DistanceKM  int      `bson:"distance_km"`   // Distância total em quilômetros
}

