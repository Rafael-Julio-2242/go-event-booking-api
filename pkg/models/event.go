package models

import (
	"errors"
	"time"

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

func (e *Event) Save(dbConn *gorm.DB) error {
	result := dbConn.Create(e)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (e Event) Update(dbConn *gorm.DB) error {
	result := dbConn.Save(&e)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (e Event) Delete(dbConn *gorm.DB) error {
	result := dbConn.Delete(&e)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (e Event) VerifyUserRegistration(userId uint, dbConn *gorm.DB) (bool, error) {
	result := dbConn.Where("user_id = ? AND event_id = ?", userId, e.ID).First(&Registration{})

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}

		return false, result.Error
	}

	return true, nil
}

func (e Event) RegisterUser(userId uint, dbConn *gorm.DB) error {

	registration := Registration{
		UserId:  userId,
		EventId: e.ID,
	}

	return registration.Save(dbConn)
}

func (e Event) UnregisterUser(userId uint, dbConn *gorm.DB) error {
	result := dbConn.Where("user_id = ? AND event_id = ?", userId, e.ID).Delete(&Registration{})
	return result.Error
}

func GetAllEvents(dbConn *gorm.DB) ([]Event, error) {
	var events []Event

	result := dbConn.Find(&events)

	if result.Error != nil {
		return nil, result.Error
	}

	return events, nil
}

func GetEventById(id int64, dbConn *gorm.DB) (*Event, error) {
	var event Event

	result := dbConn.First(&event, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &event, nil
}
