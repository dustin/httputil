// Package httputil provides a small junk drawer of http client and
// server helpers.
package httputil

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"
)

type timeSrc struct {
	now   func() time.Time
	since func(time.Time) time.Duration
}

type trackedEvent struct {
	t       time.Time
	tsrc    timeSrc
	req     *http.Request
	callers []uintptr
}

func (t trackedEvent) MarshalJSON() ([]byte, error) {
	frames := []string{}
	for _, pc := range t.callers {
		frame := runtime.FuncForPC(pc)
		fn, line := frame.FileLine(frame.Entry())
		frames = append(frames, fmt.Sprintf("%v() - %v:%v",
			frame.Name(), fn, line))
	}
	now := t.tsrc.now()
	ob := map[string]interface{}{
		"startTime":  t.t,
		"duration":   now.Sub(t.t),
		"duration_s": now.Sub(t.t).String(),
		"method":     t.req.Method,
		"url":        t.req.URL.String(),
	}
	if len(frames) > 0 {
		ob["stack"] = frames
	}
	return json.Marshal(ob)
}

// HTTPTracker is a http.RoundTripper wrapper that tracks usage of clients.
//
// The easiest way to use http tracker for a commandline tool is to
// just call InitHTTPTracker:
//
//    httputil.InitHTTPTracker(false)
//
// This wraps the current http.DefaultTransport with a tracking
// transport and installs a SIGINFO handler to report the current
// state on demand.
//
// If you have a web server that is also an HTTP client and uses
// expvar, you can publish an expvar version of the data with the
// following, similar invocation:
//
//   expvar.Publish("httpclients", httputil.InitHTTPTracker(false))
//
// The boolean parameter in the above examples determines whether
// stacks are also tracked.  See the docs for InitHTTPTracker for more
// details.
type HTTPTracker struct {
	// Next is the RoundTripper being wrapped
	Next http.RoundTripper
	// TrackStacks will record user stacks if true.
	TrackStacks bool

	tsrc     timeSrc
	mu       sync.Mutex
	inflight map[int]trackedEvent
	nextID   int
	sigch    chan os.Signal
}

// MarshalJSON provides a JSON representation of the state of tracker.
func (t *HTTPTracker) MarshalJSON() ([]byte, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	events := make([]trackedEvent, 0, len(t.inflight))
	for _, e := range t.inflight {
		events = append(events, e)
	}
	return json.Marshal(events)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

// String produces a JSON formatted representation of the tracker
// state.  This is directly useful to expvar.
func (t *HTTPTracker) String() string {
	b, err := t.MarshalJSON()
	must(err)
	return string(b)
}

func (t *HTTPTracker) count() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.inflight)
}

func (t *HTTPTracker) register(req *http.Request) int {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.inflight == nil {
		t.inflight = map[int]trackedEvent{}
	}
	_, exists := t.inflight[t.nextID]
	for exists {
		t.nextID++
		_, exists = t.inflight[t.nextID]
	}
	thisID := t.nextID
	t.nextID++

	var pcs []uintptr
	if t.TrackStacks {
		pcs = make([]uintptr, 64)
		n := runtime.Callers(3, pcs)
		pcs = pcs[:n-1]
	}
	t.inflight[thisID] = trackedEvent{t.tsrc.now(), t.tsrc, req, pcs}
	return thisID
}

func (t *HTTPTracker) unregister(id int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.inflight, id)
}

// Report writes a textual report of the current state of HTTP clients
// to the given writer.
func (t *HTTPTracker) Report(w io.Writer) {
	t.mu.Lock()
	defer t.mu.Unlock()

	fmt.Fprintf(w, "In-flight HTTP requests:\n")
	for _, e := range t.inflight {
		fmt.Fprintf(w, "  servicing %v %q for %v\n", e.req.Method, e.req.URL, t.tsrc.since(e.t))
		for _, caller := range e.callers {
			frame := runtime.FuncForPC(caller)
			fn, line := frame.FileLine(frame.Entry())
			fmt.Fprintf(w, "    - %v() - %v:%v\n", frame.Name(), fn, line)
		}
	}
}

// ReportLoop will issue a report on the given Writer whenever it
// receives a signal on the given channel.
//
// This primarily exists for signal handlers.
func (t *HTTPTracker) ReportLoop(w io.Writer, ch <-chan os.Signal) {
	for _ = range ch {
		t.Report(w)
	}
}

type trackFinalizer struct {
	b  io.ReadCloser
	t  *HTTPTracker
	id int
}

func (d *trackFinalizer) Close() error {
	d.t.unregister(d.id)
	return d.b.Close()
}

func (d *trackFinalizer) Read(b []byte) (int, error) {
	return d.b.Read(b)
}

func (d *trackFinalizer) WriteTo(w io.Writer) (n int64, err error) {
	return io.Copy(w, d.b)
}

// RoundTrip satisfies http.RoundTripper
func (t *HTTPTracker) RoundTrip(req *http.Request) (*http.Response, error) {
	id := t.register(req)
	res, err := t.Next.RoundTrip(req)
	if err == nil {
		res.Body = &trackFinalizer{res.Body, t, id}
	} else {
		t.unregister(id)
	}
	return res, err
}

// Close shuts down this tracker.
func (t *HTTPTracker) Close() error {
	stopTracker(t.sigch)
	close(t.sigch)
	return nil
}

// InitHTTPTracker wraps http.DefaultTransport with a tracking
// DefaultTransport and installs a SIGINFO handler to report progress.
//
// If trackStacks is true, the call stack will be included with
// tracking information and reports.
func InitHTTPTracker(trackStacks bool) *HTTPTracker {
	tracker := &HTTPTracker{
		Next:        http.DefaultTransport,
		TrackStacks: trackStacks,
		sigch:       make(chan os.Signal, 1),
		tsrc:        timeSrc{time.Now, time.Since},
	}

	http.DefaultTransport = tracker

	go tracker.ReportLoop(os.Stdout, tracker.sigch)
	return tracker
}
