package igdb_api

import (
	"fmt"
	"github.com/Henry-Sarabia/apicalypse"
	"github.com/markjforte2000/GameShelfAPI/internal/game"
	"github.com/markjforte2000/GameShelfAPI/internal/util"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type basicAuthClient struct {
	clientID     string
	clientSecret string
	accessToken  *token
}

type token struct {
	accessToken string
	expiration  *time.Time
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
}

// CLIENT METHODS

func (client *basicAuthClient) FindGame(title string, year string) *game.Game {
	request := client.constructGameRequest(title, year)
	httpClient := new(http.Client)
	response, err := httpClient.Do(request)
	if err != nil {
		log.Fatalf("Failed to send game request: %v\n", err)
	}
	g := client.parseGameResponse(response)
	return g
}

func (client *basicAuthClient) parseGameResponse(response *http.Response) *game.Game {
	var gameList []gameIntermediate
	err := util.ParseHTTPResponse(response, &gameList)
	if err != nil {
		log.Fatalf("Failed to decode game search response: %v\n", err)
	}
	err = response.Body.Close()
	if err != nil {
		log.Fatalf("Error closing game request response: %v\n", err)
	}
	if len(gameList) == 0 {
		return nil
	}
	topGame := gameList[0]
	g := client.translateIntermediate(&topGame)
	return g
}

func (client *basicAuthClient) translateIntermediate(intermediate *gameIntermediate) *game.Game {
	g := game.Game{
		Title:             intermediate.Name,
		ID:                intermediate.ID,
		Summary:           intermediate.Summary,
		Cover:             client.getCover(intermediate),
		ReleaseDate:       util.UnixTimestampToDate(intermediate.ReleaseDate),
		InvolvedCompanies: client.getInvolvedCompanies(intermediate),
	}
	return &g
}

func (client *basicAuthClient) getInvolvedCompanies(intermediate *gameIntermediate) []*game.InvolvedCompany {
	companies := client.getCompanyIDs(intermediate.Developers)
	client.populateInvolvedCompanyNames(companies)
	return companies
}

func (client *basicAuthClient) populateInvolvedCompanyNames(companies []*game.InvolvedCompany) {
	type nameResponse struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	request := client.constructCompanyRequest(companies)
	httpClient := new(http.Client)
	response, err := httpClient.Do(request)
	if err != nil {
		log.Fatalf("Unable to send company name request: %v\n", err)
	}
	util.PrettyPrintHTTPResponse(response)
	var rawCompanies []*nameResponse
	err = util.ParseHTTPResponse(response, &rawCompanies)
	if err != nil {
		log.Fatalf("Error decoding company name response: %v\n", err)
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
	util.PrettyPrintHTTPRequest(request)
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
	httpClient := new(http.Client)
	response, err := httpClient.Do(request)
	if err != nil {
		log.Fatalf("Unable to send involved companies request: %v\n", err)
	}
	util.PrettyPrintHTTPResponse(response)
	var rawCompanies []*companyResponse
	err = util.ParseHTTPResponse(response, &rawCompanies)
	if err != nil {
		log.Fatalf("Error decoding involved company response: %v\n", err)
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
	util.PrettyPrintHTTPRequest(request)
	return request
}

func (client *basicAuthClient) getCover(intermediate *gameIntermediate) *game.Artwork {
	request := client.constructArtworkRequest(intermediate.CoverID)
	httpClient := new(http.Client)
	response, err := httpClient.Do(request)
	if err != nil {
		log.Fatalf("Unable to send artwork request: %v\n", err)
	}
	var covers []game.Artwork
	err = util.ParseHTTPResponse(response, &covers)
	if err != nil {
		log.Fatalf("Failed to parse artwork response: %v\n", err)
	}
	if len(covers) == 0 {
		log.Fatalf("Unable to find artwork for game %v\n", intermediate.Name)
	}
	cover := covers[0]
	// replace thumbnail with full size artwork
	cover.URL = strings.Replace(cover.URL, "t_thumb", "t_cover_big", 1)
	return &cover
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
	util.PrettyPrintHTTPRequest(request)
	return request
}

func (client *basicAuthClient) constructGameRequest(title string, year string) *http.Request {
	startTimestamp, endTimestamp := getUnixTimestampRange(year)
	request, err := apicalypse.NewRequest(
		"POST",
		"https://api.igdb.com/v4/games",
		apicalypse.Limit(1),
		apicalypse.Fields("name", "genres", "first_release_date",
			"involved_companies", "summary", "cover"),
		apicalypse.Search("", title),
		apicalypse.Where(
			fmt.Sprintf("first_release_date > %v & first_release_date < %v",
				startTimestamp, endTimestamp)),
	)
	client.addIGDBHeaders(request)
	if err != nil {
		log.Fatalf("Failed to create game request: %v\n", err)
	}
	util.PrettyPrintHTTPRequest(request)
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

func (client *basicAuthClient) constructAPIRequest(endpoint string,
	body io.Reader) (*http.Request, error) {
	bearer := fmt.Sprintf("Bearer %v", client.accessToken.accessToken)
	request, err := util.CreateRequestWithHeaders(endpoint, "POST",
		map[string]string{
			"Client-ID":     client.clientID,
			"Authorization": bearer,
			"Accept":        "application/json",
		}, body)
	return request, err
}
