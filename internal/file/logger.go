package file

import "log"

func logGameFile(file *GameFile) {
	log.Printf("Game File: Name: %v\tYear: %v\tPlatform: %v\tFile Name: %v\n",
		file.Name, file.Year, file.Platform, file.FileName)
}
