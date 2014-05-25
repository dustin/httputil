package httputil

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

const maxBody = 512

type httpError struct {
	format string
	args   []interface{}
	code   int
	status string
	body   []byte
}

// Error statisfies the "error" interface.
func (e httpError) Error() string {
	r := strings.NewReplacer("%S", e.status, "%B", string(e.body))
	return fmt.Sprintf(r.Replace(e.format), e.args...)
}

// HTTPError converts an http response into an error.
//
// Note that this reads the body, so only use it when the response
// exists and you don't believe it's valid for your needs.
func HTTPError(res *http.Response) error {
	return HTTPErrorf(res, "HTTP Error %S - %B")
}

// HTTPErrorf converts an http response into an error.
//
// This allows for standard printf-style formatting with the addition
// of %S for the http status (e.g. "404 Not Found") and %B for the
// body that was returned along with the error.
//
// Note that this reads the body, so only use it when the response
// exists and you don't believe it's valid for your needs.
func HTTPErrorf(res *http.Response, format string, args ...interface{}) error {
	// The read error is explicitly ignored here because we're
	// only trying to use it to produce a prettier error.
	body, _ := ioutil.ReadAll(io.LimitReader(res.Body, maxBody))

	return httpError{
		format: format,
		args:   args,
		code:   res.StatusCode,
		status: res.Status,
		body:   body,
	}
}
