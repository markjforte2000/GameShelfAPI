package game

import "time"

type Game struct {
	ID                int
	Title             string
	ReleaseDate       time.Time
	InvolvedCompanies []*InvolvedCompany
	Summary           string
	Genres            []*Genre
	Cover             *Artwork
}

type Genre struct {
	Name string
}

type InvolvedCompany struct {
	Name      string
	ID        int
	Publisher bool
	Developer bool
}

type Artwork struct {
	URL string `json:"url"`
	ID  int    `json:"id"`
}
