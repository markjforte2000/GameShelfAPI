package main

import (
	"fmt"
	"github.com/markjforte2000/GameShelfAPI/internal/igdb_api"
	"os"
	"strings"
)

func main() {
	id, secret := getClientIDAndSecret()
	client := igdb_api.NewAuthorizedClient(id, secret)
	game := client.FindGame("Super Mario Bros", "1985")
	fmt.Printf("%+v\n", game)
	fmt.Printf("Cover: %+v\n", game.Cover)
	fmt.Printf("Involved Companies: ")
	for _, involvedCompany := range game.InvolvedCompanies {
		fmt.Printf("%+v ", involvedCompany)
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
