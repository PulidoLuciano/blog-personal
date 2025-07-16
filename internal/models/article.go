package models

import "time"

type Article struct {
	Url         string
	Title       string
	Description string
	Content     string
	CreatedAt   time.Time
	Image       string
	AuthorId    int
}
