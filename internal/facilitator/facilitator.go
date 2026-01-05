package facilitator

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/AlleyBo55/x402-facilitator-oss/internal/config"
	"github.com/AlleyBo55/x402-facilitator-oss/internal/domain"
	"github.com/AlleyBo55/x402-facilitator-oss/internal/metrics"
	"github.com/AlleyBo55/x402-facilitator-oss/internal/middleware"
	"github.com/AlleyBo55/x402-facilitator-oss/internal/webhook"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

// CAIP-2 Network Identifiers (Solana chain IDs per CAIP-2 spec)
const (
	NetworkDevnet  = "solana:EtWTRABZaYq6iMfeYKouRu166VU2xqa1"
	NetworkMainnet = "solana:5eykt4UsFv8P8NJdTREpY1vzqKqZKvdp"
	DummyFeePayer  = "11111111111111111111111111111111"

	// Default max transaction age (5 minutes)
	DefaultMaxTxAgeSeconds = 300
)

// Facilitator implements the x402 facilitator protocol for Solana
type Facilitator struct {
	rpcClient       *rpc.Client
	network         string
	Metrics         *metrics.Tracker
	RateLimiter     *middleware.RateLimiter
	Webhook         *webhook.Notifier
	APIKey          string
	commitmentLevel rpc.CommitmentType
}

// New creates a new facilitator instance
func New(cfg config.Config) *Facilitator {
	// Use Finalized for mainnet (highest security), Confirmed for devnet (faster)
	commitment := rpc.CommitmentConfirmed
	if cfg.Network == "mainnet-beta" {
		commitment = rpc.CommitmentFinalized
	}

	return &Facilitator{
		rpcClient:       rpc.New(cfg.RPCURL),
		network:         cfg.Network,
		Metrics:         metrics.New(),
		RateLimiter:     middleware.NewRateLimiter(100, time.Minute),
		Webhook:         webhook.New(cfg.WebhookURL),
		APIKey:          cfg.APIKey,
		commitmentLevel: commitment,
	}
}

// GetSupported returns supported payment schemes and networks
// Currently SOL-only for security (no SPL token attack surface)
func (f *Facilitator) GetSupported() domain.SupportedResponse {
	return domain.SupportedResponse{
		Kinds: []domain.SupportedKind{
			{
				X402Version: 2,
				Scheme:      "exact",
				Network:     NetworkDevnet,
				Extra: map[string]interface{}{
					"feePayer": DummyFeePayer,
					"assets":   []string{"SOL"}, // SOL only
				},
			},
			{
				X402Version: 2,
				Scheme:      "exact",
				Network:     NetworkMainnet,
				Extra: map[string]interface{}{
					"feePayer": DummyFeePayer,
					"assets":   []string{"SOL"}, // SOL only
				},
			},
		},
		Extensions: []string{},
		Signers:    map[string][]string{"solana:*": {DummyFeePayer}},
	}
}

// Verify verifies a payment transaction
func (f *Facilitator) Verify(payload domain.PaymentPayload, requirements domain.PaymentRequirements) domain.VerifyResponse {
	start := time.Now()

	// 1. Extract and validate signature
	sig, ok := payload.Payload["signature"].(string)
	if !ok || sig == "" {
		return domain.VerifyResponse{IsValid: false, InvalidReason: "Missing signature"}
	}

	// 2. Validate signature format (base58, 88 chars)
	if len(sig) < 80 || len(sig) > 90 {
		return domain.VerifyResponse{IsValid: false, InvalidReason: "Invalid signature format"}
	}

	// 3. Parse required amount
	amountStr := requirements.Amount
	if amountStr == "" {
		amountStr = requirements.MaxAmountRequired
	}
	if amountStr == "" {
		return domain.VerifyResponse{IsValid: false, InvalidReason: "Missing payment amount"}
	}
	requiredAmount, ok := new(big.Int).SetString(amountStr, 10)
	if !ok || requiredAmount.Sign() <= 0 {
		return domain.VerifyResponse{IsValid: false, InvalidReason: "Invalid payment amount"}
	}

	// 4. Validate payTo address
	if requirements.PayTo == "" {
		return domain.VerifyResponse{IsValid: false, InvalidReason: "Missing payTo address"}
	}
	if _, err := solana.PublicKeyFromBase58(requirements.PayTo); err != nil {
		return domain.VerifyResponse{IsValid: false, InvalidReason: "Invalid payTo address"}
	}

	// 5. Strict asset validation - SOL only
	asset := requirements.Asset
	if asset != "" && asset != "SOL" && asset != "native" {
		return domain.VerifyResponse{IsValid: false, InvalidReason: "Only SOL payments supported"}
	}

	// 6. Fetch transaction from chain
	tx, blockTime, err := f.fetchTransactionWithBlockTime(sig)
	if err != nil {
		return domain.VerifyResponse{IsValid: false, InvalidReason: fmt.Sprintf("Transaction not found or not confirmed: %v", err)}
	}

	// 7. Check transaction age (防止重放攻击)
	maxAge := DefaultMaxTxAgeSeconds
	if requirements.MaxTimeoutSeconds > 0 {
		maxAge = requirements.MaxTimeoutSeconds
	}
	if blockTime > 0 {
		txAge := time.Now().Unix() - blockTime
		if txAge > int64(maxAge) {
			return domain.VerifyResponse{IsValid: false, InvalidReason: fmt.Sprintf("Transaction too old: %ds > %ds", txAge, maxAge)}
		}
	}

	// 8. Parse native SOL transfer
	paidAmount, payer := f.parseNativeTransfer(tx, requirements.PayTo)

	// 9. Check if payment is sufficient
	latencyMs := time.Since(start).Milliseconds()
	record := domain.PaymentRecord{
		Timestamp: time.Now().Format(time.RFC3339),
		Signature: sig,
		Payer:     payer,
		Amount:    paidAmount.String(),
		Asset:     "SOL",
		Success:   paidAmount.Cmp(requiredAmount) >= 0,
		LatencyMs: latencyMs,
	}

	if paidAmount.Cmp(requiredAmount) >= 0 {
		f.Metrics.RecordVerify(true, paidAmount.Int64(), latencyMs, record)
		f.Webhook.Notify("payment.verified", map[string]interface{}{
			"signature": sig, "payer": payer, "amount": paidAmount.String(),
		})
		return domain.VerifyResponse{IsValid: true, Payer: payer}
	}

	f.Metrics.RecordVerify(false, 0, latencyMs, record)
	return domain.VerifyResponse{
		IsValid:       false,
		InvalidReason: fmt.Sprintf("Insufficient: paid %s lamports, required %s", paidAmount, requiredAmount),
		Payer:         payer,
	}
}

