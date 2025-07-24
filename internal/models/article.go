package models

import (
	"html/template"
	"time"
)

type Article struct {
	Url         string
	Title       string
	Description string
	Content     template.HTML
	CreatedAt   time.Time
	Image       string
	AuthorId    int
}
