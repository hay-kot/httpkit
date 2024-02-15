package graceful

import (
	"os"
	"time"
)

type runnerOpts struct {
	signals []os.Signal
	timeout time.Duration
	println func(...any)
}

type RunnerOptFunc func(*runnerOpts)

// WithSignals provides a list of signals to listen for
// when starting the server that will cancel the context
//
// Defaults to
//   - os.Interrupt
//   - syscall.SIGTERM
//
// Multiple calls to this option will override the previous
func WithSignals(signals ...os.Signal) RunnerOptFunc {
	return func(o *runnerOpts) {
		o.signals = signals
	}
}

// WithTimeout provides a timeout for the server to wait for
// plugins to stop before shutting down.
//
// Defaults to 5 seconds
func WithTimeout(timeout time.Duration) RunnerOptFunc {
	return func(o *runnerOpts) {
		o.timeout = timeout
	}
}

// WithPrintln provides a function to print messages
func WithPrintln(fn func(...any)) RunnerOptFunc {
	return func(o *runnerOpts) {
		o.println = fn
	}
}
