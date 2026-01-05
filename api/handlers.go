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
func RegisterRoutes(f *facilitator.Facilitator, network string) {
	// Health check
	http.HandleFunc("/health", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {
		middleware.JSON(w, http.StatusOK, domain.HealthResponse{
			Status: "ok", Version: version, Network: network,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
	}))

	// Root - API info
	http.HandleFunc("/", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {
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
	http.HandleFunc("/supported", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			middleware.JSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
			return
		}
		middleware.JSON(w, http.StatusOK, f.GetSupported())
	}))

	http.HandleFunc("/metrics", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(f.GetPrometheus()))
	}))

	http.HandleFunc("/stats", middleware.CORS(func(w http.ResponseWriter, r *http.Request) {
		middleware.JSON(w, http.StatusOK, f.GetStats())
	}))

	// Protected endpoints
	http.HandleFunc("/verify", middleware.CORS(
		middleware.RateLimit(f.RateLimiter,
			middleware.Auth(f.APIKey, func(w http.ResponseWriter, r *http.Request) {
				handleVerify(f, w, r)
			}))))

	http.HandleFunc("/settle", middleware.CORS(
		middleware.RateLimit(f.RateLimiter,
			middleware.Auth(f.APIKey, func(w http.ResponseWriter, r *http.Request) {
				handleSettle(f, w, r)
			}))))
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
