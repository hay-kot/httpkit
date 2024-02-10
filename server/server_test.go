package server

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	host = "127.0.0.1"
	port = "19245"
)

func urlpath(path string) string {
	return fmt.Sprintf("http://%s:%s%s", host, port, path)
}

func init() { // nolint:gochecknoinits
	// Set random ports
	port = fmt.Sprintf("%d", rand.Intn(10000)+10000)
}

func testServer(t *testing.T, r http.Handler) *Server {
	svr := NewServer(WithHost(host), WithPort(port))

	mux := http.NewServeMux()
	if r != nil {
		mux.Handle("/", r)
	}
	go func() {
		_ = svr.Start(mux)
	}()

	ping := func() error {
		_, err := http.Get(urlpath("/")) // nolint:bodyclose
		return err
	}

	for {
		if err := ping(); err == nil {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}

	return svr
}

func Test_ServerShutdown_Error(t *testing.T) {
	svr := NewServer(WithHost(host), WithPort(port))

	err := svr.Shutdown("test")
	assert.ErrorIs(t, err, ErrServerNotStarted)
}

func Test_ServerStarts_Error(t *testing.T) {
	svr := testServer(t, nil)

	err := svr.Start(http.NewServeMux())
	assert.ErrorIs(t, err, ErrServerAlreadyStarted)

	err = svr.Shutdown("test")
	require.NoError(t, err)
}

func Test_ServerStarts(t *testing.T) {
	svr := testServer(t, nil)
	err := svr.Shutdown("test")
	require.NoError(t, err)
}

func Test_GracefulServerShutdownWithWorkers(t *testing.T) {
	blockingChannel := make(chan struct{})
	finishedChannel := make(chan string, 1)

	svr := testServer(t, nil)

	svr.Background(func() {
		// Block until the channel is closed
		<-blockingChannel

		// Set the flag to true
		finishedChannel <- "worker finished"
	})

	// Shutdown the server
	var wg sync.WaitGroup
	var err error

	wg.Add(1)
	go func() {
		defer wg.Done()
		err = svr.Shutdown("test")
	}()

	require.NoError(t, err)
	close(blockingChannel)

	wg.Wait()
	require.NoError(t, err)

	result := <-finishedChannel

	assert.Equal(t, "worker finished", result)
}

func Test_GracefulServerShutdownWithRequests(t *testing.T) {
	mux := http.NewServeMux()

	requestStarted := make(chan struct{})
	blockingChannel := make(chan struct{})
	finishedChannel := make(chan string, 1)

	// add long running handler func
	mux.HandleFunc("/test", func(rw http.ResponseWriter, r *http.Request) {
		requestStarted <- struct{}{}

		// Block until the channel is closed
		<-blockingChannel

		// Set the flag to true
		finishedChannel <- "handler finished"
	})

	svr := testServer(t, mux)

	// Make request to "/test"
	go func() {
		_, _ = http.Get(urlpath("/test")) // nolint:bodyclose
	}()

	// Wait for the request to start
	<-requestStarted

	// Shutdown the server
	var wg sync.WaitGroup
	var err error

	wg.Add(1)
	go func() {
		defer wg.Done()
		err = svr.Shutdown("test")
	}()

	require.NoError(t, err)
	close(blockingChannel)

	wg.Wait()
	require.NoError(t, err)

	close(finishedChannel)
	result := <-finishedChannel
	assert.Equal(t, "handler finished", result)
}
