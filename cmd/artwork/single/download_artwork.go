package main

import (
	"github.com/markjforte2000/GameShelfAPI/internal/artwork"
	"github.com/markjforte2000/GameShelfAPI/internal/game"
	"log"
	"os"
	"path"
)

func main() {
	art := game.Artwork{
		RemoteURL: "https://images.igdb.com/igdb/image/upload/t_cover_big/co2362.jpg",
		ID:        555,
	}
	const storageDirectory = ".dev/artwork"
	// delete existing file
	if _, err := os.Stat(path.Join(storageDirectory, "555.jpg")); err == nil {
		err = os.Remove(path.Join(storageDirectory, "555.jpg"))
		if err != nil {
			log.Fatalf("Failed to delete existing file: %v\n", err)
		}
	}
	manager := artwork.NewManager(storageDirectory)
	log.Printf("Location: %v\n", manager.GetArtworkLocation(&art))
	// check now that art exists
	log.Printf("Location: %v\n", manager.GetArtworkLocation(&art))
}
