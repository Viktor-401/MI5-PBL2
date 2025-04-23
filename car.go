package main

import (
	"math/rand"
)

type Car struct {
	CarID               int   `json:"car_id"`
	CoordX              int   `json:"coord_x"`
	CoordY              int   `json:"coord_y"`
	BatteryLevel        int   `json:"battery_level"`       // 0-100%
	BatteryDrainRate    int   `json:"battery_drain_rate"`  // % por segundo
	Speed               int   `json:"speed"`               // m/s
	RecommendedStation  int   `json:"recommended_station"` // StationID
	ReservedStation     int   `json:"reserved_station"`    // StationID
	PaidReservedStation bool  `json:"paid_reserved_station"`
	PixCode             int   `json:"pix_code"`
	CreditCardNumber    int   `json:"credit_card_number"`
	PaymentHistory      []int `json:"payment_history"` // Slice de PaymentID
}

// Functions

func GetNewRandomCar() Car {
	// Criar um novo carro com as coordenadas geradas
	return Car{
		CarID:               rand.Intn(1000), // Exemplo: ID do carro entre 0 e 999
		CoordX:              rand.Intn(1000),
		CoordY:              rand.Intn(1000),
		BatteryLevel:        100, // Bateria cheia inicialmente
		BatteryDrainRate:    1,   // % por segundo
		Speed:               20,  // m/s
		RecommendedStation:  -1,  // Nenhuma estação recomendada inicialmente
		ReservedStation:     -1,  // Nenhuma estação reservada inicialmente
		PaidReservedStation: false,
		PixCode:             rand.Intn(1000000000), // Exemplo: código Pix aleatório
		CreditCardNumber:    rand.Intn(1000000000), // Exemplo: número do cartão de crédito aleatório
	}
}

func (c *Car) PrintPaymentHistory() {
	for _, payment := range c.PaymentHistory {
		println("Payment ID:", payment)
	}
}

func (c *Car) PrintState(paymentID int) {
	println("Car ID:", c.CarID)
	println("Coordinates:", c.CoordX, c.CoordY)
	println("Battery Level:", c.BatteryLevel)
	println("Battery Drain Rate:", c.BatteryDrainRate)
	println("Speed:", c.Speed)
	println("Recommended Station:", c.RecommendedStation)
	println("Reserved Station:", c.ReservedStation)
	println("Paid Reserved Station:", c.PaidReservedStation)
	println("Pix Code:", c.PixCode)
	println("Credit Card Number:", c.CreditCardNumber)
	c.PrintPaymentHistory()
}

// Getters
func (c *Car) GetCarID() int {
	return c.CarID
}

func (c *Car) GetCoordX() int {
	return c.CoordX
}

func (c *Car) GetCoordY() int {
	return c.CoordY
}

func (c *Car) GetBatteryLevel() int {
	return c.BatteryLevel
}

func (c *Car) GetBatteryDrainRate() int {
	return c.BatteryDrainRate
}

func (c *Car) GetSpeed() int {
	return c.Speed
}

func (c *Car) GetRecommendedStation() int {
	return c.RecommendedStation
}

func (c *Car) GetReservedStation() int {
	return c.ReservedStation
}

func (c *Car) GetPaidReservedStation() bool {
	return c.PaidReservedStation
}

func (c *Car) GetPixCode() int {
	return c.PixCode
}

func (c *Car) GetCreditCardNumber() int {
	return c.CreditCardNumber
}

func (c *Car) GetPaymentHistory() []int {
	return c.PaymentHistory
}

// Setters
func (c *Car) SetCarID(carID int) {
	c.CarID = carID
}

func (c *Car) SetCoordX(coordX int) {
	c.CoordX = coordX
}

func (c *Car) SetCoordY(coordY int) {
	c.CoordY = coordY
}

func (c *Car) SetBatteryLevel(batteryLevel int) {
	c.BatteryLevel = batteryLevel
}

func (c *Car) SetBatteryDrainRate(batteryDrainRate int) {
	c.BatteryDrainRate = batteryDrainRate
}

func (c *Car) SetSpeed(speed int) {
	c.Speed = speed
}

func (c *Car) SetRecommendedStation(recommendedStation int) {
	c.RecommendedStation = recommendedStation
}

func (c *Car) SetReservedStation(reservedStation int) {
	c.ReservedStation = reservedStation
}

func (c *Car) SetPaidReservedStation(paidReservedStation bool) {
	c.PaidReservedStation = paidReservedStation
}

func (c *Car) SetPixCode(pixCode int) {
	c.PixCode = pixCode
}

func (c *Car) SetCreditCardNumber(creditCardNumber int) {
	c.CreditCardNumber = creditCardNumber
}

func (c *Car) SetPaymentHistory(paymentHistory []int) {
	c.PaymentHistory = paymentHistory
}
