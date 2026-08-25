package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthEndpoint(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("health: got %d, want 200", w.Code)
	}
	if !strings.Contains(w.Body.String(), `"status":"ok"`) {
		t.Fatalf("health body: %s", w.Body.String())
	}
}

func TestAbsorbanceEndpoint(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	body := `{"components":[{"label":"dye","extinction":100,"concentration":0.01}],"path_length":1.0,"stray_fraction":0}`
	req := httptest.NewRequest(http.MethodPost, "/api/absorbance", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("absorbance: got %d, want 200; body: %s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "Absorbance") && !strings.Contains(w.Body.String(), "absorbance") {
		t.Logf("response: %s", w.Body.String())
	}
}

func TestAbsorbanceBadMethod(t *testing.T) {
	mux := New(Config{Addr: ":8080"})
	req := httptest.NewRequest(http.MethodGet, "/api/absorbance", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("got %d, want 405", w.Code)
	}
}
