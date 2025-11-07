package models

import (
	"errors"
	"event-booking-rest-api/internal/auth"
	"event-booking-rest-api/pkg/db"

	"gorm.io/gorm"
)

type User struct {
	ID       uint   `gorm:"primarykey;autoIncrement"`
	Email    string `binding:"required" gorm:"unique;not null"`
	Password string `binding:"required" gorm:"not null"`
}

func (u *User) Save() error {
	result := db.DB.Create(u)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (u *User) ValidateCredentials() error {
	var retrievedUser User

	result := db.DB.Where("email = ?", u.Email).First(&retrievedUser)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errors.New("credentials invalid")
		}

		return result.Error
	}

	passwordIsValid := auth.CheckPasswordHash(u.Password, retrievedUser.Password)

	if !passwordIsValid {
		return errors.New("credentials invalid")
	}

	return nil
}
