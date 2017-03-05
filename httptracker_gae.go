// +build appengine

package httputil

import (
	"net/http"
	"time"
)

func reqDeadline(req *http.Request) time.Time {
	return time.Time{}
}
