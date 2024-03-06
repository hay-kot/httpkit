package errchain

import (
	"net/http"
	"net/http/httptest"
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
