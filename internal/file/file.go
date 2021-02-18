package file

import (
	"github.com/markjforte2000/GameShelfAPI/internal/game"
)

const (
	Exists Op = 0
	Create Op = 1
	Rename Op = 2
	Delete Op = 3
)

type Op uint8

// File format for roms
// Game name (year) [platform].extension

type Manager interface {
	GetCurrentFiles() []*game.GameFile
	init(rootDirectory string, handler NewFileHandler)
}

type NewFileHandler func(file *game.GameFile, fileOp Op)

func NewManager(fileDirectory string, handler NewFileHandler) Manager {
	manager := new(fsnotifyFileManager)
	manager.init(fileDirectory, handler)
	return manager
}
