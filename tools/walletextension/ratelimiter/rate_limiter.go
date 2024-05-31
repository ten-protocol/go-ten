package ratelimiter

import (
	"github.com/ethereum/go-ethereum/common"
	"sync"
	"time"
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

func (rl *RateLimiter) Allow(userID common.Address, weightOfTheCall uint32) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Allow all requests if the threshold is 0
	if rl.threshold == 0 {
		return true
	}

	now := time.Now()
	userScore, exists := rl.users[userID]
	if !exists {
		// Create a new entry for the user if not exists
		rl.users[userID] = Score{lastRequest: now, score: weightOfTheCall}
	} else {
		// Calculate the decay based on the time passed
		decayTime := uint32(now.Sub(userScore.lastRequest).Seconds())
		decayedScore := int64(userScore.score) - int64(decayTime)*int64(rl.decay)
		newScore := decayedScore + int64(weightOfTheCall)

		// Ensure score does not become negative
		if newScore < 0 {
			newScore = 0
		}

		// Update user's score and last request time
		rl.users[userID] = Score{lastRequest: now, score: uint32(newScore)}
	}

	return rl.users[userID].score <= rl.threshold
}
