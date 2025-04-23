package main

import (
	"fmt"
	"math/rand"
	"slices"
)

type Station struct {
	StationID       int   `json:"station_id"`
	CoordX          int   `json:"coord_x"`
	CoordY          int   `json:"coord_y"`
	CarsWaitingList []int `json:"cars_waiting_list"`
	InUseBy         int   `json:"in_use"` // CarID
}

// Functions

/*
Retorna uma nova estação com dados aleatórios
*/
func GetNewRandomStation() Station {
	// Gerar coordenadas aleatórias para a estação
	coordX := rand.Intn(1000) // Exemplo: coordenadas entre 0 e 999
	coordY := rand.Intn(1000)

	// Criar uma nova estação com as coordenadas geradas
	return Station{
		StationID:       rand.Intn(1000), // Exemplo: ID da estação entre 0 e 999
		CoordX:          coordX,
		CoordY:          coordY,
		CarsWaitingList: []int{},
		InUseBy:         -1, // Nenhum carro em uso inicialmente
	}
}

// Retorna o índice do carro na lista de espera, ou -1 se não encontrado
func (c *Station) GetCarListIndex(carID int) int {
	index := slices.Index(c.CarsWaitingList, carID)
	return index
}

/*
Adiciona um carro à lista de espera
Retorna erro se o carro já estiver na lista
*/
func (c *Station) AddCarToList(carID int) error {
	// Verifica se o carro já está na lista de espera
	if c.GetCarListIndex(carID) != -1 {
		return fmt.Errorf("carro de ID %d já está na lista de espera", carID)
	} else {
		c.CarsWaitingList = append(c.CarsWaitingList, carID)
		return nil
	}
}

/*
Remove um carro da lista de espera
Retorna erro se o carro não estiver na lista
*/
func (c *Station) RemoveCarFromList(carID int) error {
	// Verifica se o carro já está na lista de espera
	index := c.GetCarListIndex(carID)
	if index != -1 {
		c.CarsWaitingList = slices.Delete(c.CarsWaitingList, index, index+1)
		return nil
	} else {
		return fmt.Errorf("carro de ID %d não encontrado na lista de espera", carID)
	}
}

func (c *Station) PrintState() {
	fmt.Printf("Estação ID: %d\n", c.StationID)
	fmt.Printf("Coordenadas: (%d, %d)\n", c.CoordX, c.CoordY)
	fmt.Printf("Carros na lista de espera: %v\n", c.CarsWaitingList)
	if c.InUseBy != -1 {
		fmt.Printf("Em uso pelo carro ID: %d\n", c.InUseBy)
	} else {
		fmt.Println("Nenhum carro usando a estação")
	}
}

// Getters
func (c *Station) GetStationID() int {
	return c.StationID
}

func (c *Station) GetCoordX() int {
	return c.CoordX
}

func (c *Station) GetCoordY() int {
	return c.CoordY
}

func (c *Station) GetCarsWaitingList() []int {
	return c.CarsWaitingList
}

func (c *Station) GetNumberOfCarsWaiting() int {
	return len(c.CarsWaitingList)
}

func (c *Station) GetInUseBy() int {
	return c.InUseBy
}

// Setters
func (c *Station) SetStationID(StationID int) {
	c.StationID = StationID
}

func (c *Station) SetCoordX(coordX int) {
	c.CoordX = coordX
}

func (c *Station) SetCoordY(coordY int) {
	c.CoordY = coordY
}

func (c *Station) SetCarsWaitingList(CarsWaitingList []int) {
	c.CarsWaitingList = CarsWaitingList
}

func (c *Station) SetInUseBy(inUseBy int) {
	c.InUseBy = inUseBy
}
