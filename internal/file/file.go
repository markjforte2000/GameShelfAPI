package file

import "github.com/fsnotify/fsnotify"

// File format for roms
// Game name (year) [platform].extension

type GameFile struct {
	Name     string
	Year     string
	Platform string
	FileName string
}

type Manager interface {
	GetCurrentFiles() []*GameFile
	init(rootDirectory string, handler NewFileHandler)
}

type NewFileHandler func(file *GameFile, op fsnotify.Op)

func NewManager(fileDirectory string, handler NewFileHandler) Manager {
	manager := new(fsnotifyFileManager)
	manager.init(fileDirectory, handler)
	return manager
}
