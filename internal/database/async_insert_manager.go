package database

import (
	"github.com/markjforte2000/GameShelfAPI/internal/game"
	"sync"
)

type asyncInsertDBManager struct {
	basicManager   Manager
	pendingInserts *sync.WaitGroup
	insertLock     *sync.Mutex
}

func (manager *asyncInsertDBManager) AccessGameData(gameFile *game.GameFile) (*game.Game, bool) {
	manager.pendingInserts.Wait()
	return manager.basicManager.AccessGameData(gameFile)
}

func (manager *asyncInsertDBManager) SaveGameData(g *game.Game) {
	manager.pendingInserts.Add(1)
	go manager.asyncSaveGameData(g)
}

func (manager *asyncInsertDBManager) asyncSaveGameData(g *game.Game) {
	manager.insertLock.Lock()
	manager.basicManager.SaveGameData(g)
	manager.pendingInserts.Done()
	manager.insertLock.Unlock()
}

func (manager *asyncInsertDBManager) init(dbFile string) {
	manager.basicManager = NewManager(dbFile)
	manager.pendingInserts = new(sync.WaitGroup)
	manager.insertLock = new(sync.Mutex)
}