// Settle settles a payment (for Solana, verify = settle since tx already on-chain)
func (f *Facilitator) Settle(payload domain.PaymentPayload, requirements domain.PaymentRequirements) domain.SettleResponse {
	result := f.Verify(payload, requirements)
	sig, _ := payload.Payload["signature"].(string)
	network := f.getNetworkCAIP()

	if !result.IsValid {
		f.Metrics.RecordSettle(false)
		return domain.SettleResponse{
			Success: false, ErrorReason: result.InvalidReason,
			Transaction: sig, Network: network,
		}
	}

	f.Metrics.RecordSettle(true)
	f.Webhook.Notify("payment.settled", map[string]interface{}{
		"signature": sig, "payer": result.Payer, "network": network,
	})

	return domain.SettleResponse{
		Success: true, Payer: result.Payer,
		Transaction: sig, Network: network,
	}
}

// GetStats returns current metrics
func (f *Facilitator) GetStats() domain.StatsResponse {
	return f.Metrics.GetStats()
}

// GetPrometheus returns Prometheus metrics
func (f *Facilitator) GetPrometheus() string {
	return f.Metrics.GetPrometheus()
}

func (f *Facilitator) getNetworkCAIP() string {
	if f.network == "mainnet-beta" {
		return NetworkMainnet
	}
	return NetworkDevnet
}

// fetchTransactionWithBlockTime fetches tx and returns block time for age check
func (f *Facilitator) fetchTransactionWithBlockTime(signature string) (*rpc.GetTransactionResult, int64, error) {
	sig, err := solana.SignatureFromBase58(signature)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid signature: %w", err)
	}

	// Exponential backoff: 2s, 4s, 8s, 16s, 32s (~60s total max)
	// Public RPCs like api.devnet.solana.com are load-balanced; different nodes may have different sync states
	for attempt := 0; attempt < 5; attempt++ {
		tx, err := f.rpcClient.GetTransaction(context.Background(), sig, &rpc.GetTransactionOpts{
			Commitment: f.commitmentLevel,
			Encoding:   solana.EncodingBase64,
		})
		if err == nil && tx != nil {
			blockTime := int64(0)
			if tx.BlockTime != nil {
				blockTime = int64(*tx.BlockTime)
			}
			return tx, blockTime, nil
		}
		time.Sleep(time.Duration(2<<attempt) * time.Second) // 2s, 4s, 8s, 16s, 32s
	}
	return nil, 0, fmt.Errorf("transaction not found after 5 retries")
}

// parseNativeTransfer extracts SOL transfer amount to payTo address
func (f *Facilitator) parseNativeTransfer(tx *rpc.GetTransactionResult, payTo string) (*big.Int, string) {
	if tx.Meta == nil || tx.Transaction == nil {
		return big.NewInt(0), ""
	}

	decodedTx, err := tx.Transaction.GetTransaction()
	if err != nil {
		return big.NewInt(0), ""
	}

	var payer string
	paidAmount := big.NewInt(0)

	accountKeys := decodedTx.Message.AccountKeys
	for i, acc := range accountKeys {
		pubkey := acc.String()

		// First signer is the payer/fee payer
		if i == 0 {
			payer = pubkey
		}

		// Check if this account is the recipient
		if pubkey == payTo && i < len(tx.Meta.PreBalances) && i < len(tx.Meta.PostBalances) {
			pre := tx.Meta.PreBalances[i]
			post := tx.Meta.PostBalances[i]
			if post > pre {
				// Positive balance change = received lamports
				paidAmount = big.NewInt(int64(post - pre))
			}
		}
	}

	return paidAmount, payer
}
