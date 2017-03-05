// +build !appengine

package httputil

import (
	"net/http"
	"time"
)

func reqDeadline(req *http.Request) time.Time {
	deadline, _ := req.Context().Deadline()
	return deadline
}
