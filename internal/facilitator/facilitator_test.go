package facilitator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlleyBo55/x402-facilitator-oss/internal/config"
	"github.com/AlleyBo55/x402-facilitator-oss/internal/domain"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
)

// Helper to generate a dummy signed transaction
func generateBase64Transaction(payer, recipient solana.PrivateKey) (string, string) {
	// Create instruction: Payer -> Recipient (1000000 lamports)
	inst := system.NewTransferInstruction(
		1000000,
		payer.PublicKey(),
		recipient.PublicKey(),
	).Build()

	// Build transaction
	tx, err := solana.NewTransaction(
		[]solana.Instruction{inst},
		solana.Hash{}, // Recent blockhash (dummy)
		solana.TransactionPayer(payer.PublicKey()),
	)
	if err != nil {
		panic(err)
	}

	// Sign
	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if key.Equals(payer.PublicKey()) {
				return &payer
			}
			return nil
		},
	)
	if err != nil {
		panic(err)
	}

	// Serialize to Base64
	return tx.MustToBase64(), tx.Signatures[0].String()
}

func TestVerifyPayment_Success(t *testing.T) {
	// 0. Generate Real Keys and Transaction
	payer := solana.NewWallet().PrivateKey
	recipient := solana.NewWallet().PrivateKey
	base64Tx, sig := generateBase64Transaction(payer, recipient)

	// 1. Setup Mock RPC
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Method string      `json:"method"`
			ID     interface{} `json:"id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if req.Method == "getTransaction" {
			idBytes, _ := json.Marshal(req.ID)
			// Return Base64 transaction
			response := fmt.Sprintf(`{
				"jsonrpc": "2.0",
				"result": {
					"slot": 100,
					"blockTime": 1600000000,
					"meta": {
						"preBalances": [2000000, 0, 1],
						"postBalances": [995000, 1000000, 1],
						"err": null
					},
					"transaction": ["%s", "base64"]
				},
				"id": %s
			}`, base64Tx, string(idBytes))
			w.Write([]byte(response))
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	// 2. Init Facilitator
	cfg := config.Config{
		RPCURL:  server.URL,
		Network: "devnet",
	}
	f := New(cfg)

	// 3. Define Test Data
	payload := domain.PaymentPayload{
		X402Version: 2,
		Payload: map[string]interface{}{
			"signature": sig,
		},
	}

	reqs := domain.PaymentRequirements{
		PayTo:             recipient.PublicKey().String(),
		Amount:            "1000000",
		Asset:             "SOL",
		MaxTimeoutSeconds: 9999999999,
	}

	// 4. Run Verify
	result := f.Verify(payload, reqs)

	if !result.IsValid {
		t.Errorf("Expected valid payment, got invalid: %s", result.InvalidReason)
	}

	if result.Payer != payer.PublicKey().String() {
		t.Errorf("Expected payer %s, got %s", payer.PublicKey().String(), result.Payer)
	}
}

func TestVerifyPayment_InsufficientAmount(t *testing.T) {
	// Re-use generation logic or simple mock for failure case
	payer := solana.NewWallet().PrivateKey
	recipient := solana.NewWallet().PrivateKey
	base64Tx, sig := generateBase64Transaction(payer, recipient)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Method string      `json:"method"`
			ID     interface{} `json:"id"`
		}
		json.NewDecoder(r.Body).Decode(&req)
		idBytes, _ := json.Marshal(req.ID)

		if req.Method == "getTransaction" {
			response := fmt.Sprintf(`{
				"jsonrpc": "2.0",
				"result": {
					"slot": 100,
					"blockTime": 1600000000,
					"meta": {
						"preBalances": [2000000, 0, 1],
						"postBalances": [995000, 1000000, 1],
						"err": null
					},
					"transaction": ["%s", "base64"]
				},
				"id": %s
			}`, base64Tx, string(idBytes))
			w.Write([]byte(response))
			return
		}
	}))
	defer server.Close()

	cfg := config.Config{RPCURL: server.URL}
	f := New(cfg)

	payload := domain.PaymentPayload{
		Payload: map[string]interface{}{"signature": sig},
	}
	// Require MORE than 1000000
	reqs := domain.PaymentRequirements{
		PayTo:             recipient.PublicKey().String(),
		Amount:            "2000000",
		Asset:             "SOL",
		MaxTimeoutSeconds: 9999999999,
	}

	result := f.Verify(payload, reqs)

	if result.IsValid {
		t.Error("Expected invalid payment (insufficient amount), got valid")
	}
}
