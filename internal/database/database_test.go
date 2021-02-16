package database

import (
	"github.com/markjforte2000/GameShelfAPI/internal/game"
	"path"
	"testing"
	"time"
)

const TestDir = "D:/Documents/Program/Web/GameLibrary/GameShelfAPI/.test/"

func TestManagerInterface(t *testing.T) {
	manager := NewManager(path.Join(TestDir, "interface_test.db"))
	original := newTestGame(t)
	manager.SaveGameData(original)
	accessed, exists := manager.AccessGameData(&game.GameFile{
		Title:    original.Title,
		Year:     "2021",
		Platform: "Unknown",
		FileName: original.Filename,
	})
	if !exists {
		t.Error("Unable to access game")
	}
	if !original.Equal(accessed) {
		t.Error("Accessed does not equal original")
	}
	altered := newTestGame(t)
	altered.Title = "Altered title"
	manager.SaveGameData(altered)
	accessed, exists = manager.AccessGameData(&game.GameFile{
		Title:    altered.Title,
		Year:     "2021",
		Platform: "Unknown",
		FileName: altered.Filename,
	})
	if !exists {
		t.Error("Unable to access altered game")
	}
	if original.Equal(accessed) {
		t.Error("Accessed does not equal altered")
	}
}

func newTestGame(t *testing.T) *game.Game {
	return &game.Game{
		ID:    0,
		Title: t.Name() + " Game Title",
		ReleaseDate: time.Date(2021, 2, 14, 12,
			0, 0, 0, time.UTC),
		InvolvedCompanies: []*game.InvolvedCompany{
			{
				Name:      t.Name() + " Company Publisher",
				ID:        1,
				Publisher: true,
				Developer: false,
			},
			{
				Name:      t.Name() + " Company Developer",
				ID:        2,
				Publisher: false,
				Developer: true,
			},
		},
		Summary: t.Name() + " Summary",
		Genres: []*game.Genre{
			{
				Name: t.Name() + " Genre 1",
				ID:   0,
			},
			{
				Name: t.Name() + " Genre 2",
				ID:   1,
			},
		},
		Cover: &game.Artwork{
			RemoteURL: "https://www." + t.Name() + ".lan",
			ID:        1,
		},
		Filename: t.Name() + ".rom",
	}
}
