package services

import (
	"context"
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common/stopcontrol"
	wecommon "github.com/ten-protocol/go-ten/tools/walletextension/common"
	"github.com/ten-protocol/go-ten/tools/walletextension/storage"
)

// SessionKeyExpirationService runs in the background and monitors users
type SessionKeyExpirationService struct {
	storage         storage.UserStorage
	activityStorage storage.SessionKeyActivityStorage
	logger          gethlog.Logger
	stopControl     *stopcontrol.StopControl
	ticker          *time.Ticker
	config          *wecommon.Config
	backendRPC      *BackendRPC
	activityTracker SessionKeyActivityTracker
	txSender        TxSender
}

// NewSessionKeyExpirationService creates a new session key expiration service
func NewSessionKeyExpirationService(storage storage.UserStorage, activityStorage storage.SessionKeyActivityStorage, logger gethlog.Logger, stopControl *stopcontrol.StopControl, config *wecommon.Config, backendRPC *BackendRPC, activityTracker SessionKeyActivityTracker, txSender TxSender) *SessionKeyExpirationService {
	logger.Info("Creating session key expiration service", "expirationThreshold", config.SessionKeyExpirationThreshold.String())

	service := &SessionKeyExpirationService{
		storage:         storage,
		activityStorage: activityStorage,
		logger:          logger,
		stopControl:     stopControl,
		config:          config,
		backendRPC:      backendRPC,
		activityTracker: activityTracker,
		txSender:        txSender,
	}

	// Start the service
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("Session key expiration service creation panicked", "error", r)
			}
		}()
		service.start()
	}()
	return service
}

// start begins the periodic user monitoring
func (s *SessionKeyExpirationService) start() {
	defer func() {
		if r := recover(); r != nil {
			s.logger.Error("Session key expiration service panicked", "error", r)
		}
	}()

	// Configure interval from config
	interval := s.config.SessionKeyExpirationInterval
	s.ticker = time.NewTicker(interval)

	// load all activities from the database to make them recoverable in case of restart
	if persisted, err := s.activityStorage.Load(); err != nil {
		s.logger.Warn("Failed to load persisted session key activities", "error", err)
	} else if len(persisted) > 0 {
		loaded := make([]wecommon.SessionKeyActivity, 0, len(persisted))
		for _, a := range persisted {
			loaded = append(loaded, wecommon.SessionKeyActivity{
				Addr:       a.Addr,
				UserID:     a.UserID,
				LastActive: a.LastActive,
			})
		}
		s.activityTracker.Load(loaded)
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				s.logger.Error("Session key expiration goroutine panicked", "error", r)
			}
			s.ticker.Stop()
		}()

		for {
			select {
			case <-s.ticker.C:
				s.sessionKeyExpiration()
			case <-s.stopControl.Done():
				s.logger.Info("Session key expiration service stopped")
				return
			}
		}
	}()
	s.logger.Info("Session key expiration service started", "expirationThreshold", s.config.SessionKeyExpirationThreshold.String(), "interval", interval.String())
}

// sessionKeyExpiration runs the monitoring logic
func (s *SessionKeyExpirationService) sessionKeyExpiration() {
	s.logger.Info("Session key expiration check started")

	cutoff := time.Now().Add(-s.config.SessionKeyExpirationThreshold)

	// Get expired candidates from in-memory LRU cache
	memoryCandidates := s.activityTracker.ListOlderThan(cutoff)

	// Also query CosmosDB for expired entries (evicted from memory but still valid)
	dbCandidates, err := s.activityStorage.ListOlderThan(cutoff)
	if err != nil {
		s.logger.Warn("Failed to query DB for expired session keys", "error", err)
		dbCandidates = nil
	}

	// Merge and deduplicate: in-memory takes precedence for recent updates
	candidates := s.mergeActivityLists(memoryCandidates, dbCandidates)
	s.logger.Info("Session key expiration check", "cutoff", cutoff, "memoryCount", len(memoryCandidates), "dbCount", len(dbCandidates), "mergedCount", len(candidates))

	for _, c := range candidates {
		s.logger.Info("Processing expired session key candidate",
			"sessionKeyAddress", c.Addr.Hex(),
			"lastActive", c.LastActive,
			"timeSinceLastActive", time.Since(c.LastActive))
		// Load the user for this session key
		user, err := s.storage.GetUser(c.UserID)
		if err != nil || user == nil {
			s.logger.Error("Failed to load user for session key candidate", "error", err)
			continue
		}

		// Ensure this session key still belongs to the user
		if _, ok := user.SessionKeys[c.Addr]; !ok {
			// The session key may have been deleted; remove from both tracker and DB
			s.deleteActivity(c.Addr)
			continue
		}

		// Transfer funds to user's primary account using TxSender (sends all minus gas)
		firstAccount, err := user.GetFirstAccount()
		if err != nil {
			s.logger.Error("No primary account found for user", "error", err, "userID", wecommon.HashForLogging(user.ID))
			continue
		}

		s.logger.Info("Attempting to recover funds from expired session key",
			"sessionKeyAddress", c.Addr.Hex(),
			"toAccount", firstAccount.Address.Hex())
		txHash, err := s.txSender.SendAllMinusGasWithSK(context.Background(), user, c.Addr, *firstAccount.Address)
		if err != nil {
			s.logger.Error("Failed to recover funds from expired session key",
				"error", err,
				"userID", wecommon.HashForLogging(user.ID),
				"sessionKeyAddress", c.Addr.Hex())
			continue
		}

		s.logger.Info("Successfully initiated fund recovery transaction",
			"sessionKeyAddress", c.Addr.Hex(),
			"txHash", txHash.Hex())

		// After successful external operation, delete from both tracker and DB
		s.deleteActivity(c.Addr)
	}

	// store all activities in the database to make them persistent and recoverable in case of restart
	allActivities := s.activityTracker.ListAll()
	_ = s.activityStorage.Save(allActivities)
}

// mergeActivityLists merges activities from memory and DB, deduplicating by address.
// In-memory entries take precedence as they have the most recent state.
func (s *SessionKeyExpirationService) mergeActivityLists(memory, db []wecommon.SessionKeyActivity) []wecommon.SessionKeyActivity {
	seen := make(map[gethcommon.Address]struct{})
	result := make([]wecommon.SessionKeyActivity, 0, len(memory)+len(db))

	// Add all memory entries first (they take precedence)
	for _, item := range memory {
		if _, exists := seen[item.Addr]; !exists {
			seen[item.Addr] = struct{}{}
			result = append(result, item)
		}
	}

	// Add DB entries that aren't already in memory
	for _, item := range db {
		if _, exists := seen[item.Addr]; !exists {
			seen[item.Addr] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}

// deleteActivity removes an activity from both the in-memory tracker and the database
func (s *SessionKeyExpirationService) deleteActivity(addr gethcommon.Address) {
	s.activityTracker.Delete(addr)
	if err := s.activityStorage.Delete(addr); err != nil {
		s.logger.Warn("Failed to delete activity from DB", "addr", addr.Hex(), "error", err)
	}
}
