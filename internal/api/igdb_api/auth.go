package igdb_api

import "github.com/markjforte2000/GameShelfAPI/internal/game"

type AuthorizedClient interface {
	GetGameData(gameFile *game.GameFile) *game.Game
	Reauthenticate()
	IsTokenExpired() bool
	init()
}

func NewAuthorizedClient(clientID string, clientSecret string) AuthorizedClient {
	client := newBasicAuthClient(clientID, clientSecret)
	client.init()
	return client
}
