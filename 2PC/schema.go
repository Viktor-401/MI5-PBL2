type Station struct {
	ID         string     `bson:"_id"`
	Location   string     `bson:"location"`
	TotalSlots int        `bson:"total_slots"`
}

type Reservation struct {
	RouteID     string    `bson:"route_id"`
	StationID   string    `bson:"station_id"`
	UserID      string    `bson:"user_id"`
	StartTime   time.Time `bson:"start_time"`
	EndTime     time.Time `bson:"end_time"`
	Status      string    `bson:"status"` // prepared, committed, canceled
}