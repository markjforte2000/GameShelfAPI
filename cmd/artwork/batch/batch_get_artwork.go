package main

import (
	"github.com/markjforte2000/GameShelfAPI/internal/artwork"
	"github.com/markjforte2000/GameShelfAPI/internal/game"
	"log"
	"os"
)

func main() {
	artworks := []*game.Artwork{
		&game.Artwork{
			RemoteURL: "https://images.igdb.com/igdb/image/upload/t_cover_big/co29pl.jpg",
			ID:        105897,
		},
		&game.Artwork{
			RemoteURL: "https://images.igdb.com/igdb/image/upload/t_cover_big/co1voj.jpg",
			ID:        87715,
		},
		&game.Artwork{
			RemoteURL: "https://images.igdb.com/igdb/image/upload/t_cover_big/co2362.jpg",
			ID:        97418,
		},
		&game.Artwork{
			RemoteURL: "https://images.igdb.com/igdb/image/upload/t_cover_big/co1vcp.jpg",
			ID:        87289,
		},
		&game.Artwork{
			RemoteURL: "https://images.igdb.com/igdb/image/upload/t_cover_big/co1uii.jpg",
			ID:        86202,
		},
		&game.Artwork{
			RemoteURL: "https://images.igdb.com/igdb/image/upload/t_cover_big/co2erg.jpg",
			ID:        112444,
		},
		&game.Artwork{
			RemoteURL: "https://images.igdb.com/igdb/image/upload/t_cover_big/co2s5t.jpg",
			ID:        129809,
		},
		&game.Artwork{
			RemoteURL: "https://images.igdb.com/igdb/image/upload/t_cover_big/co21yv.jpg",
			ID:        95863,
		},
	}
	const storageDirectory = ".dev/artwork"
	err := os.RemoveAll(storageDirectory)
	if err != nil {
		log.Fatalf("Error deleting directory contents: %v\n", err)
	}
	os.Mkdir(storageDirectory, os.ModeDir)
	batch := artwork.NewBatch(storageDirectory)
	for _, art := range artworks {
		batch.Add(art)
	}
	for resp := batch.Get(); resp != nil; resp = batch.Get() {
		log.Printf("Location for %v is %v\n", resp.Artwork.ID, resp.Location)
	}
}
