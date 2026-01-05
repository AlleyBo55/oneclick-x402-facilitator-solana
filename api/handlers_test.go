package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlleyBo55/x402-facilitator-oss/internal/config"
	"github.com/AlleyBo55/x402-facilitator-oss/internal/domain"
	"github.com/AlleyBo55/x402-facilitator-oss/internal/facilitator"
)

func TestSupportedEndpoint(t *testing.T) {
	// Setup
	cfg := config.Config{Network: "devnet"}
	f := facilitator.New(cfg) // Real facilitator, but we won't call RPC for /supported

	// Register routes
	// Need to register routes similar to RegisterRoutes but on our mux
	// To avoid circular dep issues or global state, we'll just test the handler function logic via RegisterRoutes
	// Ideally RegisterRoutes accepts a *http.ServeMux, but it uses http.HandleFunc (DefaultServeMux)
	// For this test, we'll swap DefaultServeMux temporarily or just test the logic if we refactored.

	// Since RegisterRoutes uses global http.HandleFunc, we can use httptest.NewServer with http.DefaultServeMux
	// BUT concurrent tests might clash. For safety, let's just create a new facilitator and call GetSupported directly for this unit test
	// OR, we can mock the handler logic:

	// Better approach for Clean Arch: Test the handler function if it was exported, or logic.
	// Since we want E2E, let's assume we can run it.

	// Let's rely on the fact that we can call RegisterRoutes which modifies http.DefaultServeMux
	// We'll reset it or just add to it.

	mux := http.NewServeMux()
	RegisterRoutes(mux, f, "devnet")

	// Create Request
	req := httptest.NewRequest("POST", "/supported", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", w.Code)
	}

	var resp domain.SupportedResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(resp.Kinds) == 0 {
		t.Error("Expected at least one supported kind")
	}
	if resp.Kinds[0].Network != "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1" {
		t.Errorf("Expected devnet CAIP, got %s", resp.Kinds[0].Network)
	}
}

func TestHealthEndpoint(t *testing.T) {
	cfg := config.Config{Network: "devnet"}
	f := facilitator.New(cfg)

	mux := http.NewServeMux()
	RegisterRoutes(mux, f, "devnet")

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", w.Code)
	}
}
