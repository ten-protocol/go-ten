package services

import (
	"context"
	"time"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common/stopcontrol"
	"github.com/ten-protocol/go-ten/tools/walletextension/common"
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
		loaded := make([]common.SessionKeyActivity, 0, len(persisted))
		for _, a := range persisted {
			loaded = append(loaded, common.SessionKeyActivity{
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
	candidates := s.activityTracker.ListOlderThan(cutoff)

	if len(candidates) == 0 {
		s.logger.Debug("No session keys eligible for expiration at this check", "cutoff", cutoff.UTC().Format(time.RFC3339))
		return
	}

	s.logger.Info("Found session keys eligible for expiration", "count", len(candidates), "cutoff", cutoff.UTC().Format(time.RFC3339))

	for _, c := range candidates {
		s.logger.Debug("Processing expired session key candidate", "sessionKeyAddress", c.Addr.Hex(), "userID", wecommon.HashForLogging(c.UserID), "lastActive", c.LastActive.UTC().Format(time.RFC3339))
		// Load the user for this session key
		user, err := s.storage.GetUser(c.UserID)
		if err != nil || user == nil {
			s.logger.Error("Failed to load user for session key candidate", "error", err, "userID", wecommon.HashForLogging(c.UserID), "sessionKeyAddress", c.Addr.Hex())
			continue
		}

		// Ensure this session key still belongs to the user
		if _, ok := user.SessionKeys[c.Addr]; !ok {
			// The session key may have been deleted; remove from tracker
			s.activityTracker.Delete(c.Addr)
			s.logger.Debug("Session key no longer belongs to user, removing from tracker", "sessionKeyAddress", c.Addr.Hex(), "userID", wecommon.HashForLogging(user.ID))
			continue
		}

		// Transfer funds to user's primary account using TxSender (sends all minus gas)
		// Find the first account registered with the user - we will send funds to this account
		var firstAccount *wecommon.GWAccount
		for _, account := range user.Accounts {
			firstAccount = account
			break
		}
		if firstAccount == nil || firstAccount.Address == nil {
			s.logger.Error("No primary account found for user", "userID", wecommon.HashForLogging(user.ID))
			continue
		}

		s.logger.Info("Recovering funds from expired session key", "sessionKeyAddress", c.Addr.Hex(), "to", firstAccount.Address.Hex(), "userID", wecommon.HashForLogging(user.ID))
		txHash, err := s.txSender.SendAllMinusGasWithSK(context.Background(), user, c.Addr, *firstAccount.Address)
		if err != nil {
			s.logger.Error("Failed to recover funds from expired session key",
				"error", err,
				"userID", wecommon.HashForLogging(user.ID),
				"sessionKeyAddress", c.Addr.Hex())
			continue
		}
		if (txHash == gethcommon.Hash{}) {
			// No-op (likely due to dust threshold) - we still remove from tracker to avoid infinite retries
			s.logger.Info("Skipped fund recovery (no tx sent) for expired session key", "sessionKeyAddress", c.Addr.Hex(), "userID", wecommon.HashForLogging(user.ID))
		} else {
			s.logger.Info("Submitted fund recovery tx for expired session key", "txHash", txHash.Hex(), "sessionKeyAddress", c.Addr.Hex(), "userID", wecommon.HashForLogging(user.ID))
		}
		// After successful external operation, delete from tracker
		_ = s.activityTracker.Delete(c.Addr)
	}

	// store all activities in the database to make them persistent and recoverable in case of restart
	allActivities := s.activityTracker.ListAll()
	_ = s.activityStorage.Save(allActivities)
}
