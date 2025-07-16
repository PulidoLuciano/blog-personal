package models

import "time"

type User struct {
	ID       int
	Username string
	Password string
	IsAdmin  bool
	CreatedAt time.Time
}