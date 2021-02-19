package manager

import (
	"github.com/markjforte2000/GameShelfAPI/internal/api/bulk"
	"github.com/markjforte2000/GameShelfAPI/internal/artwork"
	"github.com/markjforte2000/GameShelfAPI/internal/database"
	"github.com/markjforte2000/GameShelfAPI/internal/file"
	"github.com/markjforte2000/GameShelfAPI/internal/game"
	"sync"
)

type basicGameLibManager struct {
	gameBulk       bulk.BulkGameHandler
	artworkManager artwork.Manager
	database       database.Manager
	fileManager    file.Manager
	library        []*game.Game
	changedGames   []*game.Game
}

func (manager *basicGameLibManager) GetGameLibrary() []*game.Game {
	manager.changedGames = []*game.Game{}
	return manager.library
}

func (manager *basicGameLibManager) GetChanges() []*game.Game {
	changes := make([]*game.Game, len(manager.changedGames))
	copy(changes, manager.changedGames)
	manager.changedGames = []*game.Game{}
	return changes
}

func (manager *basicGameLibManager) AlterGame(g *game.Game) {
	manager.database.SaveGameData(g)
}

func (manager *basicGameLibManager) init() {
	clientID, clientSecret := getClientIDAndSecret()
	manager.gameBulk = bulk.NewBulkGameHandler(clientID, clientSecret)
	manager.artworkManager = artwork.NewManager(getArtworkDir())
	manager.database = database.NewAsyncInsertManager(getDatabase())
	manager.fileManager = file.NewManager(getGameDir(), manager.handleNewFile)
	manager.library = []*game.Game{}
	manager.changedGames = []*game.Game{}
	// init current games
	fileWaitGroup := new(sync.WaitGroup)
	for _, f := range manager.fileManager.GetCurrentFiles() {
		fileWaitGroup.Add(1)
		go manager.addExistingGame(f, fileWaitGroup)
	}
	fileWaitGroup.Wait()
}

func (manager *basicGameLibManager) addExistingGame(f *game.GameFile,
	fileWaitGroup *sync.WaitGroup) {
	manager.handleNewFile(f, file.Exists)
	fileWaitGroup.Done()
}

func (manager *basicGameLibManager) handleNewFile(f *game.GameFile, op file.Op) {
	if op == file.Rename || op == file.Delete {
		// TODO add rename and delete
		return
	}
	// check if file is in database
	g, exists := manager.database.AccessGameData(f)
	if exists {
		manager.addNewGame(g)
		return
	}
	// get file from bulk
	response := manager.gameBulk.Add(f)
	g = response.GetGame()
	// add file to database
	manager.database.SaveGameData(g)
	manager.addNewGame(g)
}

func (manager *basicGameLibManager) addNewGame(g *game.Game) {
	manager.changedGames = append(manager.changedGames, g)
	manager.library = append(manager.library, g)
}
