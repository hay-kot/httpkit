package graceful_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/hay-kot/httpkit/graceful"
)

func Test_Runner_FailedStartup(t *testing.T) {
	runner := graceful.NewRunner(
		graceful.WithTimeout(3*time.Millisecond),
		graceful.WithPrintln(func(args ...any) {
			t.Log(args)
		}),
	)

	runner.AddPlugin(graceful.PluginFunc("plug1", func(ctx context.Context) error {
		return errors.New("failed to start")
	}))

	err := runner.Start(context.Background())

	assert(t, err.Error(), "failed to start")
}

type plugResults struct {
	mu          sync.Mutex
	start, stop bool
}

func (p *plugResults) setStart(v bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.start = v
}

func (p *plugResults) setStop(v bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.stop = v
}

func (p *plugResults) getStart() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.start
}

func (p *plugResults) getStop() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.stop
}

func Test_Runner_LifeCycle(t *testing.T) {
	runner := graceful.NewRunner(
		graceful.WithTimeout(3*time.Millisecond),
		graceful.WithPrintln(func(args ...any) {
			t.Log(args)
		}),
	)

	plug1Got := plugResults{}
	plug2Got := plugResults{}
	plug3Got := plugResults{}

	runner.AddPlugin(graceful.PluginFunc("plug1", func(ctx context.Context) error {
		plug1Got.setStart(true)
		<-ctx.Done()
		plug1Got.setStop(true)
		return nil
	}))

	runner.AddPlugin(graceful.PluginFunc("plug2", func(ctx context.Context) error {
		plug2Got.setStart(true)
		<-ctx.Done()
		plug2Got.setStop(true)
		return nil
	}))

	runner.AddPlugin(graceful.PluginFunc("plug3", func(ctx context.Context) error {
		plug3Got.setStart(true)
		<-ctx.Done()

		// Block forever
		<-make(chan struct{})
		return nil
	}))

	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = runner.Start(ctx)
	}()

	cancel()
	wg.Wait()

	// Asserts
	assert(t, plug1Got.getStart(), true)
	assert(t, plug1Got.getStop(), true)

	assert(t, plug2Got.getStart(), true)
	assert(t, plug2Got.getStop(), true)

	assert(t, plug3Got.getStart(), true)
	assert(t, plug3Got.getStop(), false)
}

func assert[T comparable](t *testing.T, got, expect T) {
	t.Helper()
	if expect != got {
		t.Errorf("expect %v, got %v", expect, got)
	}
}
