package artwork

import (
	"github.com/markjforte2000/GameShelfAPI/internal/game"
)

type Manager interface {
	// Checks if artwork exists as local file
	// if not, manager downloads artwork from RemoteURL
	GetArtworkLocation(artwork *game.Artwork) string
}

type LocationBatch interface {
	Add(artwork *game.Artwork)
	Get() *batchResponse
	init(storageDirectory string)
}

type batchResponse struct {
	Artwork  *game.Artwork
	Location string
}

func NewManager(storageDirectory string) Manager {
	manager := new(basicArtworkManager)
	manager.storageDirectory = storageDirectory
	return manager
}

func NewBatch(storageDirectory string) LocationBatch {
	batch := new(basicLocationBatch)
	batch.init(storageDirectory)
	return batch
}
