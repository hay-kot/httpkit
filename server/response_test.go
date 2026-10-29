package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJSON_DefaultMarshaler(t *testing.T) {
	w := httptest.NewRecorder()

	err := JSON(w, http.StatusOK, map[string]any{"items": []string(nil)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got, want := w.Body.String(), `{"items":null}`; got != want {
		t.Errorf("body = %s, want %s", got, want)
	}
}

func TestSetJSONMarshaler(t *testing.T) {
	t.Cleanup(func() { SetJSONMarshaler(json.Marshal) })

	SetJSONMarshaler(func(v any) ([]byte, error) {
		return []byte(`{"custom":true}`), nil
	})

	t.Run("JSON", func(t *testing.T) {
		w := httptest.NewRecorder()

		if err := JSON(w, http.StatusOK, "ignored"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got, want := w.Body.String(), `{"custom":true}`; got != want {
			t.Errorf("body = %s, want %s", got, want)
		}
	})

	t.Run("ErrorBuilder", func(t *testing.T) {
		w := httptest.NewRecorder()

		_ = Error().Write(context.Background(), w)

		if got, want := w.Body.String(), `{"custom":true}`; got != want {
			t.Errorf("body = %s, want %s", got, want)
		}
	})
}

func TestSetJSONMarshaler_Error(t *testing.T) {
	t.Cleanup(func() { SetJSONMarshaler(json.Marshal) })

	wantErr := errors.New("marshal failed")
	SetJSONMarshaler(func(v any) ([]byte, error) { return nil, wantErr })

	w := httptest.NewRecorder()

	err := JSON(w, http.StatusOK, "data")
	if !errors.Is(err, wantErr) {
		t.Fatalf("err = %v, want %v", err, wantErr)
	}

	if w.Header().Get("Content-Type") != "" {
		t.Errorf("Content-Type set after marshal failure")
	}
}
