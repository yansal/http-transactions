package model

import "time"

type User struct {
	ID    int64
	Email string
}

type UserAction struct {
	ID         int64
	UserID     int64
	Action     string
	OccurredAt time.Time
}
