package usecase

import (
	"api/model"
	"api/repository"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ServerUsecase struct {
	serverRepo repository.ServerRepository
}

func NewServerUsecase(serverRepo repository.ServerRepository) ServerUsecase {
	return ServerUsecase{
		serverRepo: serverRepo,
	}
}

func (su *ServerUsecase) RegisterOrUpdateServer(company string, serverIP string) error {
	if company == "" || serverIP == "" {
		return fmt.Errorf("company and serverIP are required")
	}

	// Chama o repositório para registrar ou atualizar o servidor
	err := su.serverRepo.RegisterOrUpdateServer(context.Background(), company, serverIP)
	if err != nil {
		return fmt.Errorf("erro ao registrar ou atualizar servidor: %w", err)
	}

	return nil
}
func (su *ServerUsecase) GetServerByCompany(company string) (model.Server, error) {
	return su.serverRepo.GetServerByCompany(context.Background(), company)
}

// Consulta estações disponíveis em outro servidor
func (su *ServerUsecase) GetStationsFromServer(url string) ([]model.Station, error) {
	// Faz a requisição HTTP
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer requisição para o servidor remoto: %w", err)
	}
	defer resp.Body.Close()

	// Lê o corpo da resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta do servidor remoto: %w", err)
	}

	var stations []model.Station

	// Deserializa o JSON
	if err := json.Unmarshal(body, &stations); err != nil {
		return nil, fmt.Errorf("erro ao deserializar a resposta do servidor remoto: %w", err)
	}

	return stations, nil
}

// Reserva uma estação em outro servidor
func (su *ServerUsecase) ReserveStationOnServer(serverURL string, stationID int, carID int) error {
	// Cria o payload da requisição
	payload := map[string]interface{}{
		"station_id": stationID,
		"car_id":     carID,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("erro ao serializar payload: %w", err)
	}

	// Constrói a URL do endpoint remoto
	url := fmt.Sprintf("%s/stations/reserve", serverURL)

	// Faz a requisição HTTP POST
	resp, err := http.Post(url, "application/json", bytes.NewReader(payloadBytes))
	if err != nil {
		return fmt.Errorf("erro ao fazer requisição para reservar estação: %w", err)
	}
	defer resp.Body.Close()

	// Verifica o status da resposta
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body) // Lê o corpo da resposta para depuração
		return fmt.Errorf("falha ao reservar estação, status: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}
