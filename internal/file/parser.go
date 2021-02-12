package file

import (
	"github.com/markjforte2000/GameShelfAPI/internal/game"
	"log"
	"regexp"
	"strings"
)

func parseFileName(name string) *game.GameFile {
	yearRegexp, err := regexp.Compile(`\([0-9]{4}\)`)
	if err != nil {
		log.Fatalf("Failed to compile year regex: %v\n", err)
	}
	platformRegexp, err := regexp.Compile(`\[.*]`)
	if err != nil {
		log.Fatalf("Failed to compile platform regex: %v\n", err)
	}
	yearString := yearRegexp.FindString(name)
	if len(yearString) == 0 {
		log.Printf("Could not find valid year for file: %v", name)
		yearString = "(null)"
	}
	platform := platformRegexp.FindString(name)
	if len(yearString) == 0 {
		log.Printf("Could not find valid platform for file: %v", name)
		platform = "[null]"
	}
	gameName := strings.Split(name, ".")[0]
	if len(gameName) == 0 {
		log.Printf("Invalid file name: %v\n", name)
		return nil
	}
	gameName = strings.Replace(gameName, yearString, "", 1)
	gameName = strings.Replace(gameName, platform, "", 1)
	gameName = gameName[0 : len(gameName)-2]
	yearString = yearString[1 : len(yearString)-1]
	platform = platform[1 : len(platform)-1]
	gameFile := &game.GameFile{
		Title:    gameName,
		Platform: platform,
		Year:     yearString,
		FileName: name,
	}
	return gameFile
}
