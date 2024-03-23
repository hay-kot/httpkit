package errchain

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var TestErrHandler = func(h Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := h.ServeHTTP(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

func Test_Router_PathConstruction(t *testing.T) {
	tests := []struct {
		name   string
		path   string
		prefix string
		expect string
	}{
		{
			name:   "no prefix",
			path:   "/apitest/v1/endpoint",
			prefix: "",
			expect: "/apitest/v1/endpoint",
		},
		{
			name:   "with prefix",
			path:   "/apitest/v1/endpoint",
			prefix: "/api",
			expect: "/api/apitest/v1/endpoint",
		},
		{
			name:   "with trailing slash",
			path:   "/apitest/v1/endpoint/",
			prefix: "/api/",
			expect: "/api/apitest/v1/endpoint",
		},
		{
			name:   "with no slash",
			path:   "apitest/v1/endpoint",
			prefix: "api",
			expect: "/api/apitest/v1/endpoint",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := NewRouter(test.prefix, New(TestErrHandler))

			handler := func(w http.ResponseWriter, r *http.Request) error {
				w.WriteHeader(http.StatusOK)
				return nil
			}

			server := httptest.NewServer(r)

			r.MethodFunc(http.MethodGet, test.path, handler)

			client := server.Client()

			req, err := http.NewRequest(http.MethodGet, server.URL+test.expect, nil)
			if err != nil {
				t.Fatal(err)
			}

			resp, err := client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = resp.Body.Close() }()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
			}
		})
	}
}

func Test_Router_MethodFunc(t *testing.T) {
	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodOptions,
		http.MethodHead,
		http.MethodConnect,
		http.MethodTrace,
	}

	path := "/apitest/v1/endpoint"

	handler := func(w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(http.StatusOK)
		return nil
	}

	for i, method := range methods {
		t.Run(method, func(t *testing.T) {
			r := NewRouter("", New(TestErrHandler))
			r.MethodFunc(method, path, handler)

			server := httptest.NewServer(r)
			defer server.Close()

			// Make a request to the server using the method
			client := server.Client()

			// ensure other methods are not allowed
			for j, otherMethod := range methods {
				req, err := http.NewRequest(otherMethod, server.URL+path, nil)
				if err != nil {
					t.Fatal(err)
				}

				resp, err := client.Do(req)
				if err != nil {
					t.Fatal(err)
				}
				defer func() { _ = resp.Body.Close() }()

				// if the method is the same OR if the method is GET
				// and the other method is HEAD, then the status code
				// should be 200
				if i == j || (method == "GET" && otherMethod == "HEAD") {
					if resp.StatusCode != http.StatusOK {
						t.Logf("Test: [%s] %s", method, path)
						t.Errorf("[%s] %s ~ expected status code %d, got %d",
							otherMethod,
							path,
							http.StatusOK,
							resp.StatusCode,
						)
					}
					continue
				}

				if resp.StatusCode != http.StatusMethodNotAllowed {
					t.Logf("Test: [%s] %s", method, path)
					t.Errorf("[%s] %s ~ expected status code %d, got %d",
						otherMethod,
						path,
						http.StatusMethodNotAllowed,
						resp.StatusCode,
					)
				}
			}
		})
	}
}

func Test_Router_MiddlewareOrder(t *testing.T) {
	stdMiddlewares := []func(http.Handler) http.Handler{
		newStdMiddleware("mid 1"),
		newStdMiddleware("mid 2"),
	}

	customMiddlewares := []Middleware{
		newErrMiddleware("err mid 1"),
		newErrMiddleware("err mid 2"),
	}

	handlerMiddleware := []Middleware{
		newErrMiddleware("handler 1"),
		newErrMiddleware("handler 2"),
	}

	// Expected order of execution
	expected := []string{
		"mid 1",
		"mid 2",
		"err mid 1",
		"err mid 2",
		"handler 1",
		"handler 2",
	}

	chain := New(TestErrHandler)

	r := NewRouter("", chain)

	chain.Use(customMiddlewares...)
	r.Use(stdMiddlewares...)

	handler := HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
		l, ok := r.Context().Value(slicekeyValue).([]string)
		if !ok {
			l = []string{}
		}

		var failed bool
		for i, v := range l {
			if v != expected[i] {
				t.Errorf("expected %s, got %s", expected[i], v)
				failed = true
			}
		}

		if failed {
			t.Logf("expected: [%v]", strings.Join(expected, ", "))
			t.Logf("got:      [%v]", strings.Join(l, ", "))
		}

		return nil
	})

	r.Get("/test", handler, handlerMiddleware...)

	svr := httptest.NewServer(r)

	// Make a request to the server
	req, err := http.NewRequest(http.MethodGet, svr.URL+"/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	client := svr.Client()
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
}
