package manager

import (
	"github.com/markjforte2000/GameShelfAPI/internal/game"
	"log"
	"os"
	"path"
	"testing"
	"time"
)

type testPair struct {
	filename string
	gameName string
}

var existingLibrary = []*testPair{
	{
		filename: "Super Smash Bros Melee (2001) [GCN].rom",
		gameName: "Super Smash Bros. Melee",
	},
	{
		filename: "F-Zero GX (2003) [GCN].rom",
		gameName: "F-Zero GX",
	},
	{
		filename: "Pokemon Heartgold (2009) [NDS].rom",
		gameName: "Pokémon HeartGold",
	},
}

var gamesToAdd = []*testPair{
	{
		filename: "Pokemon Red (1996) [GBC].rom",
		gameName: "Pokémon Red",
	},
	{
		filename: "Super Mario 64 (1996) [N64].rom",
		gameName: "Super Mario 64",
	},
	{
		filename: "The Legend of Zelda: Majoras Mask (2000) [N64].rom",
		gameName: "The Legend of Zelda: Majora's Mask",
	},
}

func TestBasicGameLibManager(t *testing.T) {
	// delete existing game library
	path := os.Getenv("GAME_DIR")
	os.RemoveAll(path)
	os.Mkdir(path, os.ModePerm)
	// create initial game files
	addGames(existingLibrary)
	// init manager
	manager := NewGameLibManager()
	// check current lib
	time.Sleep(10 * time.Second)
	currentGames := manager.GetGameLibrary()
	checkGames(existingLibrary, currentGames, t)
	// add other games
	addGames(gamesToAdd)
	// give manager time to update
	time.Sleep(10 * time.Second)
	changedGames := manager.GetChanges()
	// check changes
	checkGames(gamesToAdd, changedGames, t)
}

func checkGames(expected []*testPair, actual []*game.Game, t *testing.T) {
	if len(expected) != len(actual) {
		t.Errorf("List sizes mismatch: expected: %v actual: %v", len(expected), len(actual))
	}
	for _, g := range expected {
		gameFound := false
		for _, gameInManager := range actual {
			if gameInManager.Title == g.gameName {
				gameFound = true
			}
		}
		if !gameFound {
			t.Errorf("Game %v not found in manager", g.gameName)
		}
	}
}

func addGames(games []*testPair) {
	dir := os.Getenv("GAME_DIR")
	for _, file := range games {
		f, err := os.Create(path.Join(dir, file.filename))
		if err != nil {
			log.Fatalf("Failed to create file %v: %v", file.filename, err)
		}
		f.Close()
	}
}
