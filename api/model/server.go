package model

import "time"

type Server struct {
	ServerIP   string    `bson:"server_ip"`   // IP do servidor
	ServerPort string    `bson:"server_port"` // Porta do servidor
	Company    string    `bson:"company"`     // Nome da empresa
	Timestamp  time.Time `bson:"timestamp"`   // Última atualização
}
