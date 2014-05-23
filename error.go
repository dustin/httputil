package httputil

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const maxBody = 512

type httpError struct {
	code   int
	status string
	body   []byte
}

// Error statisfies the "error" interface.
func (e httpError) Error() string {
	return fmt.Sprintf("HTTP Error %v - %s", e.status, e.body)
}

// HTTPError converts an http response into an error.
//
// Note that this reads the body, so only use it when the response
// exists and you don't believe it's valid for your needs.
func HTTPError(res *http.Response) error {
	// The read error is explicitly ignored here because we're
	// only trying to use it to produce a prettier error.
	body, _ := ioutil.ReadAll(io.LimitReader(res.Body, maxBody))

	return httpError{
		code:   res.StatusCode,
		status: res.Status,
		body:   body,
	}
}
