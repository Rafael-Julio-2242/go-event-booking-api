package models

import (
	"gorm.io/gorm"
)

type Registration struct {
	ID      uint `gorm:"primaryKey;autoIncrement"`
	UserId  uint `gorm:"not null"`
	EventId uint `gorm:"not null"`
}

func (r *Registration) Save(dbConn *gorm.DB) error {
	result := dbConn.Create(r)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
