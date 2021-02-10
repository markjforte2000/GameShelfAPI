package scheduling

import (
	"net/http"
)

type Scheduler interface {
	ScheduleHTTPRequest(request *http.Request, output interface{}) request
	init()
}

type request interface {
	Wait()
	Error() error
}

func NewScheduler() Scheduler {
	scheduler := new(queueScheduler)
	scheduler.init()
	return scheduler
}
