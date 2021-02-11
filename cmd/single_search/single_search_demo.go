package main

import (
	"github.com/markjforte2000/GameShelfAPI/internal/api/igdb_api"
	"github.com/markjforte2000/GameShelfAPI/internal/util"
	"os"
	"strings"
)

func main() {
	id, secret := getClientIDAndSecret()
	client := igdb_api.NewAuthorizedClient(id, secret)
	game := client.GetGameData("Super Mario Bros", "1985")
	util.PrettyPrintGame(game)
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
