package metrics

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/AlleyBo55/x402-facilitator-oss/internal/domain"
)

// Tracker tracks facilitator performance metrics
type Tracker struct {
	TotalVerifyRequests int64
	SuccessfulVerifies  int64
	FailedVerifies      int64
	TotalSettleRequests int64
	SuccessfulSettles   int64
	FailedSettles       int64
	TotalAmountVerified int64
	RequestLatencySum   int64
	RequestCount        int64
	StartTime           time.Time

	mu             sync.RWMutex
	recentPayments []domain.PaymentRecord
}

// New creates a new metrics tracker
func New() *Tracker {
	return &Tracker{
		StartTime:      time.Now(),
		recentPayments: make([]domain.PaymentRecord, 0, 100),
	}
}

// RecordVerify records a verify request
func (t *Tracker) RecordVerify(success bool, amount int64, latencyMs int64, record domain.PaymentRecord) {
	atomic.AddInt64(&t.TotalVerifyRequests, 1)
	atomic.AddInt64(&t.RequestLatencySum, latencyMs)
	atomic.AddInt64(&t.RequestCount, 1)

	if success {
		atomic.AddInt64(&t.SuccessfulVerifies, 1)
		atomic.AddInt64(&t.TotalAmountVerified, amount)
	} else {
		atomic.AddInt64(&t.FailedVerifies, 1)
	}

	t.mu.Lock()
	t.recentPayments = append(t.recentPayments, record)
	if len(t.recentPayments) > 100 {
		t.recentPayments = t.recentPayments[1:]
	}
	t.mu.Unlock()
}

// RecordSettle records a settle request
func (t *Tracker) RecordSettle(success bool) {
	atomic.AddInt64(&t.TotalSettleRequests, 1)
	if success {
		atomic.AddInt64(&t.SuccessfulSettles, 1)
	} else {
		atomic.AddInt64(&t.FailedSettles, 1)
	}
}

// GetStats returns current statistics as JSON-friendly struct
func (t *Tracker) GetStats() domain.StatsResponse {
	avgLatency := float64(0)
	if t.RequestCount > 0 {
		avgLatency = float64(t.RequestLatencySum) / float64(t.RequestCount)
	}

	t.mu.RLock()
	recent := make([]domain.PaymentRecord, len(t.recentPayments))
	copy(recent, t.recentPayments)
	t.mu.RUnlock()

	// Reverse to show newest first, limit to 20
	for i, j := 0, len(recent)-1; i < j; i, j = i+1, j-1 {
		recent[i], recent[j] = recent[j], recent[i]
	}
	if len(recent) > 20 {
		recent = recent[:20]
	}

	return domain.StatsResponse{
		Uptime:              time.Since(t.StartTime).Round(time.Second).String(),
		TotalVerifyRequests: t.TotalVerifyRequests,
		SuccessfulVerifies:  t.SuccessfulVerifies,
		FailedVerifies:      t.FailedVerifies,
		TotalSettleRequests: t.TotalSettleRequests,
		SuccessfulSettles:   t.SuccessfulSettles,
		FailedSettles:       t.FailedSettles,
		TotalVolumeVerified: t.TotalAmountVerified,
		AvgLatencyMs:        avgLatency,
		RecentPayments:      recent,
	}
}

// GetPrometheus returns Prometheus-formatted metrics
func (t *Tracker) GetPrometheus() string {
	uptime := time.Since(t.StartTime).Seconds()
	avgLatency := float64(0)
	if t.RequestCount > 0 {
		avgLatency = float64(t.RequestLatencySum) / float64(t.RequestCount)
	}

	return fmt.Sprintf(`# HELP x402_uptime_seconds Time since server start
# TYPE x402_uptime_seconds gauge
x402_uptime_seconds %.2f

# HELP x402_verify_total Total verify requests
# TYPE x402_verify_total counter
x402_verify_total{status="success"} %d
x402_verify_total{status="failed"} %d

# HELP x402_settle_total Total settle requests
# TYPE x402_settle_total counter
x402_settle_total{status="success"} %d
x402_settle_total{status="failed"} %d

# HELP x402_amount_verified_lamports Total amount verified
# TYPE x402_amount_verified_lamports counter
x402_amount_verified_lamports %d

# HELP x402_request_latency_avg_ms Average request latency
# TYPE x402_request_latency_avg_ms gauge
x402_request_latency_avg_ms %.2f
`, uptime, t.SuccessfulVerifies, t.FailedVerifies,
		t.SuccessfulSettles, t.FailedSettles,
		t.TotalAmountVerified, avgLatency)
}
