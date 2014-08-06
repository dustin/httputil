// +build !appengine

package httputil

import (
	"os"
	"os/signal"
	"syscall"
)

const sigInfo = syscall.Signal(29)

func stopTracker(ch chan os.Signal) {
	signal.Stop(ch)
}

func initTracker(ch chan os.Signal) {
	signal.Notify(tracker.sigch, sigInfo)
}
