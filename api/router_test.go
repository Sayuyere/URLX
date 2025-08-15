package api

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"urlx/logging"
	"urlx/shortener"
	"urlx/store"
)

func TestUIRootServesHTML(t *testing.T) {
	// Create a temp HTML file in ui/index.html
	os.MkdirAll("ui", 0755)
	tmpHtml := []byte("<html><body>Test UI</body></html>")
	err := ioutil.WriteFile("ui/index.html", tmpHtml, 0644)
	if err != nil {
		t.Fatalf("failed to create ui/index.html: %v", err)
	}
	defer os.RemoveAll("ui")

	s := store.NewMemoryStore()
	shortenerSvc := shortener.NewSimpleShortener()
	logger := logging.NewLogger()
	r := SetupRouter(s, shortenerSvc, logger)

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if got := w.Body.String(); got != string(tmpHtml) {
		t.Fatalf("expected body %q, got %q", string(tmpHtml), got)
	}
}

func TestUIRootNotFound(t *testing.T) {
	os.RemoveAll("ui") // Ensure ui/index.html does not exist

	s := store.NewMemoryStore()
	shortenerSvc := shortener.NewSimpleShortener()
	logger := logging.NewLogger()
	r := SetupRouter(s, shortenerSvc, logger)

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
	if got := w.Body.String(); got != "UI not found" {
		t.Fatalf("expected body 'UI not found', got %q", got)
	}
}
