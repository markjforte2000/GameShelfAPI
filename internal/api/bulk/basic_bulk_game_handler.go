package bulk

import (
	"github.com/markjforte2000/GameShelfAPI/internal/api/igdb_api"
	"github.com/markjforte2000/GameShelfAPI/internal/game"
	"github.com/markjforte2000/GameShelfAPI/internal/util"
	"sync"
)

const OutputBufferSize = 128

type basicBulkGameHandler struct {
	processedGames          chan *handlerResponse
	unprocessedGamesCounter *util.SafeCounter
	client                  igdb_api.AuthorizedClient
}

func (handler *basicBulkGameHandler) Add(gameFile *game.GameFile) *waitableResponse {
	handler.unprocessedGamesCounter.Increment()
	response := &waitableResponse{
		lock: new(sync.Mutex),
	}
	response.lock.Lock()
	go handler.asyncHandleGame(gameFile, response)
	return response
}

func (handler *basicBulkGameHandler) asyncHandleGame(gameFile *game.GameFile,
	waitable *waitableResponse) {
	g := handler.client.GetGameData(gameFile)
	if g == nil {
		waitable.lock.Unlock()
		return
	}
	response := new(handlerResponse)
	response.GameData = g
	response.Title = gameFile.Title
	response.Year = gameFile.Year
	handler.processedGames <- response
	waitable.g = g
	waitable.lock.Unlock()
}

func (handler *basicBulkGameHandler) Get() *handlerResponse {
	if handler.unprocessedGamesCounter.Get() <= 0 {
		return nil
	}
	handler.unprocessedGamesCounter.Decrement()
	response := <-handler.processedGames
	return response
}

func (handler *basicBulkGameHandler) init(clientID string, clientSecret string) {
	handler.client = igdb_api.NewAuthorizedClient(clientID, clientSecret)
	handler.processedGames = make(chan *handlerResponse, OutputBufferSize)
	handler.unprocessedGamesCounter = util.NewSafeCounter()
}
