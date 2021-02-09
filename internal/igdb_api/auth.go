package igdb_api

import "github.com/markjforte2000/GameShelfAPI/internal/game"

type AuthorizedClient interface {
	FindGame(title string, year string) *game.Game
}

func NewAuthorizedClient(clientID string, clientSecret string) AuthorizedClient {
	return newBasicAuthClient(clientID, clientSecret)
}
