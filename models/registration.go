package models

import "event-booking-rest-api/db"

type Registration struct {
	ID      uint `gorm:"primaryKey;autoIncrement"`
	UserId  uint `gorm:"not null"`
	EventId uint `gorm:"not null"`
}

func (r *Registration) Save() error {
	result := db.DB.Create(r)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
