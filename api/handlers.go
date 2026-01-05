package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/AlleyBo55/x402-facilitator-oss/internal/domain"
	"github.com/AlleyBo55/x402-facilitator-oss/internal/facilitator"
	"github.com/AlleyBo55/x402-facilitator-oss/internal/middleware"
	"github.com/AlleyBo55/x402-facilitator-oss/internal/tokens"
)

const version = "2.0.0"

// RegisterRoutes registers all HTTP routes
func RegisterRoutes(mux *http.ServeMux, f *facilitator.Facilitator, network string) {
	// Health check
	mux.HandleFunc("/health", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {
		middleware.JSON(w, http.StatusOK, domain.HealthResponse{
			Status: "ok", Version: version, Network: network,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
	}))

	// Root - API info
	mux.HandleFunc("/", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		middleware.JSON(w, http.StatusOK, map[string]interface{}{
			"name":    "x402 Facilitator (Open Source)",
			"version": version,
			"network": network,
			"assets":  tokens.GetSymbols(network),
			"endpoints": map[string]string{
				"POST /verify":    "Verify a payment",
				"POST /settle":    "Settle a payment",
				"POST /supported": "Get supported schemes",
				"GET /health":     "Health check",
				"GET /metrics":    "Prometheus metrics",
				"GET /stats":      "JSON statistics",
			},
		})
	}))

	// Public endpoints
	mux.HandleFunc("/supported", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			middleware.JSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
			return
		}
		middleware.JSON(w, http.StatusOK, f.GetSupported())
	}))

	mux.HandleFunc("/metrics", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(f.GetPrometheus()))
	}))

	mux.HandleFunc("/stats", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {
		middleware.JSON(w, http.StatusOK, f.GetStats())
	}))

	// Public endpoints (no auth, like PayAI)
	mux.HandleFunc("/verify", middleware.CORS(
		middleware.RateLimit(f.RateLimiter, func(w http.ResponseWriter, r *http.Request) {
			handleVerify(f, w, r)
		})))

	mux.HandleFunc("/settle", middleware.CORS(
		middleware.RateLimit(f.RateLimiter, func(w http.ResponseWriter, r *http.Request) {
			handleSettle(f, w, r)
		})))
}

func handleVerify(f *facilitator.Facilitator, w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		middleware.JSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	var req domain.VerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.JSON(w, http.StatusBadRequest, domain.VerifyResponse{IsValid: false, InvalidReason: "Invalid request"})
		return
	}

	result := f.Verify(req.PaymentPayload, req.PaymentRequirements)
	status := http.StatusOK
	if !result.IsValid {
		status = http.StatusBadRequest
	}
	middleware.JSON(w, status, result)
}

func handleSettle(f *facilitator.Facilitator, w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		middleware.JSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

	var req domain.VerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.JSON(w, http.StatusBadRequest, domain.SettleResponse{Success: false, ErrorReason: "Invalid request"})
		return
	}

	result := f.Settle(req.PaymentPayload, req.PaymentRequirements)
	status := http.StatusOK
	if !result.Success {
		status = http.StatusBadRequest
	}
	middleware.JSON(w, status, result)
}
