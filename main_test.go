package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"urlx/shortener"
	"urlx/store"

	"github.com/gin-gonic/gin"
)

func setupRouter(s store.Store, shortenerSvc shortener.Shortener) *gin.Engine {
	r := gin.Default()

	r.POST("/shorten", func(c *gin.Context) {
		var req ShortenRequest
		if err := c.ShouldBindJSON(&req); err != nil || req.URL == "" {
			c.Status(400)
			return
		}
		short := shortenerSvc.Shorten(req.URL)
		s.Set(short, req.URL)
		c.JSON(200, ShortenResponse{Short: short})
	})

	r.GET("/:short", func(c *gin.Context) {
		short := c.Param("short")
		if long, ok := s.Get(short); ok {
			c.Redirect(302, long)
			return
		}
		c.Status(404)
	})

	r.DELETE("/delete/:short", func(c *gin.Context) {
		short := c.Param("short")
		s.Delete(short)
		c.Status(204)
	})

	return r
}

func TestShortenAndRedirect(t *testing.T) {
	gin.SetMode(gin.TestMode)
	s := store.NewMemoryStore()
	shortenerSvc := shortener.NewSimpleShortener()
	r := setupRouter(s, shortenerSvc)

	// Test shorten
	body := bytes.NewBufferString(`{"url":"https://example.com"}`)
	req, _ := http.NewRequest(http.MethodPost, "/shorten", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp ShortenResponse
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
	r := setupRouter(s, shortenerSvc)
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
