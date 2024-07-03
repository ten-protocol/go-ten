package ratelimiter

import (
	"log"
	"math"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

// Package ratelimiter implements a simple rate limiting mechanism
// using a score-based approach. Each user has a score that decays
// over time. Requests are allowed if the user's score is below a
// specified threshold. The score increases with each request and
// decays based on the time since the last request.

type Score struct {
	lastRequest        time.Time
	score              uint32
	concurrentRequests uint32
}

type RateLimiter struct {
	mu                    sync.Mutex
	users                 map[common.Address]Score
	threshold             uint32
	decay                 float64
	maxConcurrentRequests uint32
	totalRequests         uint64
	rateLimitedRequests   uint64
}

// NewRateLimiter creates a new RateLimiter with the specified threshold, decay rate, and maximum score.
// Parameters:
//   - threshold: The maximum score a user can have to be allowed to make a request. If a user's score
//     exceeds this threshold, their request will be denied. Setting the threshold to 0 disables rate limiting.
//   - Decay: The rate at which a user's score decays over time. It represents the amount by which the score
//     has decreased per millisecond since the last request. This helps in gradually lowering the score over time.
//   - concurrentRequestsLimit: The maximum number of concurrent requests a user can make. If a user exceeds this
//     limit, their request will be denied.
func NewRateLimiter(threshold uint32, decay float64, concurrentRequestsLimit uint32) *RateLimiter {
	rl := &RateLimiter{
		users:                 make(map[common.Address]Score),
		threshold:             threshold,
		decay:                 decay,
		maxConcurrentRequests: concurrentRequestsLimit,
	}
	go rl.logRateLimitedStats()
	return rl
}

// Allow checks if the user is allowed to make a request based on the rate limit threshold
// before comparing to the threshold also decays the score of the user based on the decay rate
func (rl *RateLimiter) Allow(userID common.Address) bool {
	// If the threshold is 0, allow all requests (rate limiting is disabled)
	if rl.threshold == 0 {
		return true
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.totalRequests++
	now := time.Now()
	userScore, exists := rl.users[userID]

	if !exists {
		// Create a new entry for the user if not exists
		rl.users[userID] = Score{lastRequest: now, score: 0, concurrentRequests: 1}
		return true
	} else {
		// Increase the number of concurrent requests by 1
		rl.users[userID] = Score{lastRequest: now, score: userScore.score, concurrentRequests: userScore.concurrentRequests + 1}
	}

	// Check if the user's score is below the threshold (taking into account the decay)
	timeSinceLastRequest := float64(now.Sub(userScore.lastRequest).Milliseconds())
	newScore := int64(userScore.score) - int64(timeSinceLastRequest*rl.decay)
	if newScore > int64(rl.threshold) || userScore.concurrentRequests > rl.maxConcurrentRequests {
		rl.rateLimitedRequests++
		return false
	}
	return true
}

// UpdateScore updates the score of the user based on the execution duration of the request.
func (rl *RateLimiter) UpdateScore(userID common.Address, executionDuration uint32) {
	// If the threshold is 0,rate limiting is disabled and there is no need for updating the score
	if rl.threshold == 0 {
		return
	}
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Decay the score of the user based on the time since the last request
	now := time.Now()
	userScore, exists := rl.users[userID]
	scoreDecay := uint32(0)
	// if user exists decay the score based on the time since the last request
	if exists {
		// Decay the score based on the time since the last request and the decay rate
		timeSinceLastRequest := float64(now.Sub(userScore.lastRequest).Milliseconds())
		// limit the decay to the user's current score
		scoreDecay = min(uint32(timeSinceLastRequest*rl.decay), userScore.score)
	}

	// Increase the score of the user based on the execution duration of the request
	// and update the last request time
	newScore := rl.users[userID].score + executionDuration - scoreDecay
	rl.users[userID] = Score{lastRequest: now, score: newScore, concurrentRequests: userScore.concurrentRequests - 1}
}

func (rl *RateLimiter) logRateLimitedStats() {
	for {
		time.Sleep(30 * time.Minute)
		rl.mu.Lock()
		totalRequests := rl.totalRequests
		rateLimitedRequests := rl.rateLimitedRequests
		rl.totalRequests = 0
		rl.rateLimitedRequests = 0
		rl.mu.Unlock()

		rateLimitedPercentage := float64(rateLimitedRequests) / float64(totalRequests) * 100
		if math.IsNaN(rateLimitedPercentage) {
			rateLimitedPercentage = 0
		}
		log.Printf("Total requests: %d, Rate-limited requests: %d (%.4f%%)", totalRequests, rateLimitedRequests, rateLimitedPercentage)
	}
}
