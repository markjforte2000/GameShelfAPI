package main

import (
	"github.com/fsnotify/fsnotify"
	"github.com/markjforte2000/GameShelfAPI/internal/file"
	"log"
)

func main() {
	const rootDirectory = "./.dev/files/"
	manager := file.NewManager(rootDirectory, handler)
	manager.GetCurrentFiles()
	for {
		continue
	}
}

func handler(gameFile *file.GameFile, op fsnotify.Op) {
	log.Printf("New game file found! Name: %v\tFile: %v\n", gameFile.Name, gameFile.FileName)
}
