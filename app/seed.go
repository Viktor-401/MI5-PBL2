package main

import (
    "context"
    "fmt"
    "time"
)

// Car e Station são os modelos usados no seed

type Car struct {
    CarID               int      `bson:"car_id"`
    User                string   `bson:"user"`
    Password            string   `bson:"password"`
    CoordX              int      `bson:"coord_x"`
    CoordY              int      `bson:"coord_y"`
    BatteryLevel        int      `bson:"battery_level"`
    RecomendedStation   int      `bson:"recomended_station"`
    ReservedStation     int      `bson:"reserved_station"`
    PaidReservedStation bool     `bson:"paid_reserved_station"`
    PixCode             int      `bson:"pix_code"`
    CreditCardNumber    int      `bson:"credit_card_number"`
    PaymentHistory      []string `bson:"payment_history"`
}

type Station struct {
    StationID   int    `bson:"station_id"`
    CoordX      int    `bson:"coord_x"`
    CoordY      int    `bson:"coord_y"`
    CarList     []int  `bson:"car_list"`
    CarsWaiting int    `bson:"cars_waiting"`
    InUse       int    `bson:"in_use"`
    Company     string `bson:"company"` // "A", "B" ou "C"
}

// SeedData limpa e insere dados iniciais em cars e stations
func SeedData() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    cars := []interface{}{
        Car{401, "401", "401", 50, 100, 100, 0, 0, false, 401, 0, []string{}},
        Car{7, "Alisson", "als", 3, 10, 100, 0, 0, false, 0, 7, nil},
        Car{3, "Dk", "dk", 25, 50, 16, 0, 0, false, 7, 0, nil},
    }

    stations := []interface{}{
        Station{1, 2, 2, []int{}, 0, 0, "A"},
        Station{10, 50, 100, []int{}, 0, 0, "A"},
        Station{11, 25, 50, []int{}, 0, 3, "B"},
        Station{2, 3, 3, nil, 0, 0, "B"},
        Station{3, 4, 4, nil, 0, 0, "C"},
    }

    // limpa coleções
    GetCarCollection().Drop(ctx)
    GetStationCollection().Drop(ctx)

    // insere
    if _, err := GetCarCollection().InsertMany(ctx, cars); err != nil {
        panic(err)
    }
    if _, err := GetStationCollection().InsertMany(ctx, stations); err != nil {
        panic(err)
    }

    fmt.Println("✅ Banco populado com sucesso")
}
