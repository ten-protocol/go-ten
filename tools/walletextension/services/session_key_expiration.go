package services

import (
	"time"

	gethlog "github.com/ethereum/go-ethereum/log"
	"github.com/ten-protocol/go-ten/go/common/stopcontrol"
	wecommon "github.com/ten-protocol/go-ten/tools/walletextension/common"
	"github.com/ten-protocol/go-ten/tools/walletextension/storage"
)

// SessionKeyExpirationService runs in the background and monitors users
type SessionKeyExpirationService struct {
	storage     storage.UserStorage
	logger      gethlog.Logger
	stopControl *stopcontrol.StopControl
	ticker      *time.Ticker
}

// NewSessionKeyExpirationService creates a new session key expiration service
func NewSessionKeyExpirationService(storage storage.UserStorage, logger gethlog.Logger, stopControl *stopcontrol.StopControl) *SessionKeyExpirationService {
	service := &SessionKeyExpirationService{
		storage:     storage,
		logger:      logger,
		stopControl: stopControl,
	}

	// Start the service
	go service.start()

	return service
}

// start begins the periodic user monitoring
func (s *SessionKeyExpirationService) start() {
	// TODO: Make this interval configurable
	s.ticker = time.NewTicker(10 * time.Second)

	go func() {
		defer s.ticker.Stop()

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

	s.logger.Info("Session key expiration service started")
}

// sessionKeyExpiration runs the monitoring logic
func (s *SessionKeyExpirationService) sessionKeyExpiration() {
	users, err := s.storage.GetAllUsers()
	if err != nil {
		s.logger.Error("Failed to get all users", "error", err)
		return
	}

	s.logger.Info("Session key expiration service running - checking users", "count", len(users))

	for _, user := range users {
		s.logger.Info("User found",
			"userID", wecommon.HashForLogging(user.ID),
			"accountsCount", len(user.Accounts),
			"sessionKeysCount", len(user.SessionKeys))
	}
}
