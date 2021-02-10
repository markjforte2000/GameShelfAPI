package igdb_api

import "github.com/markjforte2000/GameShelfAPI/internal/game"

type AuthorizedClient interface {
	GetGameData(title string, year string) *game.Game
	AsyncGetGameDate(title string, year string) (AsyncWaiter, *game.Game)
	init()
}

type AsyncWaiter interface {
	Wait()
}

func NewAuthorizedClient(clientID string, clientSecret string) AuthorizedClient {
	client := newBasicAuthClient(clientID, clientSecret)
	client.init()
	return client
}
