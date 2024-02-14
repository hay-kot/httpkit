package graceful_test

import (
	"context"
	"errors"
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

func Test_Runner_LifeCycle(t *testing.T) {
	type plugResults struct {
		start, stop bool
	}

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
		plug1Got.start = true
		<-ctx.Done()
		plug1Got.stop = true
		return nil
	}))

	runner.AddPlugin(graceful.PluginFunc("plug2", func(ctx context.Context) error {
		plug2Got.start = true
		<-ctx.Done()
		plug2Got.stop = true
		return nil
	}))

	runner.AddPlugin(graceful.PluginFunc("plug3", func(ctx context.Context) error {
		plug3Got.start = true
		<-ctx.Done()
		<-time.After(10 * time.Second)
		plug3Got.stop = true
		return nil
	}))

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})

	go func() {
		_ = runner.Start(ctx)
		done <- struct{}{}
	}()

	cancel()
	<-done

	// Asserts
	assert(t, plug1Got.start, true)
	assert(t, plug1Got.stop, true)

	assert(t, plug2Got.start, true)
	assert(t, plug2Got.stop, true)

	assert(t, plug3Got.start, true)
	assert(t, plug3Got.stop, false)
}

func assert[T comparable](t *testing.T, got, expect T) {
	t.Helper()
	if expect != got {
		t.Errorf("expect %v, got %v", expect, got)
	}
}
