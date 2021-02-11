package util

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

func ParseHTTPResponse(response *http.Response, output interface{}) error {
	err := json.NewDecoder(response.Body).Decode(output)
	return err
}

func CopyRequestBody(request *http.Request) string {
	body, content := copyBody(request.Body)
	request.Body = body
	return content
}

func CopyResponseBody(response *http.Response) string {
	body, content := copyBody(response.Body)
	response.Body = body
	return content
}

func copyBody(body io.ReadCloser) (io.ReadCloser, string) {
	if body == nil {
		return nil, ""
	}
	b, err := ioutil.ReadAll(body)
	if err != nil {
		log.Fatalf("Unable to read http body: %v\n", err)
	}
	body.Close()
	return ioutil.NopCloser(bytes.NewBuffer(b)), string(b)
}
