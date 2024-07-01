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
	lastRequest time.Time
	score       uint32
}

type RateLimiter struct {
	mu                  sync.Mutex
	users               map[common.Address]Score
	threshold           uint32
	decay               float64
	maxScore            uint32
	totalRequests       uint64
	rateLimitedRequests uint64
}

// NewRateLimiter creates a new RateLimiter with the specified threshold, decay rate, and maximum score.
// Parameters:
//   - threshold: The maximum score a user can have to be allowed to make a request. If a user's score
//     exceeds this threshold, their request will be denied. Setting the threshold to 0 disables rate limiting.
//   - Decay: The rate at which a user's score decays over time. It represents the amount by which the score
//     decreases per millisecond since the last request. This helps in gradually lowering the score over time.
//   - maxScore: The maximum score a user can accumulate. This prevents the score from growing indefinitely
//     and allows for controlling the upper limit of a user's score.
func NewRateLimiter(threshold uint32, decay float64, maxScore uint32) *RateLimiter {
	rl := &RateLimiter{
		users:     make(map[common.Address]Score),
		threshold: threshold,
		decay:     decay,
		maxScore:  maxScore,
	}
	go rl.logRateLimitedStats()
	return rl
}

// Allow checks if the user is allowed to make a request based on the rate limit threshold
// before comparing to the threshold also decays the score of the user based on the decay rate
func (rl *RateLimiter) Allow(userID common.Address) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.totalRequests++

	// If the threshold is 0, allow all requests (rate limiting is disabled)
	if rl.threshold == 0 {
		return true
	}

	now := time.Now()
	userScore, exists := rl.users[userID]

	if !exists {
		// Create a new entry for the user if not exists
		rl.users[userID] = Score{lastRequest: now, score: 0}
	} else {
		// Decay the score based on the time since the last request
		timeSinceLastRequest := float64(now.Sub(userScore.lastRequest).Milliseconds())

		// calculate the new score by subtracting the decay from the user's score
		newScore := int64(userScore.score) - int64(timeSinceLastRequest*rl.decay)

		// Ensure score does not become negative
		if newScore < 0 {
			newScore = 0
		}

		// Update user's score and last request time if new score is less than the current score
		// use min value of new score and max score to ensure the score does not exceed the max score
		// we check if the score changes before updating to avoid updating time and keeping the score same
		if uint32(newScore) < userScore.score {
			rl.users[userID] = Score{lastRequest: now, score: min(uint32(newScore), rl.maxScore)}
		}
	}

	// Check if the user's score is below the threshold
	if rl.users[userID].score > rl.threshold {
		rl.rateLimitedRequests++
		return false
	}
	return true
}

// UpdateScore updates the score of the user based on the additional score (time taken to process the request)
func (rl *RateLimiter) UpdateScore(userID common.Address, additionalScore uint32) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	newScore := rl.users[userID].score + additionalScore
	rl.users[userID] = Score{lastRequest: time.Now(), score: newScore}
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
