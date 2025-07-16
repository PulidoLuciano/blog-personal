package models

type Project struct {
	ID            int
	Title         string
	Description   string
	UrlRepository string
	UrlWebsite    string
	Image         string
	IsHidden      bool
}
