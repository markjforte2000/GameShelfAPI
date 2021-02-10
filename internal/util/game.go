package util

import (
	"encoding/json"
	"github.com/markjforte2000/GameShelfAPI/internal/game"
	"log"
)

func GameToPrettyString(g *game.Game) string {
	out, err := json.MarshalIndent(g, "", "     ")
	if err != nil {
		log.Fatalf("Error marshaling game into json: %v\n", err)
	}
	return string(out)
}

func PrettyPrintGame(g *game.Game) {
	log.Printf("Game:\n%v\n", GameToPrettyString(g))
}
