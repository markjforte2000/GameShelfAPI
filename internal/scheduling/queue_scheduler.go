package scheduling

import (
	"github.com/markjforte2000/GameShelfAPI/internal/logging"
	"github.com/markjforte2000/GameShelfAPI/internal/util"
	"net/http"
	"sync"
	"time"
)

const MaxRequestsInQueue = 128
const MaxRequestsPerSecond = 4

type queueScheduler struct {
	queue              chan *queueRequest
	recentRequests     []*queueRequest
	recentRequestsLock *sync.RWMutex
	httpClient         *http.Client
}

type queueRequest struct {
	executionTime time.Time
	releaseTime   time.Time
	lock          *sync.Mutex
	httpRequest   *http.Request
	output        interface{}
	done          bool
	err           error
}

func (request *queueRequest) Wait() {
	request.lock.Lock()
	request.lock.Unlock()
}

func (request *queueRequest) Error() error {
	return request.err
}

func (scheduler *queueScheduler) ScheduleHTTPRequest(httpRequest *http.Request,
	output interface{}) request {
	request := new(queueRequest)
	request.lock = new(sync.Mutex)
	request.lock.Lock()
	request.httpRequest = httpRequest
	request.output = output
	scheduler.queue <- request
	return request
}

func (scheduler *queueScheduler) init() {
	scheduler.queue = make(chan *queueRequest, MaxRequestsInQueue)
	scheduler.recentRequests = []*queueRequest{}
	scheduler.httpClient = new(http.Client)
	scheduler.recentRequestsLock = new(sync.RWMutex)
	go manageQueueScheduler(scheduler)
}

func manageQueueScheduler(scheduler *queueScheduler) {
	for {
		manageRecentRequests(scheduler)
		scheduler.recentRequestsLock.RLock()
		if len(scheduler.recentRequests) > MaxRequestsPerSecond {
			scheduler.recentRequestsLock.RUnlock()
			continue
		}
		scheduler.recentRequestsLock.RUnlock()
		request := <-scheduler.queue
		scheduler.recentRequestsLock.Lock()
		scheduler.recentRequests = append(scheduler.recentRequests, request)
		scheduler.recentRequestsLock.Unlock()
		go executeRequest(scheduler, request)
	}
}

func executeRequest(scheduler *queueScheduler, request *queueRequest) {
	logging.LogHTTPRequest(request.httpRequest)
	response, err := scheduler.httpClient.Do(request.httpRequest)
	request.err = err
	if err == nil {
		logging.LogHTTPResponse(request.httpRequest, response)
		err = util.ParseHTTPResponse(response, request.output)
		request.err = err
		err = response.Body.Close()
		if request.err == nil {
			request.err = err
		}
	}
	request.executionTime = time.Now()
	request.releaseTime = time.Now().Add(time.Second)
	request.done = true
	request.lock.Unlock()
}

func manageRecentRequests(scheduler *queueScheduler) {
	scheduler.recentRequestsLock.RLock()
	var toRelease []int
	currentTime := time.Now()
	for i, recentRequest := range scheduler.recentRequests {
		if recentRequest.done && currentTime.After(recentRequest.releaseTime) {
			toRelease = append(toRelease, i)
		}
	}
	scheduler.recentRequestsLock.RUnlock()
	if len(toRelease) == 0 {
		return
	}
	scheduler.recentRequestsLock.Lock()
	for _, index := range toRelease {
		scheduler.recentRequests = append(
			scheduler.recentRequests[0:index],
			scheduler.recentRequests[index+1:]...,
		)
	}
	scheduler.recentRequestsLock.Unlock()
}
