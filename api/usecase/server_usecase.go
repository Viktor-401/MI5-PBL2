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

// ServerUsecase é responsável por interagir com o repositório de servidores
func NewServerUsecase(serverRepo repository.ServerRepository) ServerUsecase {
	return ServerUsecase{
		serverRepo: serverRepo,
	}
}

// RegisterOrUpdateServer registra ou atualiza um servidor no banco de dados
func (su *ServerUsecase) RegisterOrUpdateServer(company string, serverIP string, port string) error {
	if company == "" || serverIP == "" || port == "" {
		return fmt.Errorf("company and serverIP are required")
	}

	// Chama o repositório para registrar ou atualizar o servidor
	err := su.serverRepo.RegisterOrUpdateServer(context.Background(), company, serverIP, port)
	if err != nil {
		return fmt.Errorf("erro ao registrar ou atualizar servidor: %w", err)
	}

	return nil
}

// GetServerByCompany busca um servidor registrado pela companhia
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

// Prepara uma estação em outro servidor (2PC/prepare)
func (su *ServerUsecase) PrepareStationOnServer(url string, carID int) error {
	// Cria o payload para a requisição
	requestBody := map[string]int{"car_id": carID}
	payload, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("erro ao marshalling o corpo da requisição: %v", err)
	}

	// Faz a requisição HTTP PUT para o servidor remoto
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("erro ao criar requisição PUT: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao enviar requisição para o servidor: %v", err)
	}
	defer resp.Body.Close()

	// Verifica o status da resposta
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("erro ao preparar estação no servidor, status code: %d", resp.StatusCode)
	}

	return nil
}

// Reserva uma estação em outro servidor usando Commit do 2PC
func (su *ServerUsecase) CommitStationOnServer(url string, carID int) error {
	// Cria o payload para a requisição
	requestBody := map[string]int{"car_id": carID}
	payload, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("erro ao marshalling o corpo da requisição: %v", err)
	}

	// Faz a requisição HTTP PUT
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("erro ao criar requisição PUT: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao enviar requisição para commit estação: %v", err)
	}
	defer resp.Body.Close()

	// Verifica o status da resposta
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body) // Lê o corpo da resposta para depuração
		return fmt.Errorf("falha ao commit estação, status: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

// Libera uma estação em outro servidor
func (su *ServerUsecase) ReleaseStationOnServer(url string, carID int) error {
	// Cria o payload para a requisição
	requestBody := map[string]int{"car_id": carID}
	payload, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("erro ao marshalling o corpo da requisição: %v", err)
	}

	// Faz a requisição HTTP PUT para o servidor remoto
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("erro ao criar requisição PUT: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao enviar requisição para liberar estação: %v", err)
	}
	defer resp.Body.Close()

	// Verifica o status da resposta
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("falha ao liberar estação, status: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}
