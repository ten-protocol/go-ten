package ratelimiter

import (
	"log"
	"math"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/ethereum/go-ethereum/common"
)

// RequestInterval represents an interval for a request with a start and optional end timestamp.
type RequestInterval struct {
	Start time.Time
	End   *time.Time // can be nil if the request is not over yet
}

// RateLimitUser represents a user with a map of current requests.
type RateLimitUser struct {
	CurrentRequests map[uuid.UUID]RequestInterval
}

// zeroUUID is a zero UUID returned when no new request is added.
var zeroUUID uuid.UUID

// AddRequest adds a new request interval to a user's current requests and returns the UUID.
func (rl *RateLimiter) AddRequest(userID common.Address, interval RequestInterval) uuid.UUID {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	user, exists := rl.users[userID]
	if !exists {
		user = &RateLimitUser{
			CurrentRequests: make(map[uuid.UUID]RequestInterval),
		}
		rl.users[userID] = user
	}
	id := uuid.New()
	user.CurrentRequests[id] = interval
	return id
}

// UpdateRequest updates the end time of a request interval given its UUID.
func (rl *RateLimiter) UpdateRequest(userID common.Address, id uuid.UUID) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if user, userExists := rl.users[userID]; userExists {
		if request, requestExists := user.CurrentRequests[id]; requestExists {
			now := time.Now()
			request.End = &now
			user.CurrentRequests[id] = request
		} else {
			log.Printf("Request with ID %s not found for user %s.", id, userID.Hex())
		}
	} else {
		log.Printf("User %s not found while trying to update the request.", userID.Hex())
	}
}

// CountOpenRequests counts the number of requests without an End time set.
func (rl *RateLimiter) CountOpenRequests(userID common.Address) int {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	var count int
	if user, exists := rl.users[userID]; exists {
		for _, interval := range user.CurrentRequests {
			if interval.End == nil {
				count++
			}
		}
	}
	return count
}

// SumComputeTime sums the compute time for requests within the rate limiter's window
// and returns it as uint32 milliseconds.
func (rl *RateLimiter) SumComputeTime(userID common.Address) uint32 {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	var totalComputeTime time.Duration
	if user, exists := rl.users[userID]; exists {
		cutoff := time.Now().Add(-time.Duration(rl.window) * time.Millisecond)
		for _, interval := range user.CurrentRequests {
			if interval.End != nil && interval.End.After(cutoff) {
				totalComputeTime += interval.End.Sub(interval.Start)
			}
		}
	}
	return uint32(totalComputeTime / time.Millisecond)
}

type RateLimiter struct {
	mu                    sync.Mutex
	users                 map[common.Address]*RateLimitUser
	userComputeTime       uint32
	window                uint32
	maxConcurrentRequests uint32
	totalRequests         uint64
	rateLimitedRequests   uint64
}

func NewRateLimiter(rateLimitUserComputeTime uint32, rateLimitWindow uint32, concurrentRequestsLimit uint32) *RateLimiter {
	rl := &RateLimiter{
		users:                 make(map[common.Address]*RateLimitUser),
		userComputeTime:       rateLimitUserComputeTime,
		window:                rateLimitWindow,
		maxConcurrentRequests: concurrentRequestsLimit,
	}
	go rl.logRateLimitedStats()
	go rl.periodicPrune()
	return rl
}

// Allow checks if the user is allowed to make a request based on the rate limit threshold
// before comparing to the threshold also decays the score of the user based on the decay rate
func (rl *RateLimiter) Allow(userID common.Address) (bool, uuid.UUID) {
	// If the userComputeTime is 0, allow all requests (rate limiting is disabled)
	if rl.userComputeTime == 0 {
		return true, zeroUUID
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.totalRequests++

	// Check if the user has reached the maximum number of concurrent requests
	if uint32(rl.CountOpenRequests(userID)) >= rl.maxConcurrentRequests {
		rl.rateLimitedRequests++
		return false, zeroUUID
	}

	// Check if user is in limits of rate limiting
	userComputeTimeForUser := rl.SumComputeTime(userID)
	if userComputeTimeForUser <= rl.userComputeTime {
		requestUUID := rl.AddRequest(userID, RequestInterval{Start: time.Now()})
		return true, requestUUID
	}

	// If the user has exceeded the rate limit, increment the rateLimitedRequests counter
	rl.rateLimitedRequests++
	return false, zeroUUID
}

// PruneRequests deletes all requests that have ended before the rate limiter's window.
func (rl *RateLimiter) PruneRequests() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// delete all the requests that have
	cutoff := time.Now().Add(-time.Duration(rl.window) * time.Millisecond)
	for userID, user := range rl.users {
		for id, interval := range user.CurrentRequests {
			if interval.End != nil && interval.End.Before(cutoff) {
				delete(user.CurrentRequests, id)
			}
		}
		if len(user.CurrentRequests) == 0 {
			delete(rl.users, userID)
		}
	}
}

// periodically prunes the requests that have ended before the rate limiter's window every 10 * window milliseconds
func (rl *RateLimiter) periodicPrune() {
	for {
		time.Sleep(time.Duration(rl.window) * time.Millisecond * 10)
		rl.PruneRequests()
	}
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
