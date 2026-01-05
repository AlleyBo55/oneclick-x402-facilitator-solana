package metrics

import (
	"sync"
	"testing"
	"time"

	"github.com/AlleyBo55/x402-facilitator-oss/internal/domain"
)

func TestMetricsConcurrency(t *testing.T) {
	tracker := New()
	var wg sync.WaitGroup

	// Simulate concurrent requests
	interactions := 100
	wg.Add(interactions * 2)

	for i := 0; i < interactions; i++ {
		go func() {
			defer wg.Done()
			tracker.RecordVerify(true, 1000, 10, domain.PaymentRecord{
				Signature: "sig",
				Success:   true,
			})
		}()

		go func() {
			defer wg.Done()
			tracker.RecordSettle(false)
		}()
	}

	wg.Wait()

	stats := tracker.GetStats()

	if stats.TotalVerifyRequests != int64(interactions) {
		t.Errorf("Expected %d verify requests, got %d", interactions, stats.TotalVerifyRequests)
	}

	if stats.SuccessfulVerifies != int64(interactions) {
		t.Errorf("Expected %d successful verifies, got %d", interactions, stats.SuccessfulVerifies)
	}

	if stats.TotalSettleRequests != int64(interactions) {
		t.Errorf("Expected %d settle requests, got %d", interactions, stats.TotalSettleRequests)
	}
}

func TestRecentPaymentsLimit(t *testing.T) {
	tracker := New()

	// Add 150 payments
	for i := 0; i < 150; i++ {
		tracker.RecordVerify(true, 100, 1, domain.PaymentRecord{
			Signature: "sig",
			Timestamp: time.Now().String(),
		})
	}

	stats := tracker.GetStats()

	// Should only return top 20
	if len(stats.RecentPayments) > 20 {
		t.Errorf("Expected max 20 recent payments, got %d", len(stats.RecentPayments))
	}
}
