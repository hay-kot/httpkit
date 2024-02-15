package graceful

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	ErrRunnerNotStarted     = errors.New("server not started")
	ErrRunnerAlreadyStarted = errors.New("server already started")
)

// Runner is the orchestrator of the plugins provided. It will start and cancel
// the plugins based on the context and os.Signals provided.
type Runner struct {
	started  bool
	plugins  []Plugin
	shutdown chan struct{}
	opts     *runnerOpts
}

func NewRunner(opts ...RunnerOptFunc) *Runner {
	o := &runnerOpts{
		signals: []os.Signal{os.Interrupt, syscall.SIGTERM},
		timeout: 5 * time.Second,
		println: func(v ...any) {}, // NOOP
	}
	for _, opt := range opts {
		opt(o)
	}

	return &Runner{
		opts:     o,
		shutdown: make(chan struct{}),
	}
}

func (*Runner) Name() string {
	return "runner"
}

// AddPlugin adds a plugin to the server during construction.
// This returns the server for chaining.
func (svr *Runner) AddPlugin(p ...Plugin) {
	svr.plugins = append(svr.plugins, p...)
}

// Start start the server with a context provided for cancellation
// if the root context is cancelled, the server signal stops to all
// plugins registered.
//
// Note that a new context is created with the provided signals defined
// when creating the server.
func (svr *Runner) Start(ctx context.Context) error {
	if svr.started {
		return ErrRunnerAlreadyStarted
	}

	// TODO: add options for signals
	ctx, cancel := signal.NotifyContext(ctx, svr.opts.signals...)
	defer cancel()

	// Start Plugins
	var (
		wg          = sync.WaitGroup{}
		pluginErrCh = make(chan error)
		wgChannel   = make(chan struct{})
	)

	wg.Add(len(svr.plugins))

	go func() {
		wg.Wait()
		close(wgChannel)
	}()

	var plugErr error
	for _, p := range svr.plugins {
		if plugErr != nil {
			break
		}
		go func(p Plugin) {
			defer func() {
				wg.Done()
			}()

			err := p.Start(ctx)
			if err != nil {
				plugErr = err

				// safely write to the channel
				// if the channel is full, we don't want to block
				select {
				case pluginErrCh <- err:
				default:
				}
			}
		}(p)
	}

	go func() {
		<-svr.shutdown
		cancel()
	}()

	svr.started = true
	defer func() {
		svr.started = false
	}()

	// block until the context is done
	select {
	case <-ctx.Done():
		newTimer := time.NewTimer(svr.opts.timeout)
		defer newTimer.Stop()

		svr.opts.println("server received signal, shutting down")
		select {
		case <-wgChannel:
			svr.opts.println("all plugins have stopped, shutting down")
			return nil
		case <-newTimer.C:
			svr.opts.println("timeout waiting for plugins to stop, shutting down")
			return context.DeadlineExceeded
		}
	case err := <-pluginErrCh:
		svr.opts.println("plugin error:", err)
		return err
	}
}

// Shutdown sends a signal to the server to stop all plugins and
// the server itself. This function returns immediately after the
// signal is sent.
func (svr *Runner) Shutdown() {
	close(svr.shutdown)
}
