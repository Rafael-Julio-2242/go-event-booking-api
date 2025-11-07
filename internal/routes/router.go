package routes

import (
	"event-booking-rest-api/internal/auth"
	"event-booking-rest-api/internal/event"
	"event-booking-rest-api/internal/registration"
	"event-booking-rest-api/internal/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(server *gin.Engine, dbConn *gorm.DB) {

	userHandler := user.NewUserHandler(dbConn)
	eventHandler := event.NewEventHandler(dbConn)
	registrationHandler := registration.NewRegistrationHandler(dbConn)

	server.GET("/events", eventHandler.GetEvents)
	server.GET("/events/:id", eventHandler.GetEvent)

	authenticated := server.Group("/")
	authenticated.Use(auth.Authenticate)
	authenticated.POST("/events", eventHandler.CreateEvent)
	authenticated.PUT("/events/:id", eventHandler.UpdateEvent)
	authenticated.DELETE("/events/:id", eventHandler.DeleteEvent)
	authenticated.POST("/events/:id/register", registrationHandler.RegisterForEvent)
	authenticated.DELETE("/events/:id/register", registrationHandler.CancelRegistration)

	server.POST("/signup", userHandler.Signup)
	server.POST("/login", userHandler.Login)
}
