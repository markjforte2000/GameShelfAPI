package igdb_api

import (
	"fmt"
	"github.com/Henry-Sarabia/apicalypse"
	"github.com/markjforte2000/GameShelfAPI/internal/api/scheduling"
	"github.com/markjforte2000/GameShelfAPI/internal/game"
	"github.com/markjforte2000/GameShelfAPI/internal/logging"
	"github.com/markjforte2000/GameShelfAPI/internal/util"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type basicAuthClient struct {
	clientID     string
	clientSecret string
	accessToken  *token
	scheduler    scheduling.Scheduler
}

type token struct {
	accessToken string
	expiration  time.Time
	tokenType   string
}

type gameIntermediate struct {
	Name        string `json:"name"`
	ID          int    `json:"id"`
	ReleaseDate int64  `json:"first_release_date"`
	Developers  []int  `json:"involved_companies"`
	Summary     string `json:"summary"`
	CoverID     int    `json:"cover"`
	Genres      []int  `json:"genres"`
	waitGroup   *sync.WaitGroup
}

func (client *basicAuthClient) Reauthenticate() {
	newToken := getAccessToken(client.clientID, client.clientSecret)
	client.accessToken = newToken
}

func (client *basicAuthClient) IsTokenExpired() bool {
	now := time.Now()
	return now.After(client.accessToken.expiration)
}

func (client *basicAuthClient) init() {
	client.scheduler = scheduling.NewScheduler()
}

func (client *basicAuthClient) GetGameData(gameFile *game.GameFile) *game.Game {
	if client.IsTokenExpired() {
		client.Reauthenticate()
	}
	g := client.parseGameResponse(client.getGameList(gameFile))
	if g == nil {
		log.Printf("Could not find game for file: %v", gameFile.FileName)
		return nil
	}
	g.Filename = gameFile.FileName
	return g
}

func (client *basicAuthClient) getGameList(gameFile *game.GameFile) []gameIntermediate {
	request := client.constructGameRequest(gameFile)
	var gameList []gameIntermediate
	response := client.scheduler.ScheduleHTTPRequest(request, &gameList)
	response.Wait()
	if response.Error() != nil {
		log.Fatalf("Failed to do game request: %v\n", response.Error())
	}
	return gameList
}

func (client *basicAuthClient) parseGameResponse(gameList []gameIntermediate) *game.Game {
	if len(gameList) == 0 {
		return nil
	}
	topGame := gameList[0]
	topGame.waitGroup = new(sync.WaitGroup)
	return client.translateIntermediate(&topGame)
}

func (client *basicAuthClient) translateIntermediate(intermediate *gameIntermediate) *game.Game {
	g := new(game.Game)
	g.Title = intermediate.Name
	g.ID = intermediate.ID
	g.Summary = intermediate.Summary
	g.ReleaseDate = util.UnixTimestampToDate(intermediate.ReleaseDate)
	intermediate.waitGroup.Add(3)
	go client.loadGenres(intermediate, g)
	go client.loadInvolvedCompanies(intermediate, g)
	go client.loadCover(intermediate, g)
	intermediate.waitGroup.Wait()
	return g
}

func (client *basicAuthClient) loadGenres(intermediate *gameIntermediate, g *game.Game) {
	request := client.constructGenresRequest(intermediate.Genres)
	response := client.scheduler.ScheduleHTTPRequest(request, &g.Genres)
	response.Wait()
	if response.Error() != nil {
		log.Fatalf("Unable to execute genre request: %v\n", response.Error())
	}
	intermediate.waitGroup.Done()
}

func (client *basicAuthClient) constructGenresRequest(genreIDs []int) *http.Request {
	whereTerms := ""
	for _, genreID := range genreIDs {
		whereTerms += fmt.Sprintf(" | id = %v", genreID)
	}
	whereTerms = whereTerms[3:]
	request, err := apicalypse.NewRequest(
		"POST",
		"https://api.igdb.com/v4/genres",
		apicalypse.Fields("name"),
		apicalypse.Where(whereTerms),
	)
	if err != nil {
		log.Fatalf("Error creating genre request: %v\n", err)
	}
	client.addIGDBHeaders(request)
	logging.LogHTTPRequest(request)
	return request
}

func (client *basicAuthClient) loadInvolvedCompanies(intermediate *gameIntermediate, g *game.Game) {
	companies := client.getCompanyIDs(intermediate.Developers)
	client.populateInvolvedCompanyNames(companies)
	g.InvolvedCompanies = companies
	intermediate.waitGroup.Done()
}

func (client *basicAuthClient) populateInvolvedCompanyNames(companies []*game.InvolvedCompany) {
	type nameResponse struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	request := client.constructCompanyRequest(companies)
	var rawCompanies []*nameResponse
	response := client.scheduler.ScheduleHTTPRequest(request, &rawCompanies)
	response.Wait()
	if response.Error() != nil {
		log.Fatalf("Unable to execute company name request: %v\n", response.Error())
	}
	for _, rawCompany := range rawCompanies {
		for _, company := range companies {
			if rawCompany.ID == company.ID {
				company.Name = rawCompany.Name
			}
		}
	}
}

func (client *basicAuthClient) constructCompanyRequest(companies []*game.InvolvedCompany) *http.Request {
	whereTerms := ""
	for _, company := range companies {
		whereTerms += fmt.Sprintf(" | id = %v", company.ID)
	}
	whereTerms = whereTerms[2:]
	request, err := apicalypse.NewRequest(
		"POST",
		"https://api.igdb.com/v4/companies",
		apicalypse.Fields("name"),
		apicalypse.Where(whereTerms),
	)
	if err != nil {
		log.Fatalf("Error creating company name request: %v\n", err)
	}
	client.addIGDBHeaders(request)
	logging.LogHTTPRequest(request)
	return request
}

func (client *basicAuthClient) getCompanyIDs(involvedCompanyIDs []int) []*game.InvolvedCompany {
	type companyResponse struct {
		ID        int  `json:"id"`
		Company   int  `json:"company"`
		Developer bool `json:"developer"`
		Publisher bool `json:"publisher"`
	}
	request := client.constructInvolvedCompanyRequest(involvedCompanyIDs)
	var rawCompanies []*companyResponse
	response := client.scheduler.ScheduleHTTPRequest(request, &rawCompanies)
	response.Wait()
	if response.Error() != nil {
		log.Fatalf("Unable to send involved companies request: %v\n", response.Error())
	}
	var companies []*game.InvolvedCompany
	for _, rawCompany := range rawCompanies {
		company := new(game.InvolvedCompany)
		company.ID = rawCompany.Company
		company.Developer = rawCompany.Developer
		company.Publisher = rawCompany.Publisher
		companies = append(companies, company)
	}
	return companies
}

func (client *basicAuthClient) constructInvolvedCompanyRequest(involvedCompanyIDs []int) *http.Request {
	whereTerms := ""
	for _, id := range involvedCompanyIDs {
		whereTerms += fmt.Sprintf(" | id = %v", id)
	}
	whereTerms = whereTerms[2:]
	request, err := apicalypse.NewRequest(
		"POST",
		"https://api.igdb.com/v4/involved_companies",
		apicalypse.Fields("developer, company, publisher"),
		apicalypse.Where(whereTerms),
	)
	if err != nil {
		log.Fatalf("Error creating involved companies request: %v\n", err)
	}
	client.addIGDBHeaders(request)
	logging.LogHTTPRequest(request)
	return request
}

func (client *basicAuthClient) loadCover(intermediate *gameIntermediate, g *game.Game) {
	request := client.constructArtworkRequest(intermediate.CoverID)
	var covers []game.Artwork
	response := client.scheduler.ScheduleHTTPRequest(request, &covers)
	response.Wait()
	if response.Error() != nil {
		log.Fatalf("Unable to send artwork request: %v\n", response.Error())
	}
	if len(covers) == 0 {
		log.Fatalf("Unable to find artwork for game %v\n", intermediate.Name)
	}
	cover := covers[0]
	// replace thumbnail with full size artwork
	cover.RemoteURL = strings.Replace(cover.RemoteURL, "t_thumb", "t_cover_big", 1)
	cover.RemoteURL = "https:" + cover.RemoteURL
	g.Cover = &cover
	intermediate.waitGroup.Done()
}

func (client *basicAuthClient) constructArtworkRequest(artworkID int) *http.Request {
	request, err := apicalypse.NewRequest(
		"POST",
		"https://api.igdb.com/v4/covers",
		apicalypse.Limit(1),
		apicalypse.Fields("url"),
		apicalypse.Where(fmt.Sprintf("id = %v", artworkID)),
	)
	if err != nil {
		log.Fatalf("Error creating artwork request: %v\n", err)
	}
	client.addIGDBHeaders(request)
	logging.LogHTTPRequest(request)
	return request
}

func (client *basicAuthClient) constructGameRequest(gameFile *game.GameFile) *http.Request {
	startTimestamp, endTimestamp := getUnixTimestampRange(gameFile.Year)
	request, err := apicalypse.NewRequest(
		"POST",
		"https://api.igdb.com/v4/games",
		apicalypse.Limit(1),
		apicalypse.Fields("name", "genres", "first_release_date",
			"involved_companies", "summary", "cover"),
		apicalypse.Search("", gameFile.Title),
		apicalypse.Where(
			fmt.Sprintf("first_release_date > %v & first_release_date < %v & category = 0",
				startTimestamp, endTimestamp)),
	)
	client.addIGDBHeaders(request)
	if err != nil {
		log.Fatalf("Failed to create game request: %v\n", err)
	}
	logging.LogHTTPRequest(request)
	return request
}

func getUnixTimestampRange(yearString string) (int64, int64) {
	year, err := strconv.Atoi(yearString)
	if err != nil {
		log.Fatalf("Error parsing year %v: %v", yearString, err)
	}
	startYear := year
	endYear := year + 1
	return util.YearToUnixTimestamp(startYear), util.YearToUnixTimestamp(endYear)
}

func (client *basicAuthClient) addIGDBHeaders(request *http.Request) {
	bearer := fmt.Sprintf("Bearer %v", client.accessToken.accessToken)
	request.Header.Set("Client-ID", client.clientID)
	request.Header.Set("Authorization", bearer)
	request.Header.Set("Accept", "application/json")
}
