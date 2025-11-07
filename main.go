package main

import (
	"event-booking-rest-api/db"
	"event-booking-rest-api/models"
	"event-booking-rest-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	db.DB.AutoMigrate(
		&models.User{},
		&models.Event{},
		&models.Registration{},
	)

	server := gin.Default()

	routes.RegisterRoutes(server)

	server.Run(":8080")
}
