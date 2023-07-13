package entities

import (
	"time"
)

type User struct {
	ID        int64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	Email     string
	Passhash  string

	Sessions []Session
}
