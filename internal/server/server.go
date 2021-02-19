package server

import (
	"fmt"
	"github.com/markjforte2000/GameShelfAPI/internal/game"
	"github.com/markjforte2000/GameShelfAPI/internal/logging"
	"github.com/markjforte2000/GameShelfAPI/internal/manager"
	"github.com/markjforte2000/GameShelfAPI/internal/util"
	"log"
	"net/http"
	"strconv"
)

type gameLibServer struct {
	gameLibManager manager.GameLibManager
}

type gameSupplier func() []*game.Game

func ListenAndServer() {
	gameServer := gameLibServer{gameLibManager: manager.NewGameLibManager()}
	http.HandleFunc("/library", gameServer.handleLibraryRequest)
	http.HandleFunc("/changes", gameServer.handleLibraryChangeRequest)
	http.HandleFunc("/download", gameServer.handleFileDownloadRequest)
	log.Fatal(http.ListenAndServe(":5000", nil))
}

func (s *gameLibServer) handleLibraryRequest(w http.ResponseWriter, r *http.Request) {
	s.handleGameManagerRequest(w, r, s.gameLibManager.GetGameLibrary)
}

func (s *gameLibServer) handleLibraryChangeRequest(w http.ResponseWriter, r *http.Request) {
	s.handleGameManagerRequest(w, r, s.gameLibManager.GetChanges)
}

func (s *gameLibServer) handleFileDownloadRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		log.Printf("Download Request Got Unsupported Method: %v", r.Method)
		w.Write([]byte("Unsupported method"))
		return
	}
	keys, exists := r.URL.Query()["id"]
	if !exists || len(keys[0]) < 1 {
		log.Printf("Download Request Did not Specify ID")
		w.Write([]byte("No ID Specified"))
		return
	}
	id, err := strconv.Atoi(keys[0])
	if err != nil {
		log.Printf("Download request for invalid ID: %v\n", id)
		w.Write([]byte("Invalid ID"))
		return
	}
	filename, path := s.gameLibManager.GetGameFileNameAndLocation(id)
	if path == "" {
		log.Printf("Could not find game with id: %v\n", id)
		w.Write([]byte("Could not find game with specified ID"))
		return
	}
	log.Printf("Found file for id %v: %v", id, path)
	w.Header().Set("Content-Disposition",
		fmt.Sprintf(`inline; filename="%v"`, filename))
	http.ServeFile(w, r, path)
}

func (s *gameLibServer) handleGameManagerRequest(w http.ResponseWriter,
	r *http.Request, supplier gameSupplier) {
	logging.LogHTTPRequest(r)
	games := supplier()
	json := util.GameListToJSON(games)
	w.Header().Set("Content-Type", "Application/json")
	w.Write([]byte(json))
}
