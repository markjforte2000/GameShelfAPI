package game

import (
	"time"
)

type Game struct {
	ID                int                `json:"id"`
	Title             string             `json:"title"`
	ReleaseDate       time.Time          `json:"releaseDate"`
	InvolvedCompanies []*InvolvedCompany `json:"involvedCompanies"`
	Summary           string             `json:"summary"`
	Genres            []*Genre           `json:"genres"`
	Cover             *Artwork           `json:"cover"`
	Filename          string             `json:"filename"`
}

type Artwork struct {
	RemoteURL string `json:"url"`
	ID        int    `json:"id"`
}

type Genre struct {
	Name string `json:"name" db:"name"`
	ID   int    `json:"id" db:"id"`
}

type InvolvedCompany struct {
	Name      string `json:"name"`
	ID        int    `json:"id"`
	Publisher bool   `json:"publisher"`
	Developer bool   `json:"developer"`
}

type Platform struct {
	Name         string `json:"name"`
	ID           string `json:"id"`
	Abbreviation string `json:"abbreviation"`
}

type GameFile struct {
	Title    string
	Year     string
	Platform string
	FileName string
}
