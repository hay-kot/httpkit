// Package errbus provides a simple error bus for passing errors between
package errbus

import "net/http"

// ErrorBus is an error bus for passing errors between middleware and handlers up the stack.
// You can provide a ErrHandler middleware that will be the outermost middleware in the stack.
// You can also provide a list of GlobalMW middleware that will be applied to all handlers.
type ErrorBus struct {
	ErrHandler ErrMiddleware // Error handler
	GlobalMW   []Middleware  // Global middleware
}

func NewErrorBus() *ErrorBus {
	return &ErrorBus{}
}

func (b *ErrorBus) Use(mw ...Middleware) {
	b.GlobalMW = append(b.GlobalMW, mw...)
}

func (b *ErrorBus) ToHandler(h Handler, mw ...Middleware) http.Handler {
	h = wrapMiddleware(h, append(b.GlobalMW, mw...))
	return b.ErrHandler(h)
}

type ErrMiddleware func(Handler) http.Handler

type Middleware func(Handler) Handler

// wrapMiddleware creates a new handler by wrapping middleware around a final
// handler. The middlewares' Handlers will be executed by requests in the order
// they are provided.
func wrapMiddleware(handler Handler, mw []Middleware) Handler {
	// Loop backwards through the middleware invoking each one. Replace the
	// handler with the new wrapped handler. Looping backwards ensures that the
	// first middleware of the slice is the first to be executed by requests.
	for i := len(mw) - 1; i >= 0; i-- {
		h := mw[i]
		if h != nil {
			handler = h(handler)
		}
	}

	return handler
}

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return f(w, r)
}

type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request) error
}
