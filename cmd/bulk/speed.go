package main

import (
	"bufio"
	"github.com/markjforte2000/GameShelfAPI/internal/api/bulk"
	"github.com/markjforte2000/GameShelfAPI/internal/api/igdb_api"
	"github.com/markjforte2000/GameShelfAPI/internal/game"
	"log"
	"os"
	"strings"
	"time"
)

const ListFile = "gamelist.txt"

func main() {
	id, secret := getClientIDAndSecret()
	bulkHandler := bulk.NewBulkGameHandler(id, secret)
	client := igdb_api.NewAuthorizedClient(id, secret)

	inputFile1, err := os.Open(ListFile)
	inputFile2, err := os.Open(ListFile)
	if err != nil {
		log.Fatal(err)
	}
	defer inputFile1.Close()
	defer inputFile2.Close()
	bulkSpeed := testBulk(bulkHandler, inputFile1)
	sequentialSpeed := testSequential(client, inputFile2)
	log.Printf("Bulk Speed: %v\tSequential Speed: %v\n", bulkSpeed, sequentialSpeed)
}

func testSequential(client igdb_api.AuthorizedClient, file *os.File) float64 {
	start := time.Now()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ",")
		title := parts[0]
		year := parts[1]
		_ = client.GetGameData(&game.GameFile{
			Title:    title,
			Year:     year,
			Platform: "",
			FileName: "",
		})
	}
	end := time.Now()
	dur := end.Sub(start).Seconds()
	log.Printf("Sequential Method took: %v seconds", dur)
	return dur
}

func testBulk(handler bulk.BulkGameHandler, file *os.File) float64 {
	start := time.Now()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ",")
		title := parts[0]
		year := parts[1]
		handler.Add(&game.GameFile{
			Title:    title,
			Year:     year,
			Platform: "",
			FileName: "",
		})
	}

	for response := handler.Get(); response != nil; response = handler.Get() {
	}
	end := time.Now()
	dur := end.Sub(start).Seconds()
	log.Printf("Bulk Method took: %v seconds", dur)
	return dur
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
