package bulk

import (
	"github.com/markjforte2000/GameShelfAPI/internal/game"
	"sync"
)

type BulkGameHandler interface {
	Add(gameFile *game.GameFile) *waitableResponse
	Get() *handlerResponse
	init(clientID string, clientSecret string)
}

type handlerResponse struct {
	Title    string
	Year     string
	GameData *game.Game
}

type waitableResponse struct {
	lock *sync.Mutex
	g    *game.Game
}

func (resp *waitableResponse) GetGame() *game.Game {
	resp.lock.Lock()
	resp.lock.Unlock()
	return resp.g
}

func NewBulkGameHandler(clientID string, clientSecret string) BulkGameHandler {
	handler := new(basicBulkGameHandler)
	handler.init(clientID, clientSecret)
	return handler
}
