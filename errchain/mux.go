package errchain

import "net/http"

// HandlerHook is a function that can be used to add hooks to the mux.
// this was implemented to allow for the otelhttp.WithRouteTag method to
// be easily added to the mux, but may serve other purposes as well.
type HandlerHook func(pattern string, handler http.Handler) http.Handler

// Mux is a simple wrapper around the standard http.ServeMux with it's new
// improvements in go 1.22. It extends the http.ServeMux with the following
// features:
//
// 1. allows prefixing of all routes
// 2. allows for using errchain.Handler types directly
// 3. provides method based functions for GET, POST, etc...
//
// Because route embeds the http.Servermux, it can be used as a drop-in
// replacement for the http.ServeMux. However, because of this, if you
// want to use errchain.Handler types you need to use the ErrHandler
// method to add the handler to the route.
type Mux struct {
	mw []func(http.Handler) http.Handler
	http.ServeMux
	prefix string
	chain  *ErrChain

	// Hook is a function that can be used to add hooks to the mux. This
	// was implemented to allow for the otelhttp.WithRouteTag method to be
	// easily added to the mux, but may serve other purposes as well.
	//
	// Do not use this for middleware. Middleware should be added using the
	// Use method, added to the ErrChain, or added to the individual routes.
	Hook HandlerHook
}

// NewMux creates a new Mux with the provided prefix and ErrChain.
// If the ErrChain is nil, this function will panic.
//
// if the prefix is not empty, it will be preprocessed to ensure that it
// starts with a '/'. If the prefix is only a single '/' then
// it is left as is.
//
// If chain is nil, this function will panic.
func NewMux(chain *ErrChain) *Mux {
	if chain == nil {
		panic("errchain: ErrChain is nil")
	}

	return &Mux{chain: chain}
}

func (r *Mux) handle(path string, h Handler, mw ...Middleware) {
	hdlr := r.chain.ToHandler(h, mw...)

	if len(r.mw) > 0 {
		for i := len(r.mw) - 1; i >= 0; i-- {
			hdlr = r.mw[i](hdlr)
		}
	}

	if r.Hook != nil {
		hdlr = r.Hook(path, hdlr)
	}

	r.ServeMux.Handle(path, hdlr)
}

// UsePrefix sets the prefix for the Mux. This prefix will be prepended to all
//
// Example:
//
//	mux := errchain.New(chain)
//	mux = mux.UsePrefix("/api")
func (r *Mux) UsePrefix(prefix string) *Mux {
	// preprocess the prefix to ensure it starts and ends with a '/'
	if prefix == "/" {
		prefix = ""
	}

	if prefix != "" {
		// Ensure that the prefix starts with a '/'
		if prefix[0] != '/' {
			prefix = "/" + prefix
		}

		// Ensure that the prefix _does not_ end with a '/'
		if prefix[len(prefix)-1] == '/' {
			prefix = prefix[:len(prefix)-1]
		}
	}

	r.prefix = prefix
	return r
}

// Use adds middleware to the mux. This middleware will be applied to all
// routes that are added to the mux, regardless of the method. These are
// applied outside of the error chain middleware.
//
// Example
//
//	mux.Use(middleware1, middleware2)
//	mux.Get("/path", handler, middleware3)
//
// in this the call order is:
//  1. middleware1
//  2. middleware2
//  3. err chain middleware
//  4. middleware3
//  5. handler
func (r *Mux) Use(mw ...func(http.Handler) http.Handler) {
	r.mw = append(r.mw, mw...)
}

// Method adds a handler to the mux for the provided method and path. The
// path is automatically prefixed with the mux's prefix set during creation.
// The handler is wrapped in the error chain middleware and any additional
// middleware provided is applied to the handler.
//
// Example:
//
//	mux.Method(http.MethodGet, "/path", handler, middleware1, middleware2)
func (r *Mux) Method(method, path string, h Handler, mw ...Middleware) {
	if len(path) > 0 && path[0] != '/' {
		path = "/" + path
	}

	path = r.prefix + path
	if method != "" {
		path = method + " " + path
	}

	r.handle(path, h, mw...)
}

// ErrHandle adds a handler to the mux for the provided path. It does not
// bind the handler to a specific method. The path is automatically prefixed
// with the mux's prefix set during creation. The handler is wrapped in the
// error chain middleware and any additional middleware provided is applied to
// the handler.
//
// This is a errchain.Handler specific version of the http.ServeMux.Handle.
func (r *Mux) ErrHandle(path string, h Handler, mw ...Middleware) {
	r.Method("", path, h, mw...)
}

// MethodFunc is a shortcut for Method(method, path, h, mw...) but uses a
// HandlerFunc instead of a Handler.
func (r *Mux) MethodFunc(method, path string, h HandlerFunc, mw ...Middleware) {
	r.Method(method, path, h, mw...)
}

// Get is a shortcut for Method(http.MethodGet, path, h, mw...)
func (r *Mux) Get(path string, fn HandlerFunc, mw ...Middleware) {
	r.Method(http.MethodGet, path, fn, mw...)
}

// Post is a shortcut for Method(http.MethodPost, path, h, mw...)
func (r *Mux) Post(path string, fn HandlerFunc, mw ...Middleware) {
	r.Method(http.MethodPost, path, fn, mw...)
}

// Put is a shortcut for Method(http.MethodPut, path, h, mw...)
func (r *Mux) Put(path string, fn HandlerFunc, mw ...Middleware) {
	r.Method(http.MethodPut, path, fn, mw...)
}

// Patch is a shortcut for Method(http.MethodPatch, path, h, mw...)
func (r *Mux) Patch(path string, fn HandlerFunc, mw ...Middleware) {
	r.Method(http.MethodPatch, path, fn, mw...)
}

// Delete is a shortcut for Method(http.MethodDelete, path, h, mw...)
func (r *Mux) Delete(path string, fn HandlerFunc, mw ...Middleware) {
	r.Method(http.MethodDelete, path, fn, mw...)
}

// Options is a shortcut for Method(http.MethodOptions, path, h, mw...)
func (r *Mux) Options(path string, fn HandlerFunc, mw ...Middleware) {
	r.Method(http.MethodOptions, path, fn, mw...)
}

// Head is a shortcut for Method(http.MethodHead, path, h, mw...)
func (r *Mux) Head(path string, fn HandlerFunc, mw ...Middleware) {
	r.Method(http.MethodHead, path, fn, mw...)
}

// Connect is a shortcut for Method(http.MethodConnect, path, h, mw...)
func (r *Mux) Connect(path string, fn HandlerFunc, mw ...Middleware) {
	r.Method(http.MethodConnect, path, fn, mw...)
}

// Trace is a shortcut for Method(http.MethodTrace, path, h, mw...)
func (r *Mux) Trace(path string, fn HandlerFunc, mw ...Middleware) {
	r.Method(http.MethodTrace, path, fn, mw...)
}
