package models

import "time"

type Article struct {
	Url       string
	Title     string
	Content   string
	CreatedAt time.Time
	Author    User
}