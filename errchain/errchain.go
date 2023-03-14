// Package errchain provides a simple error bus for passing errors between
package errchain

import (
	"net/http"
)

// Handler is an interface that is an alternative to the http.Handler interface.
// It is the same as http.Handler except that it returns an error.
type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request) error
}

// HandlerFunc is an adapter to allow the use of ordinary functions as handlers.
type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return f(w, r)
}

// ErrorHandler defines your outermost middleware that will be used to handle errors from
// the handler and middleware chain.
type ErrorHandler func(Handler) http.Handler

type Middleware func(Handler) Handler

// ErrChain is an error bus for passing errors between middleware and handlers up the stack.
// You can provide a ErrHandler middleware that will be the outermost middleware in the stack.
// You can also provide a list of GlobalMW middleware that will be applied to all handlers.
type ErrChain struct {
	errorHandler ErrorHandler // Error handler
	globalMW     []Middleware // Global middleware
}

// New creates a new ErrChain with the provided ErrHandler.
// If the ErrHandler is nil, this function will panic.
func New(hdlr ErrorHandler) *ErrChain {
	if hdlr == nil {
		panic("errchain: ErrHandler is nil")
	}

	return &ErrChain{
		errorHandler: hdlr,
	}
}

// Use adds middleware to the global middleware chain. This middleware will be
// applied to all handlers. Note that the middleware will be applied in the order
// that it is provided. So the first middleware provided will be the first to be
// executed by requests and the last middleware provided will be the last to be
// executed by requests.
//
// When called multiple times, the middleware will be added to the end of the
// chain.
//
// Example:
//
//	chain.Use(middleware1, middleware2, middleware3)
//
// The above example will result in the following execution order:
//
// middleware1 -> middleware2 -> middleware3 -> {handler specific} -> handler
func (b *ErrChain) Use(mw ...Middleware) {
	b.globalMW = append(b.globalMW, mw...)
}

// ToHandler converts a Handler to an http.Handler. This function will wrap the
// provided handler with the global middleware and the provided middleware.
// The middleware will be applied in the order that it is provided. So the first
// middleware provided will be the first to be executed by requests and the last
// middleware provided will be the last to be executed by requests.
//
// Example:
//
//	chain.ToHandler(handler, middleware1, middleware2, middleware3)
//
// The above example will result in the following execution order:
//
// middleware1 -> middleware2 -> middleware3 -> {global} -> handler
func (b *ErrChain) ToHandler(h Handler, mw ...Middleware) http.Handler {
	h = wrapMiddleware(h, append(b.globalMW, mw...))
	return b.errorHandler(h)
}

// ToHandlerFunc does the same thing as ToHandler except that it returns a http.HandlerFunc.
func (b *ErrChain) ToHandlerFunc(h HandlerFunc, mw ...Middleware) http.HandlerFunc {
	return b.ToHandler(h, mw...).ServeHTTP
}

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
