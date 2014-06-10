package httputil

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"testing"
	"time"
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

func TestTrackerStringing(t *testing.T) {
	defer func(prev http.RoundTripper) { http.DefaultTransport = prev }(http.DefaultTransport)
	x := InitHTTPTracker(false)
	defer x.Close()

	got := x.String()
	if got != "[]" {
		t.Fatalf(`Expected "[]", got %q`, got)
	}
}

func TestEventMarshaling(t *testing.T) {
	u, err := url.Parse("http://www.spy.net/")
	must(err)
	start, err := time.Parse(time.RFC3339, "2014-06-10T09:24:00Z")
	must(err)
	end := start.Add(19 * time.Millisecond).UTC()
	ev := &trackedEvent{
		start,
		timeSrc{
			now: func() time.Time { return end },
		},
		&http.Request{Method: "GET", URL: u},
		nil,
	}

	j, err := json.Marshal(ev)
	must(err)

	exp := `{"duration":19000000,"duration_s":"19ms","method":"GET","startTime":"2014-06-10T09:24:00Z","url":"http://www.spy.net/"}`
	if string(j) != exp {
		t.Errorf("Expected %q, got %q", exp, j)
	}
}

func TestRegistration(t *testing.T) {
	defer func(prev http.RoundTripper) { http.DefaultTransport = prev }(http.DefaultTransport)
	x := InitHTTPTracker(true)
	defer x.Close()

	u, err := url.Parse("http://www.spy.net/")
	must(err)
	one := x.register(&http.Request{Method: "GET", URL: u})
	two := x.register(&http.Request{Method: "GET", URL: u})

	if n := x.count(); n != 2 {
		t.Errorf("Expected two tracked items, got %v", n)
	}
	x.unregister(one)
	if n := x.count(); n != 1 {
		t.Errorf("Expected one tracked items, got %v", n)
	}
	x.unregister(one)
	if n := x.count(); n != 1 {
		t.Errorf("Expected one tracked items, got %v", n)
	}
	x.unregister(two)
	if n := x.count(); n != 0 {
		t.Errorf("Expected zero tracked items, got %v", n)
	}

}

func TestRegistrationIDCollision(t *testing.T) {
	defer func(prev http.RoundTripper) { http.DefaultTransport = prev }(http.DefaultTransport)
	x := InitHTTPTracker(true)
	defer x.Close()

	x.inflight = map[int]trackedEvent{}
	x.inflight[0] = trackedEvent{}
	x.inflight[1] = trackedEvent{}
	x.inflight[2] = trackedEvent{}
	x.inflight[3] = trackedEvent{}

	u, err := url.Parse("http://www.spy.net/")
	must(err)
	id := x.register(&http.Request{Method: "GET", URL: u})

	if id != 4 {
		t.Errorf("Expected ID 4, got %v", id)
	}
}
