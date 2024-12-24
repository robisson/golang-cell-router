package router

import (
	"cell-router/config"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRouter(t *testing.T) {
	router := NewRouter()
	if router.config == nil {
		t.Error("Expected config to be loaded, but it was nil")
	}
}

func TestServeHTTP(t *testing.T) {
	cfg := &config.Config{
		Cells: []config.CellConfig{
			{Name: "Cell1", Endpoint: "http://example.com", RangeFrom: 1, RangeTo: 100},
		},
	}
	router := &Router{config: cfg}

	req := httptest.NewRequest("GET", "/?client_id=50", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}
}

func TestServeHTTPInvalidClientID(t *testing.T) {
	router := NewRouter()

	req := httptest.NewRequest("GET", "/?client_id=invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status BadRequest, got %v", resp.Status)
	}
}

func TestServeHTTPNoMatchingCell(t *testing.T) {
	cfg := &config.Config{
		Cells: []config.CellConfig{
			{Name: "Cell1", Endpoint: "http://example.com", RangeFrom: 1, RangeTo: 100},
		},
	}
	router := &Router{config: cfg}

	req := httptest.NewRequest("GET", "/?client_id=200", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status NotFound, got %v", resp.Status)
	}
}

func TestServeHTTPProxyError(t *testing.T) {
	cfg := &config.Config{
		Cells: []config.CellConfig{
			{Name: "Cell1", Endpoint: "http://invalid-url", RangeFrom: 1, RangeTo: 100},
		},
	}
	router := &Router{config: cfg}

	req := httptest.NewRequest("GET", "/?client_id=50", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status InternalServerError, got %v", resp.Status)
	}
}

func TestServeHTTPProxyRequestCreationError(t *testing.T) {
	cfg := &config.Config{
		Cells: []config.CellConfig{
			{
				Name:      "Cell1",
				Endpoint:  ":/invalid-url", // Invalid URL format that will fail parsing
				RangeFrom: 1,
				RangeTo:   100,
			},
		},
	}
	router := &Router{config: cfg}

	req := httptest.NewRequest("GET", "/?client_id=50", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status InternalServerError, got %v", resp.Status)
	}
}
