package igdb_api

import (
	"fmt"
	"github.com/markjforte2000/GameShelfAPI/internal/logging"
	"github.com/markjforte2000/GameShelfAPI/internal/util"
	"log"
	"net/http"
	"time"
)

func newBasicAuthClient(clientID string, clientSecret string) *basicAuthClient {
	client := new(basicAuthClient)
	client.clientID = clientID
	client.clientSecret = clientSecret
	client.accessToken = getAccessToken(clientID, clientSecret)
	return client
}

func getAccessToken(clientID string, clientSecret string) *token {
	formattedURL := buildRequestURL(clientID, clientSecret)
	request, err := http.NewRequest("POST", formattedURL, nil)
	if err != nil {
		log.Fatalf("Failed to create token request: %v\n", request)
	}
	request.Header.Set("Content-Type", "Application/json")
	logging.LogHTTPRequest(request)
	resp, err := http.Post(formattedURL, "Application/json", nil)
	if err != nil {
		log.Fatalf("Error getting Access Token: %v", err)
	}
	logging.LogHTTPResponse(request, resp)
	accessToken := parseAccessTokenResponse(resp)
	err = resp.Body.Close()
	if err != nil {
		log.Fatalf("Failed to close token response body: %v\n", err)
	}
	return accessToken
}

func parseAccessTokenResponse(response *http.Response) *token {
	// parse json response
	type resp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		TokenType   string `json:"token_type"`
	}
	parsedResponse := resp{}
	err := util.ParseHTTPResponse(response, &parsedResponse)
	if err != nil {
		log.Fatalf("Error parsing access token response: %v\n", err)
	}
	accessToken := &token{
		tokenType:   parsedResponse.TokenType,
		accessToken: parsedResponse.AccessToken,
	}
	// set expiration time
	expiresInStr := fmt.Sprintf("%vs", parsedResponse.ExpiresIn)
	expiresIn, err := time.ParseDuration(expiresInStr)
	if err != nil {
		log.Fatalf("Unable to determine expiration time for token: %v\n", err)
	}
	expiration := time.Now().Add(expiresIn)
	accessToken.expiration = expiration
	return accessToken
}

func buildRequestURL(clientID string, clientSecret string) string {
	const requestURL = "https://id.twitch.tv/oauth2/token?client_id=%v&client_secret=%v&grant_type=client_credentials"
	return fmt.Sprintf(requestURL, clientID, clientSecret)
}
