package limiter

import (
	"context"
	"time"

	"github.com/voldemarq/sentinel/internal/storage"
	"github.com/voldemarq/sentinel/internal/telemetry"
)

// Result represents the result of a rate limit check
type Result struct {
	// Allowed indicates whether the request is allowed
	Allowed bool

	// Remaining is the number of requests remaining in the current window
	Remaining int64

	// ResetAt is when the rate limit window resets
	ResetAt time.Time

	// RetryAfter is the suggested time to wait before retrying
	RetryAfter time.Duration
}

// RateLimiter defines the interface for rate limiting algorithms
type RateLimiter interface {
	// Allow checks if a request is allowed and updates counters
	// The key uniquely identifies the rate limited entity (e.g., user ID, IP address)
	// Limit is the maximum number of requests allowed in the window
	// Window is the time window for the rate limit
	Allow(ctx context.Context, key string, limit int64, window time.Duration) (*Result, error)

	// Status returns the current rate limit status without updating counters
	Status(ctx context.Context, key string, limit int64, window time.Duration) (*Result, error)
}

// BaseRateLimiter provides common functionality for rate limiters
type BaseRateLimiter struct {
	store   storage.Storage
	metrics *telemetry.Metrics
}

// newBaseRateLimiter creates a new base rate limiter
func newBaseRateLimiter(store storage.Storage, metrics *telemetry.Metrics) BaseRateLimiter {
	return BaseRateLimiter{
		store:   store,
		metrics: metrics,
	}
}

// recordMetrics records rate limiting metrics
func (b *BaseRateLimiter) recordMetrics(key string, allowed bool) {
	if b.metrics == nil {
		return
	}

	if allowed {
		b.metrics.AllowedRequests.WithLabelValues(key).Inc()
	} else {
		b.metrics.BlockedRequests.WithLabelValues(key).Inc()
	}
}
