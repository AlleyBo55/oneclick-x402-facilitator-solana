package domain

// ========== x402 Types (Spec-Compliant) ==========
// See: https://x402.gitbook.io/x402/core-concepts/facilitator

// PaymentPayload represents a payment proof from the client
type PaymentPayload struct {
	X402Version int                    `json:"x402Version"`
	Resource    *ResourceInfo          `json:"resource,omitempty"`
	Accepted    *PaymentRequirements   `json:"accepted,omitempty"`
	Payload     map[string]interface{} `json:"payload"`
	Extensions  map[string]interface{} `json:"extensions,omitempty"`
	Scheme      string                 `json:"scheme,omitempty"`
	Network     string                 `json:"network,omitempty"`
}

// ResourceInfo describes the resource being paid for
type ResourceInfo struct {
	URL         string `json:"url"`
	Description string `json:"description"`
	MimeType    string `json:"mimeType"`
}

// PaymentRequirements specifies what the server expects
type PaymentRequirements struct {
	Scheme            string                 `json:"scheme"`
	Network           string                 `json:"network"`
	Asset             string                 `json:"asset,omitempty"`
	Amount            string                 `json:"amount,omitempty"`
	PayTo             string                 `json:"payTo"`
	MaxTimeoutSeconds int                    `json:"maxTimeoutSeconds,omitempty"`
	Extra             map[string]interface{} `json:"extra,omitempty"`
	MaxAmountRequired string                 `json:"maxAmountRequired,omitempty"`
}

// VerifyRequest is the x402 verify request format
type VerifyRequest struct {
	PaymentPayload      PaymentPayload      `json:"paymentPayload"`
	PaymentRequirements PaymentRequirements `json:"paymentRequirements"`
}

// VerifyResponse matches x402 spec
type VerifyResponse struct {
	IsValid       bool   `json:"valid"` // x402 spec uses "valid", not "isValid"
	InvalidReason string `json:"invalidReason,omitempty"`
	Payer         string `json:"payer,omitempty"`
}

// SettleResponse matches x402 spec
type SettleResponse struct {
	Success     bool   `json:"success"`
	ErrorReason string `json:"errorReason,omitempty"`
	Payer       string `json:"payer,omitempty"`
	Transaction string `json:"transaction"`
	Network     string `json:"network"`
}

// SupportedKind represents a supported payment scheme
type SupportedKind struct {
	X402Version int                    `json:"x402Version"`
	Scheme      string                 `json:"scheme"`
	Network     string                 `json:"network"`
	Extra       map[string]interface{} `json:"extra,omitempty"`
}

// SupportedResponse lists all supported payment methods
type SupportedResponse struct {
	Kinds      []SupportedKind     `json:"kinds"`
	Extensions []string            `json:"extensions"`
	Signers    map[string][]string `json:"signers"`
}

// HealthResponse for health checks
type HealthResponse struct {
	Status    string `json:"status"`
	Version   string `json:"version"`
	Network   string `json:"network"`
	Timestamp string `json:"timestamp"`
}

// StatsResponse for /stats API endpoint
type StatsResponse struct {
	Uptime              string          `json:"uptime"`
	TotalVerifyRequests int64           `json:"totalVerifyRequests"`
	SuccessfulVerifies  int64           `json:"successfulVerifies"`
	FailedVerifies      int64           `json:"failedVerifies"`
	TotalSettleRequests int64           `json:"totalSettleRequests"`
	SuccessfulSettles   int64           `json:"successfulSettles"`
	FailedSettles       int64           `json:"failedSettles"`
	TotalVolumeVerified int64           `json:"totalVolumeVerified"`
	AvgLatencyMs        float64         `json:"avgLatencyMs"`
	RecentPayments      []PaymentRecord `json:"recentPayments"`
}

// PaymentRecord for tracking payments
type PaymentRecord struct {
	Timestamp string `json:"timestamp"`
	Signature string `json:"signature"`
	Payer     string `json:"payer"`
	Amount    string `json:"amount"`
	Asset     string `json:"asset"`
	Success   bool   `json:"success"`
	LatencyMs int64  `json:"latencyMs"`
}
