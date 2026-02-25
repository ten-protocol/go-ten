package ratelimiter

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	gethlog "github.com/ethereum/go-ethereum/log"
	"golang.org/x/time/rate"
)

// HTTPRateLimiter provides two-layer rate limiting for HTTP endpoints:
// 1. Global rate limit - limits total requests per second across all clients
// 2. Per-IP rate limit - limits requests per second from a single IP address
type HTTPRateLimiter struct {
	globalLimiter *rate.Limiter
	perIPLimiters map[string]*ipLimiter // map of IP -> limiter
	perIPRate     rate.Limit
	perIPBurst    int
	mu            sync.RWMutex
	logger        gethlog.Logger
	logLimiter    *rate.Limiter
	cleanupTicker *time.Ticker
	stopCleanup   chan struct{}
}

// ipLimiter holds a rate limiter and last seen time for a specific IP
type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

const (
	// cleanupInterval is how often we check for stale IP entries
	// TODO: determine good values for interval and threshold
	cleanupInterval = 5 * time.Minute
	// staleThreshold is how long an IP entry can be idle before removal
	staleThreshold = 10 * time.Minute
	// logRateLimit limits rate limit warning logs to 1 per second
	logRateLimit = 1.0
)

// NewHTTPRateLimiter creates a new HTTP rate limiter.
// Set globalRate to 0 to disable global rate limiting.
// Set perIPRate to 0 to disable per-IP rate limiting.
// If both are 0, all requests are allowed (same functionality as before).
// Burst is automatically set to 1.5x the rate for both limiters.
func NewHTTPRateLimiter(globalRate float64, perIPRate float64, logger gethlog.Logger) *HTTPRateLimiter {
	// Calculate per-IP burst as 1.5x the rate
	perIPBurst := int(perIPRate * 1.5)
	if perIPBurst < 1 && perIPRate > 0 {
		perIPBurst = 1
	}

	rl := &HTTPRateLimiter{
		perIPLimiters: make(map[string]*ipLimiter),
		perIPRate:     rate.Limit(perIPRate),
		perIPBurst:    perIPBurst,
		logger:        logger,
		logLimiter:    rate.NewLimiter(logRateLimit, 1),
		stopCleanup:   make(chan struct{}),
	}

	// Create global limiter if enabled (burst = 1.5x rate)
	var globalBurst int
	if globalRate > 0 {
		globalBurst = int(globalRate * 1.5)
		if globalBurst < 1 {
			globalBurst = 1
		}
		rl.globalLimiter = rate.NewLimiter(rate.Limit(globalRate), globalBurst)
	}

	// Start cleanup goroutine if per-IP limiting is enabled
	if perIPRate > 0 {
		rl.startCleanup()
	}

	logger.Info("HTTP rate limiter initialized",
		"globalRate", globalRate,
		"globalBurst", globalBurst,
		"perIPRate", perIPRate,
		"perIPBurst", perIPBurst,
		"globalEnabled", globalRate > 0,
		"perIPEnabled", perIPRate > 0)

	return rl
}

// Allow checks if a request from the given IP should be allowed.
// Returns (allowed, retryAfter) where retryAfter is the suggested wait time if not allowed.
func (rl *HTTPRateLimiter) Allow(ip string) (bool, time.Duration) {
	// If both limiters are disabled, always allow
	if rl.globalLimiter == nil && rl.perIPRate <= 0 {
		return true, 0
	}

	// Check global limit first (if enabled)
	if rl.globalLimiter != nil {
		if !rl.globalLimiter.Allow() {
			rl.logRateLimited("Global rate limit exceeded")
			return false, rl.calculateRetryAfter(rl.globalLimiter)
		}
	}

	// Check per-IP limit (if enabled)
	if rl.perIPRate > 0 {
		limiter := rl.getOrCreateIPLimiter(ip)
		if !limiter.Allow() {
			rl.logRateLimited("Per-IP rate limit exceeded", "ip", ip)
			return false, rl.calculateRetryAfter(limiter)
		}
	}

	return true, 0
}

// getOrCreateIPLimiter returns the rate limiter for the given IP, creating one if needed.
func (rl *HTTPRateLimiter) getOrCreateIPLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if ipLim, exists := rl.perIPLimiters[ip]; exists {
		ipLim.lastSeen = time.Now()
		return ipLim.limiter
	}

	// Create new limiter for this IP
	limiter := rate.NewLimiter(rl.perIPRate, rl.perIPBurst)
	rl.perIPLimiters[ip] = &ipLimiter{
		limiter:  limiter,
		lastSeen: time.Now(),
	}
	return limiter
}

// calculateRetryAfter estimates when the next request would be allowed.
func (rl *HTTPRateLimiter) calculateRetryAfter(limiter *rate.Limiter) time.Duration {
	// Reserve a token and immediately cancel to get the delay
	reservation := limiter.Reserve()
	delay := reservation.Delay()
	reservation.Cancel()

	// Minimum 1 second retry-after
	if delay < time.Second {
		delay = time.Second
	}
	return delay
}

// logRateLimited logs a rate limit event, but rate-limits the logging itself.
func (rl *HTTPRateLimiter) logRateLimited(msg string, ctx ...interface{}) {
	if rl.logLimiter.Allow() {
		rl.logger.Warn(msg, ctx...)
	}
}

// startCleanup starts the background goroutine that removes stale IP entries.
func (rl *HTTPRateLimiter) startCleanup() {
	rl.cleanupTicker = time.NewTicker(cleanupInterval)
	go func() {
		for {
			select {
			case <-rl.cleanupTicker.C:
				rl.cleanupStaleEntries()
			case <-rl.stopCleanup:
				rl.cleanupTicker.Stop()
				return
			}
		}
	}()
}

// cleanupStaleEntries removes IP entries that haven't been seen recently.
func (rl *HTTPRateLimiter) cleanupStaleEntries() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	cutoff := time.Now().Add(-staleThreshold)
	removedCount := 0
	for ip, ipLim := range rl.perIPLimiters {
		if ipLim.lastSeen.Before(cutoff) {
			delete(rl.perIPLimiters, ip)
			removedCount++
		}
	}

	if removedCount > 0 {
		rl.logger.Debug("Cleaned up stale IP rate limiters", "removed", removedCount, "remaining", len(rl.perIPLimiters))
	}
}

// Stop stops the cleanup goroutine. Call this when shutting down.
func (rl *HTTPRateLimiter) Stop() {
	if rl.stopCleanup != nil {
		close(rl.stopCleanup)
	}
}

// GetClientIP extracts the client IP from an HTTP request.
// It respects X-Forwarded-For header for requests behind a proxy.
func GetClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (set by proxies/load balancers)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For can contain multiple IPs: "client, proxy1, proxy2"
		// The first one is the original client IP
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			ip := strings.TrimSpace(parts[0])
			if ip != "" {
				return ip
			}
		}
	}

	// Check X-Real-IP header (alternative header used by some proxies)
	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		return strings.TrimSpace(xrip)
	}

	// Fall back to RemoteAddr
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		// RemoteAddr might not have a port
		return r.RemoteAddr
	}
	return host
}

// IsEnabled returns true if any rate limiting is enabled.
func (rl *HTTPRateLimiter) IsEnabled() bool {
	return rl.globalLimiter != nil || rl.perIPRate > 0
}
