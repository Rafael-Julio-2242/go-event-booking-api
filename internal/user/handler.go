package user

import (
	"event-booking-rest-api/internal/auth"
	"event-booking-rest-api/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	dbConn *gorm.DB
}

func NewUserHandler(dbConn *gorm.DB) *UserHandler {
	return &UserHandler{dbConn: dbConn}
}

func (uh *UserHandler) Signup(context *gin.Context) {
	var user models.User

	err := context.ShouldBindBodyWithJSON(&user)

	if err != nil {
		context.JSON(http.StatusOK, gin.H{"message": "Could not parse request data."})
		return
	}

	hashedPassword, err := auth.HashPassword(user.Password)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not hash password!"})
		return
	}

	user.Password = hashedPassword

	err = user.Save(uh.dbConn)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not save user!"})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func (uh *UserHandler) Login(context *gin.Context) {
	var user models.User

	err := context.ShouldBindBodyWithJSON(&user)

	if err != nil {
		context.JSON(http.StatusOK, gin.H{"message": "Could not parse request data."})
		return
	}

	err = user.ValidateCredentials(uh.dbConn)

	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid Credentials"})
		return
	}

	token, err := auth.GenerateToken(user.Email, user.ID)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not authenticate User"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Logged in!", "token": token})
}
