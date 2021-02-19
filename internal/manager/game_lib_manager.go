package manager

import (
	"github.com/markjforte2000/GameShelfAPI/internal/game"
	"github.com/markjforte2000/GameShelfAPI/internal/util"
)

type GameLibManager interface {
	GetGameLibrary() []*game.Game
	GetChanges() []*game.Game
	AlterGame(g *game.Game)
	init()
}

func NewGameLibManager() GameLibManager {
	manager := new(basicGameLibManager)
	manager.init()
	return manager
}

func getGameDir() string {
	return util.GetEnvironOrFail("GAME_DIR")
}

func getDatabase() string {
	return util.GetEnvironOrFail("DATABASE")
}

func getArtworkDir() string {
	return util.GetEnvironOrFail("ARTWORK_DIR")
}

func getClientIDAndSecret() (string, string) {
	key := util.GetEnvironOrFail("CLIENT_KEY")
	value := util.GetEnvironOrFail("CLIENT_SECRET")
	return key, value
}
