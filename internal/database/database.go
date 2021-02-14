package database

import (
	"github.com/markjforte2000/GameShelfAPI/internal/game"
)

type Manager interface {
	AccessGameDate(gameFile *game.GameFile) *game.Game
	SaveGameData(g *game.Game)
	init(dbFile string)
}

func NewManager(dbFile string) Manager {
	m := new(sqliteDBManager)
	m.init(dbFile)
	return m
}

type genreAssociation struct {
	genreID int
	gameID  int
}
