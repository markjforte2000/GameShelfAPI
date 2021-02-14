package main

import (
	"github.com/markjforte2000/GameShelfAPI/internal/database"
	"github.com/markjforte2000/GameShelfAPI/internal/game"
	"os"
	"time"
)

func main() {
	os.Remove("./.dev/database.db")
	manager := database.NewManager("./.dev/database.db")
	g := game.Game{
		ID:          1,
		Title:       "Test Game",
		ReleaseDate: time.Now(),
		InvolvedCompanies: []*game.InvolvedCompany{
			{
				Name:      "Test Company 1",
				ID:        0,
				Publisher: false,
				Developer: true,
			},
			{
				Name:      "Test Company 2",
				ID:        1,
				Publisher: true,
				Developer: false,
			},
		},
		Summary: "Test Summary",
		Genres: []*game.Genre{
			{
				Name: "Action",
				ID:   1,
			},
			{
				Name: "Adventure",
				ID:   2,
			},
		},
		Cover: &game.Artwork{
			RemoteURL: "https://www.example.com",
			ID:        1,
		},
		Filename: "test.rom",
	}
	manager.SaveGameData(&g)
	g.Filename = "update.rom"
	g.InvolvedCompanies[0].Name = "Updated name"
	manager.SaveGameData(&g)
}
