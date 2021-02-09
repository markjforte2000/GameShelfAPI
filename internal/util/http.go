package util

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func ParseHTTPResponse(response *http.Response, output interface{}) error {
	err := json.NewDecoder(response.Body).Decode(output)
	return err
}

func PrettyPrintHTTPRequest(request *http.Request) {
	log.Printf("%v %v\n", request.Method, request.URL.Path)
	for name, headers := range request.Header {
		name = strings.ToLower(name)
		for _, header := range headers {
			log.Printf("%v: %v\n", name, header)
		}
	}
	b, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Fatalf("Unable to read http request: %v\n", err)
	}
	request.Body.Close()
	request.Body = ioutil.NopCloser(bytes.NewBuffer(b))
	log.Printf("Body: %v\n", string(b))
}

func PrettyPrintHTTPResponse(response *http.Response) {
	for name, headers := range response.Header {
		name = strings.ToLower(name)
		for _, header := range headers {
			log.Printf("%v: %v\n", name, header)
		}
	}
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Unable to read http response: %v\n", err)
	}
	response.Body.Close()
	response.Body = ioutil.NopCloser(bytes.NewBuffer(b))
	log.Printf("Body: %v\n", string(b))
}

func CreateRequestWithHeaders(url string, method string,
	headers map[string]string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	for header, value := range headers {
		request.Header.Set(header, value)
	}
	return request, nil
}
