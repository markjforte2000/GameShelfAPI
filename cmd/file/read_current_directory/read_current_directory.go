package main

import (
	"github.com/fsnotify/fsnotify"
	"github.com/markjforte2000/GameShelfAPI/internal/file"
	"github.com/markjforte2000/GameShelfAPI/internal/game"
	"log"
)

func main() {
	const rootDirectory = "./.dev/library/"
	manager := file.NewManager(rootDirectory, handler)
	manager.GetCurrentFiles()
	for {
		continue
	}
}

func handler(gameFile *game.GameFile, op fsnotify.Op) {
	log.Printf("New game file found! Title: %v\tFile: %v\n", gameFile.Title, gameFile.FileName)
}
