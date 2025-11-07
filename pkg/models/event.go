package models

import (
	"errors"
	"time"

	"event-booking-rest-api/pkg/db"

	"gorm.io/gorm"
)

type Event struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	Name        string    `binding:"required" gorm:"not null"`
	Description string    `binding:"required" gorm:"not null"`
	Location    string    `binding:"required" gorm:"not null"`
	DateTime    time.Time `binding:"required" gorm:"not null" json:"dateTime"`
	UserID      uint      `gorm:"not null"`
}

func (e *Event) Save() error {
	result := db.DB.Create(e)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (e Event) Update() error {
	result := db.DB.Save(&e)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (e Event) Delete() error {
	result := db.DB.Delete(&e)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (e Event) VerifyUserRegistration(userId uint) (bool, error) {
	result := db.DB.Where("user_id = ? AND event_id = ?", userId, e.ID).First(&Registration{})

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}

		return false, result.Error
	}

	return true, nil
}

func (e Event) RegisterUser(userId uint) error {

	registration := Registration{
		UserId:  userId,
		EventId: e.ID,
	}

	return registration.Save()
}

func (e Event) UnregisterUser(userId uint) error {
	result := db.DB.Where("user_id = ? AND event_id = ?", userId, e.ID).Delete(&Registration{})
	return result.Error
}

func GetAllEvents() ([]Event, error) {
	var events []Event

	result := db.DB.Find(&events)

	if result.Error != nil {
		return nil, result.Error
	}

	return events, nil
}

func GetEventById(id int64) (*Event, error) {
	var event Event

	result := db.DB.First(&event, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &event, nil
}
