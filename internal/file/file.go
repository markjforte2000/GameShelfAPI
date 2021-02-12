package file

import (
	"github.com/fsnotify/fsnotify"
	"github.com/markjforte2000/GameShelfAPI/internal/game"
)

// File format for roms
// Game name (year) [platform].extension

type Manager interface {
	GetCurrentFiles() []*game.GameFile
	init(rootDirectory string, handler NewFileHandler)
}

type NewFileHandler func(file *game.GameFile, op fsnotify.Op)

func NewManager(fileDirectory string, handler NewFileHandler) Manager {
	manager := new(fsnotifyFileManager)
	manager.init(fileDirectory, handler)
	return manager
}
