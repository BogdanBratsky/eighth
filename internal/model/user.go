package model

import "time"

type User struct {
	ID           int
	Login        string
	Email        string
	PasswordHash string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
