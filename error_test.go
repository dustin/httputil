package httputil

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestError(t *testing.T) {
	res := &http.Response{
		StatusCode: 404,
		Status:     "404 Not Found",
		Body:       ioutil.NopCloser(strings.NewReader("omg, not found")),
	}

	got := HTTPError(res).Error()
	exp := `HTTP Error 404 Not Found - omg, not found`
	if got != exp {
		t.Errorf("Expected %q, got %q", exp, got)
	}
}

func TestErrorf(t *testing.T) {
	res := &http.Response{
		StatusCode: 404,
		Status:     "404 Not Found",
		Body:       ioutil.NopCloser(strings.NewReader("omg, not found")),
	}

	got := HTTPErrorf(res, "wtf %q - %B on %S", "quote me").Error()
	exp := `wtf "quote me" - omg, not found on 404 Not Found`
	if got != exp {
		t.Errorf("Expected %q, got %q", exp, got)
	}
}

func TestIsStatus(t *testing.T) {
	if IsHTTPStatus(io.EOF, 404) {
		t.Errorf("Expected false on io.EOF, was somehow a 404")
	}

	res := &http.Response{
		StatusCode: 404,
		Status:     "404 Not Found",
		Body:       ioutil.NopCloser(strings.NewReader("omg, not found")),
	}

	err := HTTPError(res)
	if !IsHTTPStatus(err, 404) {
		t.Errorf("Expected 404 to be a 404, but wasn't. %v", err)
	}

	err = HTTPError(res)
	if IsHTTPStatus(err, 500) {
		t.Errorf("Expected 500 to not be a 404, but was. %v", err)
	}
}
