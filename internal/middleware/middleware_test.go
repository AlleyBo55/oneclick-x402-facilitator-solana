package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
	// Allow 2 requests per second
	rl := NewRateLimiter(2, time.Second)

	ip := "127.0.0.1"

	// Request 1: Should pass
	if !rl.Allow(ip) {
		t.Error("Request 1 should be allowed")
	}

	// Request 2: Should pass
	if !rl.Allow(ip) {
		t.Error("Request 2 should be allowed")
	}

	// Request 3: Should fail
	if rl.Allow(ip) {
		t.Error("Request 3 should be denied")
	}
}

func TestMiddlewareChain(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	// Test CORS
	req := httptest.NewRequest("OPTIONS", "/", nil)
	w := httptest.NewRecorder()

	CORS(handler)(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("CORS header missing")
	}

	// Test Auth - Success
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-API-Key", "secret")
	w = httptest.NewRecorder()

	Auth("secret", handler)(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", w.Code)
	}

	// Test Auth - Failure
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-API-Key", "wrong")
	w = httptest.NewRecorder()

	Auth("secret", handler)(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 Unauthorized, got %d", w.Code)
	}
}
