package logging

import (
	"github.com/markjforte2000/GameShelfAPI/internal/util"
	"log"
	"net/http"
)

func LogHTTPRequest(request *http.Request) {
	log.Printf("%v\t%v\t%v:\t\t\t%v\n",
		request.Method, request.URL, request.Proto, util.CopyRequestBody(request))
}

func LogHTTPResponse(originalRequest *http.Request, response *http.Response) {
	log.Printf("RESPONSE\t%v\t%v:\t\t\t%v\n",
		originalRequest.URL, response.Proto, util.FlattenJSONString(util.CopyResponseBody(response)))
}
