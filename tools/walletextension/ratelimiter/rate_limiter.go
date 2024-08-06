package ratelimiter

import (
	"math"
	"sync"
	"time"

	gethlog "github.com/ethereum/go-ethereum/log"

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
	// If the userComputeTime is 0, do nothing (rate limiting is disabled)
	if rl.GetUserComputeTime() == 0 {
		return zeroUUID
	}
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

// SetRequestEnd updates the end time of a request interval given its UUID.
func (rl *RateLimiter) SetRequestEnd(userID common.Address, id uuid.UUID) {
	// If the userComputeTime is 0, do nothing (rate limiting is disabled)
	if rl.GetUserComputeTime() == 0 {
		return
	}

	if user, userExists := rl.users[userID]; userExists {
		if request, requestExists := user.CurrentRequests[id]; requestExists {
			rl.mu.Lock()
			defer rl.mu.Unlock()
			now := time.Now()
			request.End = &now
			user.CurrentRequests[id] = request
		} else {
			rl.logger.Info("Request with ID %s not found for user %s.", id, userID.Hex())
		}
	} else {
		rl.logger.Info("User %s not found while trying to update the request.", userID.Hex())
	}
}

// CountOpenRequests counts the number of requests without an End time set.
func (rl *RateLimiter) CountOpenRequests(userID common.Address) int {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

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
func (rl *RateLimiter) SumComputeTime(userID common.Address) time.Duration {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	var totalComputeTime time.Duration
	if user, exists := rl.users[userID]; exists {
		cutoff := time.Now().Add(-rl.window)
		for _, interval := range user.CurrentRequests {
			// if the request has ended and it's within the window, add the compute time
			if interval.End != nil && interval.End.After(cutoff) {
				totalComputeTime += interval.End.Sub(interval.Start)
			}
			// if the request hasn't ended yet, add the compute time until now
			if interval.End == nil {
				totalComputeTime += time.Since(interval.Start)
			}
		}
	}
	return totalComputeTime
}

type RateLimiter struct {
	mu                    sync.RWMutex
	users                 map[common.Address]*RateLimitUser
	userComputeTime       time.Duration
	window                time.Duration
	maxConcurrentRequests uint32
	totalRequests         uint64
	rateLimitedRequests   uint64
	logger                gethlog.Logger
}

// IncrementTotalRequests increments the total requests counter by 1 with thread safety.
func (rl *RateLimiter) IncrementTotalRequests() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.totalRequests++
}

// IncrementRateLimitedRequests increments the total requests counter by 1 with thread safety.
func (rl *RateLimiter) IncrementRateLimitedRequests() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.rateLimitedRequests++
}

// GetMaxConcurrentRequest returns the maximum number of concurrent requests allowed.
func (rl *RateLimiter) GetMaxConcurrentRequest() uint32 {
	rl.mu.RLock()
	defer rl.mu.RUnlock()
	return rl.maxConcurrentRequests
}

// GetUserComputeTime returns the user compute time
func (rl *RateLimiter) GetUserComputeTime() time.Duration {
	rl.mu.RLock()
	defer rl.mu.RUnlock()
	return rl.userComputeTime
}

func NewRateLimiter(rateLimitUserComputeTime time.Duration, rateLimitWindow time.Duration, concurrentRequestsLimit uint32, logger gethlog.Logger) *RateLimiter {
	rl := &RateLimiter{
		users:                 make(map[common.Address]*RateLimitUser),
		userComputeTime:       rateLimitUserComputeTime,
		window:                rateLimitWindow,
		maxConcurrentRequests: concurrentRequestsLimit,
		logger:                logger,
	}

	// If the userComputeTime is 0 (rate limiting is disabled) we don't need to prune and log rate limited stats
	if rl.GetUserComputeTime() != 0 {
		go rl.logRateLimitedStats()
		go rl.periodicPrune()
	}

	return rl
}

// Allow checks if the user is allowed to make a request based on the rate limit threshold
// before comparing to the threshold also decays the score of the user based on the decay rate
func (rl *RateLimiter) Allow(userID common.Address) (bool, uuid.UUID) {
	// If the userComputeTime is 0, allow all requests (rate limiting is disabled)
	if rl.GetUserComputeTime() == 0 {
		return true, zeroUUID
	}
	// Increment the total requests counter for statistics
	rl.IncrementTotalRequests()

	// Check if the user has reached the maximum number of concurrent requests
	if uint32(rl.CountOpenRequests(userID)) >= rl.GetMaxConcurrentRequest() {
		rl.IncrementRateLimitedRequests()
		rl.logger.Info("User %s has reached the maximum number of concurrent requests.", userID.Hex())
		return false, zeroUUID
	}

	// Check if user is in limits of rate limiting
	userComputeTimeForUser := rl.SumComputeTime(userID)
	if userComputeTimeForUser > rl.userComputeTime {
		rl.IncrementRateLimitedRequests()
		rl.logger.Info("User %s has reached the rate limit threshold.", userID.Hex())
		return false, zeroUUID
	}

	requestUUID := rl.AddRequest(userID, RequestInterval{Start: time.Now()})
	return true, requestUUID
}

// PruneRequests deletes all requests that have ended before the rate limiter's window.
func (rl *RateLimiter) PruneRequests() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	startTime := time.Now()
	// delete all the requests that have
	cutoff := time.Now().Add(-rl.window)
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
	timeTaken := time.Since(startTime)
	if timeTaken > 1*time.Second {
		rl.logger.Warn("PruneRequests completed in %s", timeTaken)
	}
}

// periodically prunes the requests that have ended before the rate limiter's window every 10 * window milliseconds
func (rl *RateLimiter) periodicPrune() {
	for {
		time.Sleep(rl.window / 2)
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
		rl.logger.Info("Total requests: %d, Rate-limited requests: %d (%.4f%%)", totalRequests, rateLimitedRequests, rateLimitedPercentage)
	}
}
