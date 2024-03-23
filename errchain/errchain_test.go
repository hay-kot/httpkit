package errchain

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_wrapMiddleware(t *testing.T) {
	middlewares := []Middleware{
		newErrMiddleware("mid 1"),
		newErrMiddleware("mid 2"),
		newErrMiddleware("mid 3"),
		newErrMiddleware("mid 4"),
	}

	// Expected order of execution
	expected := []string{
		"mid 1",
		"mid 2",
		"mid 3",
		"mid 4",
	}

	handler := HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		l, ok := r.Context().Value(slicekeyValue).([]string)
		if !ok {
			l = []string{}
		}

		for i, v := range l {
			if v != expected[i] {
				t.Errorf("expected %s, got %s", expected[i], v)
			}
		}

		return nil
	})

	wrapped := wrapMiddleware(handler, middlewares)

	writer := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/", nil)

	err := wrapped.ServeHTTP(writer, request)
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func Test_AdaptMiddleware(t *testing.T) {
	key := slicekey("AdaptMiddleware")

	httpMid := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l, ok := r.Context().Value(key).([]string)
			if !ok {
				l = []string{}
			}

			l = append(l, "http mid")

			r = r.WithContext(context.WithValue(r.Context(), key, l))
			h.ServeHTTP(w, r)
		})
	}

	newMid := AdaptMiddleware(httpMid)

	handler := HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		l, ok := r.Context().Value(key).([]string)
		if !ok {
			l = []string{}
		}

		if len(l) != 1 {
			t.Errorf("expected 1, got %d", len(l))
		}

		if l[0] != "http mid" {
			t.Errorf("expected http mid, got %s", l[0])
		}

		return nil
	})

	wrapped := newMid(handler)

	writer := httptest.NewRecorder()
	request := httptest.NewRequest("GET", "/", nil)

	err := wrapped.ServeHTTP(writer, request)
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}
