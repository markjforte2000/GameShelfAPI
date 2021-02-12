package artwork

import (
	"fmt"
	"github.com/markjforte2000/GameShelfAPI/internal/game"
	"github.com/markjforte2000/GameShelfAPI/internal/logging"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

type basicArtworkManager struct {
	storageDirectory string
}

func (manager *basicArtworkManager) GetArtworkLocation(artwork *game.Artwork) string {
	if manager.doesArtworkExistLocally(artwork) {
		log.Printf("Artwork %v exists locally", artwork.ID)
		return manager.getArtworkFullPath(artwork)
	}
	log.Printf("Artwork %v does not exist locally", artwork.ID)
	return manager.downloadArtwork(artwork)
}

func (manager *basicArtworkManager) downloadArtwork(artwork *game.Artwork) string {
	request, err := http.NewRequest("GET", artwork.RemoteURL, nil)
	if err != nil {
		log.Fatalf("Error creating download artwork request: %v\n", err)
	}
	logging.LogHTTPRequest(request)
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatalf("Error downloading artwork: %v\n", err)
	}
	downloadLocation := manager.getArtworkFullPath(artwork)
	out, err := os.Create(downloadLocation)
	if err != nil {
		log.Fatalf("Error creating file for artwork: %v\n", err)
	}
	defer out.Close()
	_, err = io.Copy(out, response.Body)
	if err != nil {
		log.Fatalf("Error downloading artwork to file: %v\n", err)
	}
	return downloadLocation
}

func (manager *basicArtworkManager) doesArtworkExistLocally(artwork *game.Artwork) bool {
	fullArtworkPath := manager.getArtworkFullPath(artwork)
	_, err := os.Stat(fullArtworkPath)
	if err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else {
		log.Fatalf("Error checking for local artwork: %v\n", err)
	}
	return false
}

func (manager *basicArtworkManager) getArtworkFullPath(artwork *game.Artwork) string {
	return path.Join(manager.storageDirectory, getArtworkFileName(artwork))
}

func getArtworkFileName(artwork *game.Artwork) string {
	fileParts := strings.Split(artwork.RemoteURL, ".")
	if len(fileParts) == 0 || len(fileParts) == 1 {
		log.Fatalf("Invalid remote url for %+v\n", artwork)
	}
	extension := fileParts[len(fileParts)-1]
	fileName := fmt.Sprintf("%v.%v", artwork.ID, extension)
	return fileName
}
