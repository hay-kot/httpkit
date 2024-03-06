package errchain

import "net/http"

// Router is a simple wrapper around the standard http.ServeMux with it's new
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
type Router struct {
	http.ServeMux
	prefix string
	chain  *ErrChain
}

// NewRouter creates a new Router with the provided prefix and ErrChain.
// If the ErrChain is nil, this function will panic.
//
// if the prefix is not empty, it will be preprocessed to ensure that it
// starts with a '/'. If the prefix is only a single '/' then
// it is left as is.
func NewRouter(prefix string, chain *ErrChain) *Router {
	if chain == nil {
		panic("errchain: ErrChain is nil")
	}

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

	return &Router{prefix: prefix, chain: chain}
}

// applyMethod constructs a applyMethod string from the method and applyMethod. This
// assumes that r.prefix has a '/' at the beginning and end of
// the string.
func (r *Router) applyMethod(method string, path string) string {
	// If the path does not start with a '/', then add it
	if len(path) > 0 && path[0] != '/' {
		path = "/" + path
	}

	return method + " " + r.prefix + path
}

func (r *Router) applyPrefix(path string) string {
	// If the path does not start with a '/', then add it
	if len(path) > 0 && path[0] != '/' {
		path = "/" + path
	}

	return r.prefix + path
}

func (r *Router) ErrHandler(path string, h Handler) {
	r.ServeMux.Handle(r.applyPrefix(path), r.chain.ToHandler(h))
}

func (r *Router) Method(method, path string, h Handler, mw ...Middleware) {
	r.ServeMux.Handle(r.applyMethod(method, path), r.chain.ToHandler(h, mw...))
}

func (r *Router) MethodFunc(method, path string, fn HandlerFunc, mw ...Middleware) {
	r.ServeMux.Handle(r.applyMethod(method, path), r.chain.ToHandlerFunc(fn, mw...))
}

func (r *Router) Get(path string, fn HandlerFunc, mw ...Middleware) {
	r.MethodFunc(http.MethodGet, path, fn, mw...)
}

func (r *Router) Post(path string, fn HandlerFunc, mw ...Middleware) {
	r.MethodFunc(http.MethodPost, path, fn, mw...)
}

func (r *Router) Put(path string, fn HandlerFunc, mw ...Middleware) {
	r.MethodFunc(http.MethodPut, path, fn, mw...)
}

func (r *Router) Patch(path string, fn HandlerFunc, mw ...Middleware) {
	r.MethodFunc(http.MethodPatch, path, fn, mw...)
}

func (r *Router) Delete(path string, fn HandlerFunc, mw ...Middleware) {
	r.MethodFunc(http.MethodDelete, path, fn, mw...)
}

func (r *Router) Options(path string, fn HandlerFunc, mw ...Middleware) {
	r.MethodFunc(http.MethodOptions, path, fn, mw...)
}

func (r *Router) Head(path string, fn HandlerFunc, mw ...Middleware) {
	r.MethodFunc(http.MethodHead, path, fn, mw...)
}

func (r *Router) Connect(path string, fn HandlerFunc, mw ...Middleware) {
	r.MethodFunc(http.MethodConnect, path, fn, mw...)
}

func (r *Router) Trace(path string, fn HandlerFunc, mw ...Middleware) {
	r.MethodFunc(http.MethodTrace, path, fn, mw...)
}
