package models

import "time"

type Event struct {
	ID          int
	Name        string    `binding:"required"`
	Description string    `binding:"required"`
	Location    string    `binding:"required"`
	DateTime    time.Time `binding:"required"`
	UserID      int
}

var events []Event

func (e Event) Save() error {
	// later: add it to database
	events = append(events, e)
	return nil
}

func GetAllEvents() []Event {
	return events
}
