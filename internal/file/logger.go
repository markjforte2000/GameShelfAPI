package file

import (
	"github.com/markjforte2000/GameShelfAPI/internal/game"
	"log"
)

func logGameFile(file *game.GameFile) {
	log.Printf("Game File: Title: %v\tYear: %v\tPlatform: %v\tFile Title: %v\n",
		file.Title, file.Year, file.Platform, file.FileName)
}
