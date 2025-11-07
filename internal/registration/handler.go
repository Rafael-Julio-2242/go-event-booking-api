package registration

import (
	"event-booking-rest-api/pkg/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RegistrationHandler struct {
	dbConn *gorm.DB
}

func NewRegistrationHandler(dbConn *gorm.DB) *RegistrationHandler {
	return &RegistrationHandler{dbConn: dbConn}
}

func (rh *RegistrationHandler) RegisterForEvent(context *gin.Context) {
	id := context.Param("id")

	if id == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Event ID not informed!"})
		return
	}

	userId := context.GetUint("userId")

	eventId, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to parse information"})
		return
	}

	event, err := models.GetEventById(eventId, rh.dbConn)

	if event == nil && err == nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Event not found"})
		return
	}

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch Event"})
		return
	}

	isRegistered, err := event.VerifyUserRegistration(userId, rh.dbConn)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to verify user registration"})
		return
	}

	if isRegistered {
		context.JSON(http.StatusBadRequest, gin.H{"message": "User already registered in event"})
		return
	}

	err = event.RegisterUser(userId, rh.dbConn)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to register user"})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "User Registered successfully"})
}

func (rh *RegistrationHandler) CancelRegistration(context *gin.Context) {
	id := context.Param("id")

	if id == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Event ID not informed!"})
		return
	}

	userId := context.GetUint("userId")

	eventId, err := strconv.ParseInt(id, 10, 64)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to parse information"})
		return
	}

	event, err := models.GetEventById(eventId, rh.dbConn)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch Event"})
		return
	}

	isRegistered, err := event.VerifyUserRegistration(userId, rh.dbConn)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to verify user registration"})
		return
	}

	if !isRegistered {
		context.JSON(http.StatusBadRequest, gin.H{"message": "User is not Registered in the event!"})
		return
	}

	err = event.UnregisterUser(userId, rh.dbConn)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to register user"})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "User Unregistered successfully"})
}
