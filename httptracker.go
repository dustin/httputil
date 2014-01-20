// Package httputil provides a small junk drawer of http client and
// server helpers.
package httputil

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"
)

const sigInfo = syscall.Signal(29)

type trackedEvent struct {
	t       time.Time
	req     *http.Request
	callers []uintptr
}

// HTTPTracker is a http.RoundTripper wrapper that tracks usage of clients.
type HTTPTracker struct {
	// Next is the RoundTripper being wrapped
	Next http.RoundTripper
	// TrackStacks will record user stacks if true.
	TrackStacks bool

	mu       sync.Mutex
	inflight map[int]trackedEvent
	nextID   int
}

func (t *HTTPTracker) register(req *http.Request) int {
	t.mu.Lock()
	defer t.mu.Unlock()
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
	t.inflight[thisID] = trackedEvent{time.Now(), req, pcs}
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
		fmt.Fprintf(w, "  servicing %v %q for %v\n", e.req.Method, e.req.URL, time.Since(e.t))
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

// InitHTTPTracker wraps http.DefaultTransport with a tracking
// DefaultTransport and installs a SIGINFO handler to report progress.
//
// If trackStacks is true, the call stack will be included with
// tracking information and reports.
func InitHTTPTracker(trackStacks bool) {
	http.DefaultTransport = &HTTPTracker{
		Next:        http.DefaultTransport,
		TrackStacks: trackStacks,
		inflight:    map[int]trackedEvent{},
	}

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, sigInfo)

	go http.DefaultTransport.(*HTTPTracker).ReportLoop(os.Stdout, sigch)
}
