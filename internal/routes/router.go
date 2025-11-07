package routes

import (
	"event-booking-rest-api/internal/auth"
	eventHandler "event-booking-rest-api/internal/event"
	registrationHandler "event-booking-rest-api/internal/registration"
	userHandler "event-booking-rest-api/internal/user"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {

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
