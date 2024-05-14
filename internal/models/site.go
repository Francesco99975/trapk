package models

import "time"

type SEO struct {
	Description string
	Keywords    string
}
type Site struct {
	AppName  string
	Title    string
	Metatags SEO
	Year     int
}

func GetDefaultSite(title string) Site {
	return Site{
		AppName:  "trapk",
		Title:    title,
		Metatags: SEO{Description: "API", Keywords: "tool"},
		Year:     time.Now().Year(),
	}
}
