package main

import (
	"event-booking-rest-api/internal/routes"
	"event-booking-rest-api/pkg/db"
	"event-booking-rest-api/pkg/models"

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
