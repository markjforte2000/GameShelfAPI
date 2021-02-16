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

func (g *Game) Equal(other *Game) bool {
	if g.ID != other.ID || g.Filename != other.Filename ||
		!g.ReleaseDate.Truncate(time.Second).Equal(other.ReleaseDate.Truncate(time.Second)) ||
		g.Summary != other.Summary || g.Title != other.Title {
		return false
	}
	if !g.Cover.Equal(other.Cover) {
		return false
	}
	if len(g.Genres) != len(other.Genres) {
		return false
	}
	for _, genre := range g.Genres {
		match := false
		for _, otherGenre := range other.Genres {
			if genre.Equal(otherGenre) {
				match = true
			}
		}
		if !match {
			return false
		}
	}
	if len(g.InvolvedCompanies) != len(other.InvolvedCompanies) {
		return false
	}
	for _, company := range g.InvolvedCompanies {
		match := false
		for _, otherCompany := range other.InvolvedCompanies {
			if company.Equal(otherCompany) {
				match = true
			}
		}
		if !match {
			return false
		}
	}
	return true
}

type Artwork struct {
	RemoteURL string `json:"url"`
	ID        int    `json:"id"`
}

func (art *Artwork) Equal(other *Artwork) bool {
	return art.ID == other.ID && art.RemoteURL == other.RemoteURL
}

type Genre struct {
	Name string `json:"name" db:"name"`
	ID   int    `json:"id" db:"id"`
}

func (genre *Genre) Equal(other *Genre) bool {
	return genre.ID == other.ID && genre.Name == other.Name
}

type InvolvedCompany struct {
	Name      string `json:"name"`
	ID        int    `json:"id"`
	Publisher bool   `json:"publisher"`
	Developer bool   `json:"developer"`
}

func (company *InvolvedCompany) Equal(other *InvolvedCompany) bool {
	return company.Developer == other.Developer && company.Publisher == other.Publisher &&
		company.ID == other.ID && company.Name == other.Name
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
