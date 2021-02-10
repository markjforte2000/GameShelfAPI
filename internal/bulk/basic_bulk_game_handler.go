package bulk

import (
	"github.com/markjforte2000/GameShelfAPI/internal/igdb_api"
	"log"
	"sync"
)

const OutputBufferSize = 128

type basicBulkGameHandler struct {
	processedGames          chan *handlerResponse
	unprocessedGamesCounter *safeCounter
	client                  igdb_api.AuthorizedClient
}

type safeCounter struct {
	lock  *sync.RWMutex
	count int
}

func (handler *basicBulkGameHandler) Add(title string, year string) {
	handler.unprocessedGamesCounter.increment()
	go handler.asyncHandleGame(title, year)
}

func (handler *basicBulkGameHandler) asyncHandleGame(title string, year string) {
	game := handler.client.GetGameData(title, year)
	response := new(handlerResponse)
	response.GameData = game
	response.Title = title
	response.Year = year
	handler.processedGames <- response
	log.Printf("added to channel")
}

func (handler *basicBulkGameHandler) Get() *handlerResponse {
	if handler.unprocessedGamesCounter.count <= 0 {
		return nil
	}
	handler.unprocessedGamesCounter.decrement()
	response := <-handler.processedGames
	return response
}

func (handler *basicBulkGameHandler) init(clientID string, clientSecret string) {
	handler.unprocessedGamesCounter = new(safeCounter)
	handler.client = igdb_api.NewAuthorizedClient(clientID, clientSecret)
	handler.processedGames = make(chan *handlerResponse, OutputBufferSize)
	handler.unprocessedGamesCounter.init()
}

func (counter *safeCounter) init() {
	counter.lock = new(sync.RWMutex)
	counter.count = 0
}

func (counter *safeCounter) increment() {
	counter.lock.Lock()
	counter.count++
	counter.lock.Unlock()
}

func (counter *safeCounter) decrement() {
	counter.lock.Lock()
	counter.count--
	counter.lock.Unlock()
}

func (counter *safeCounter) get() int {
	counter.lock.RLock()
	defer counter.lock.RUnlock()
	return counter.count
}
