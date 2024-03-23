package errchain

import (
	"context"
	"net/http"
)

type slicekey string

var slicekeyValue = slicekey("slicekey")

func newErrMiddleware(str string) Middleware {
	return func(h Handler) Handler {
		return HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
			l, ok := r.Context().Value(slicekeyValue).([]string)
			if !ok {
				l = []string{}
			}

			l = append(l, str)

			r = r.WithContext(context.WithValue(r.Context(), slicekeyValue, l))
			return h.ServeHTTP(w, r)
		})
	}
}

func newStdMiddleware(str string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l, ok := r.Context().Value(slicekeyValue).([]string)
			if !ok {
				l = []string{}
			}

			l = append(l, str)

			r = r.WithContext(context.WithValue(r.Context(), slicekeyValue, l))
			h.ServeHTTP(w, r)
		})
	}
}
