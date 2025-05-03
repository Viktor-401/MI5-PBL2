package model

type Route struct {
    RouteID          int    `json:"id"`
    Origin      string `json:"origin"`
    Destination string `json:"destination"`
    Stops       []Station  `json:"stops"` // IDs das estações de recarga
}


