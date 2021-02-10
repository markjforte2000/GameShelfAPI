package main

import (
	"bufio"
	"github.com/markjforte2000/GameShelfAPI/internal/igdb_api"
	"github.com/markjforte2000/GameShelfAPI/internal/util"
	"log"
	"os"
	"strings"
)

const ListFile = "gamelist.txt"
const OutputFile = "gamelist_out.txt"

func main() {
	id, secret := getClientIDAndSecret()
	client := igdb_api.NewAuthorizedClient(id, secret)

	file, err := os.Open(ListFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	output, err := os.Create(OutputFile)

	if err != nil {
		log.Fatalf("Failed to create output file: %v\n", err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ",")
		title := parts[0]
		year := parts[1]
		game := client.GetGameData(title, year)
		gameString := util.GameToPrettyString(game)
		output.WriteString(gameString + "\n")
	}
}

func getClientIDAndSecret() (string, string) {
	var id string
	var secret string
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if pair[0] == "CLIENT-ID" {
			id = pair[1]
		} else if pair[0] == "CLIENT-SECRET" {
			secret = pair[1]
		}
	}
	return id, secret
}
