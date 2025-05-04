package main

import (
    "context"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/bson"
)

func GetCars(c *gin.Context) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    var cars []Car
    cursor, err := GetCarCollection().Find(ctx, bson.M{})
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer cursor.Close(ctx)

    cursor.All(ctx, &cars)
    c.JSON(http.StatusOK, cars)
}

func GetStations(c *gin.Context) {
    company := c.Query("company")
    filter := bson.M{}
    if company != "" {
        filter["company"] = company
    }

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    var stations []Station
    cursor, err := GetStationCollection().Find(ctx, filter)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer cursor.Close(ctx)

    cursor.All(ctx, &stations)
    c.JSON(http.StatusOK, stations)
}
