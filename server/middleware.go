package server

import (
	"net/http"
	"strings"
)

// StripTrailingSlash is a middleware that will strip trailing slashes from the request path.
//
// Example: /api/v1/ -> /api/v1
//
// Deprecated: Use alternative middleware from Chi or Gorilla.
func StripTrailingSlash(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		}
		next.ServeHTTP(w, r)
	})
}
