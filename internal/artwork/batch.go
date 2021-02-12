package artwork

import (
	"github.com/markjforte2000/GameShelfAPI/internal/game"
	"github.com/markjforte2000/GameShelfAPI/internal/util"
)

const ArtworkQueueSize = 128

type basicLocationBatch struct {
	processedArtwork        chan *batchResponse
	unprocessedArtworkCount *util.SafeCounter
	manager                 Manager
}

func (batch *basicLocationBatch) init(storageDirectory string) {
	batch.processedArtwork = make(chan *batchResponse, ArtworkQueueSize)
	batch.unprocessedArtworkCount = util.NewSafeCounter()
	batch.manager = NewManager(storageDirectory)
}

func (batch *basicLocationBatch) Add(artwork *game.Artwork) {
	batch.unprocessedArtworkCount.Increment()
	go batch.asyncProcessArtwork(artwork)
}

func (batch *basicLocationBatch) asyncProcessArtwork(art *game.Artwork) {
	location := batch.manager.GetArtworkLocation(art)
	response := new(batchResponse)
	response.Artwork = art
	response.Location = location
	batch.processedArtwork <- response
	batch.unprocessedArtworkCount.Decrement()
}

func (batch *basicLocationBatch) Get() *batchResponse {
	if batch.unprocessedArtworkCount.Get() <= 0 {
		return nil
	}
	return <-batch.processedArtwork
}
