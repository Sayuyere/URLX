package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"urlx/api"
	"urlx/shortener"
	"urlx/store"

	"github.com/gin-gonic/gin"
)

func TestShortenAndRedirect(t *testing.T) {
	gin.SetMode(gin.TestMode)
	s := store.NewMemoryStore()
	shortenerSvc := shortener.NewSimpleShortener()
	r := api.SetupRouter(s, shortenerSvc)

	// Test shorten
	body := bytes.NewBufferString(`{"url":"https://example.com"}`)
	req, _ := http.NewRequest(http.MethodPost, "/shorten", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp api.ShortenResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil || resp.Short == "" {
		t.Fatalf("invalid response: %v", w.Body.String())
	}

	// Test redirect
	req2, _ := http.NewRequest(http.MethodGet, "/"+resp.Short, nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	if w2.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d", w2.Code)
	}
	if loc := w2.Header().Get("Location"); loc != "https://example.com" {
		t.Fatalf("expected redirect to https://example.com, got %s", loc)
	}
}

func TestDelete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	s := store.NewMemoryStore()
	shortenerSvc := shortener.NewSimpleShortener()
	r := api.SetupRouter(s, shortenerSvc)
	short := "abc123"
	s.Set(short, "https://example.com")

	req, _ := http.NewRequest(http.MethodDelete, "/delete/"+short, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
	if _, ok := s.Get(short); ok {
		t.Fatalf("expected url to be deleted")
	}
}
