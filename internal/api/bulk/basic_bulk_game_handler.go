package bulk

import (
	"github.com/markjforte2000/GameShelfAPI/internal/api/igdb_api"
	"github.com/markjforte2000/GameShelfAPI/internal/util"
)

const OutputBufferSize = 128

type basicBulkGameHandler struct {
	processedGames          chan *handlerResponse
	unprocessedGamesCounter *util.SafeCounter
	client                  igdb_api.AuthorizedClient
}

func (handler *basicBulkGameHandler) Add(title string, year string) {
	handler.unprocessedGamesCounter.Increment()
	go handler.asyncHandleGame(title, year)
}

func (handler *basicBulkGameHandler) asyncHandleGame(title string, year string) {
	game := handler.client.GetGameData(title, year)
	response := new(handlerResponse)
	response.GameData = game
	response.Title = title
	response.Year = year
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
