package model

type Station struct {
	StationID int    `bson:"station_id"`
	Company   string `bson:"company"` // Nome da empresa
	ServerIP string `bson:"server_ip"`// Servidor ao qual a estação está conectada
	InUseBy   int    `bson:"in_use"`  // CarID
	IsActive  bool   `bson:"is_active"` // Indica se a estação está ativa
}
