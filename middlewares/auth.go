package middlewares

import (
	"event-booking-rest-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authenticate(context *gin.Context) {
	tokenValue := context.Request.Header.Get("Authorization")

	if tokenValue == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	if tokenValue[:7] != "Bearer " {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	token := tokenValue[7:]

	userId, err := utils.VerifyToken(token)

	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	context.Set("userId", userId)
	context.Next()
}
