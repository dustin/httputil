package httputil

import (
	"errors"
	"net/http"
	"testing"
)

func TestTrackerInit(t *testing.T) {
	defer func(prev http.RoundTripper) { http.DefaultTransport = prev }(http.DefaultTransport)
	x := InitHTTPTracker(false)
	defer x.Close()
	if http.DefaultTransport != x {
		t.Errorf("Expected transport to be installed. Was %v", http.DefaultTransport)
	}
	if x.TrackStacks {
		t.Errorf("Expected to track stacks. Won't")
	}
}

func TestTrackerInitNoTrack(t *testing.T) {
	defer func(prev http.RoundTripper) { http.DefaultTransport = prev }(http.DefaultTransport)
	x := InitHTTPTracker(true)
	defer x.Close()
	if http.DefaultTransport != x {
		t.Errorf("Expected transport to be installed. Was %v", http.DefaultTransport)
	}
	if !x.TrackStacks {
		t.Errorf("Expected to not track stacks. Will")
	}
}

func TestMust(t *testing.T) {
	must(nil) // be nice to not panic here
	panicked := false
	func() {
		defer func() { _, panicked = recover().(error) }()
		must(errors.New("break here"))
	}()
	if !panicked {
		t.Fatalf("Expected a panic, but didn't get one.")
	}
}
