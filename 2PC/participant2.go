func prepareHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RouteID    string    `json:"route_id"`
		StationID  string    `json:"station_id"`
		ReserveTime time.Time `json:"reserve_time"`
		Duration   int       `json:"duration"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	session, err := client.StartSession()
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer session.EndSession(context.Background())

	_, err = session.WithTransaction(context.Background(), func(ctx mongo.SessionContext) (interface{}, error) {
		// Check existing reservations
		endTime := req.ReserveTime.Add(time.Duration(req.Duration) * time.Minute)
		
		overlappingReservations, err := reservationsCollection.CountDocuments(ctx, bson.M{
			"station_id": req.StationID,
			"status":     "committed",
			"$or": []bson.M{
				{"start_time": bson.M{"$lt": endTime}},
				{"end_time": bson.M{"$gt": req.ReserveTime}},
			},
		})

		if err != nil {
			return nil, err
		}

		// Get station capacity
		var station Station
		if err := stationsCollection.FindOne(ctx, bson.M{"_id": req.StationID}).Decode(&station); err != nil {
			return nil, fmt.Errorf("station not found")
		}

		if int(overlappingReservations) >= station.TotalSlots {
			return nil, fmt.Errorf("no available slots")
		}

		// Create temporary reservation
		reservation := Reservation{
			RouteID:   req.RouteID,
			StationID: req.StationID,
			StartTime: req.ReserveTime,
			EndTime:   endTime,
			Status:    "prepared",
		}

		if _, err := reservationsCollection.InsertOne(ctx, reservation); err != nil {
			return nil, err
		}

		return nil, nil
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func commitHandler(w http.ResponseWriter, r *http.Request) {
	var req struct{ RouteID string }
	json.NewDecoder(r.Body).Decode(&req)

	_, err := reservationsCollection.UpdateOne(
		context.Background(),
		bson.M{"route_id": req.RouteID},
		bson.M{"$set": bson.M{"status": "committed"}},
	)

	if err != nil {
		http.Error(w, "Commit failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func abortHandler(w http.ResponseWriter, r *http.Request) {
	var req struct{ RouteID string }
	json.NewDecoder(r.Body).Decode(&req)

	_, err := reservationsCollection.DeleteOne(
		context.Background(),
		bson.M{"route_id": req.RouteID, "status": "prepared"},
	)

	if err != nil {
		http.Error(w, "Abort failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}