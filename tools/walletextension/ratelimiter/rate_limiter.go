package ratelimiter

import (
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type Score struct {
	lastRequest time.Time
	score       uint32
}

type RateLimiter struct {
	mu        sync.Mutex
	users     map[common.Address]Score
	threshold uint32
	decay     uint32
}

func NewRateLimiter(threshold uint32, decay uint32) *RateLimiter {
	return &RateLimiter{
		users:     make(map[common.Address]Score),
		threshold: threshold,
		decay:     decay,
	}
}

// Allow checks if the user is allowed to make a request based on the rate limit threshold
// before comparing to the threshold also decays the score of the user based on the decay rate
func (rl *RateLimiter) Allow(userID common.Address) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

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
		timeSinceLastRequest := uint32(now.Sub(userScore.lastRequest).Milliseconds())

		// calculate the new score by subtracting the decay from the user's score
		newScore := int64(userScore.score) - int64(timeSinceLastRequest)*int64(rl.decay)

		// Ensure score does not become negative
		if newScore < 0 {
			newScore = 0
		}

		// Update user's score and last request time
		rl.users[userID] = Score{lastRequest: now, score: uint32(newScore)}
	}

	// Check if the user's score is below the threshold
	return rl.users[userID].score <= rl.threshold
}

// UpdateScore updates the score of the user based on the additional score (time taken to process the request)
func (rl *RateLimiter) UpdateScore(userID common.Address, additionalScore uint32) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	newScore := rl.users[userID].score + additionalScore
	rl.users[userID] = Score{lastRequest: time.Now(), score: newScore}
}
