package util

import (
	"encoding/json"
	"github.com/markjforte2000/GameShelfAPI/internal/game"
	"log"
)

func PrettyPrintGame(g *game.Game) {
	out, err := json.MarshalIndent(g, "", "     ")
	if err != nil {
		log.Fatalf("Error marshaling game into json: %v\n", err)
	}
	log.Printf("Game:\n%v\n", string(out))
}
