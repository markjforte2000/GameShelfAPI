package bulk

import (
	"github.com/markjforte2000/GameShelfAPI/internal/api/igdb_api"
	"github.com/markjforte2000/GameShelfAPI/internal/game"
	"github.com/markjforte2000/GameShelfAPI/internal/util"
)

const OutputBufferSize = 128

type basicBulkGameHandler struct {
	processedGames          chan *handlerResponse
	unprocessedGamesCounter *util.SafeCounter
	client                  igdb_api.AuthorizedClient
}

func (handler *basicBulkGameHandler) Add(gameFile *game.GameFile) {
	handler.unprocessedGamesCounter.Increment()
	go handler.asyncHandleGame(gameFile)
}

func (handler *basicBulkGameHandler) asyncHandleGame(gameFile *game.GameFile) {
	g := handler.client.GetGameData(gameFile)
	response := new(handlerResponse)
	response.GameData = g
	response.Title = gameFile.Title
	response.Year = gameFile.Year
	handler.processedGames <- response
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
