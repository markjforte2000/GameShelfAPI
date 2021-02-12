package bulk

import "github.com/markjforte2000/GameShelfAPI/internal/game"

type BulkGameHandler interface {
	Add(gameFile *game.GameFile)
	Get() *handlerResponse
	init(clientID string, clientSecret string)
}

type handlerResponse struct {
	Title    string
	Year     string
	GameData *game.Game
}

func NewBulkGameHandler(clientID string, clientSecret string) BulkGameHandler {
	handler := new(basicBulkGameHandler)
	handler.init(clientID, clientSecret)
	return handler
}
