// +build !appengine

package httputil

import (
	"os"
	"os/signal"
	"syscall"
)

const sigInfo = syscall.Signal(29)

func initTracker(ch chan os.Signal) {
	signal.Notify(tracker.sigch, sigInfo)
}
