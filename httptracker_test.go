package httputil

import (
	"net/http"
	"testing"
)

func TestTrackerInit(t *testing.T) {
	defer func(prev http.RoundTripper) { http.DefaultTransport = prev }(http.DefaultTransport)
	x := InitHTTPTracker(false)
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
	if http.DefaultTransport != x {
		t.Errorf("Expected transport to be installed. Was %v", http.DefaultTransport)
	}
	if !x.TrackStacks {
		t.Errorf("Expected to not track stacks. Will")
	}
}
