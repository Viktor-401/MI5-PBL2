package usecase

import (
	"api/model"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ServerUsecase struct{}

func NewServerUsecase() ServerUsecase {
	return ServerUsecase{}
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

	// Estrutura auxiliar para deserializar a resposta
	var response struct {
		Stations []model.Station `json:"stations"`
	}

	// Deserializa o JSON
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("erro ao deserializar a resposta do servidor remoto: %w", err)
	}

	return response.Stations, nil
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
