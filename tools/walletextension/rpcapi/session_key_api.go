package rpcapi

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ten-protocol/go-ten/tools/walletextension/services"
)

type SessionKeyAPI struct {
	we *services.Services
}

func NewSessionKeyAPI(we *services.Services) *SessionKeyAPI {
	return &SessionKeyAPI{we}
}

// Create - returns hex-encoded checksum address of the newly created SK
func (api *SessionKeyAPI) Create(ctx context.Context) (string, error) {
	user, err := extractUserForRequest(ctx, api.we)
	if err != nil {
		return "", err
	}

	sk, err := api.we.SKManager.CreateSessionKey(user)
	if err != nil {
		return "", fmt.Errorf("unable to create session key: %w", err)
	}
	return (*sk.Account.Address).Hex(), nil
}

func (api *SessionKeyAPI) Delete(ctx context.Context, sessionKeyAddr string) (bool, error) {
	user, err := extractUserForRequest(ctx, api.we)
	if err != nil {
		return false, err
	}

	if !common.IsHexAddress(sessionKeyAddr) {
		return false, fmt.Errorf("invalid session key address: %s", sessionKeyAddr)
	}

	addr := common.HexToAddress(sessionKeyAddr)
	return api.we.SKManager.DeleteSessionKey(user, addr)
}

// Get returns information about a specific session key
func (api *SessionKeyAPI) Get(ctx context.Context, sessionKeyAddr string) (string, error) {
	user, err := extractUserForRequest(ctx, api.we)
	if err != nil {
		return "", err
	}

	if !common.IsHexAddress(sessionKeyAddr) {
		return "", fmt.Errorf("invalid session key address: %s", sessionKeyAddr)
	}

	addr := common.HexToAddress(sessionKeyAddr)
	sessionKey, err := api.we.SKManager.GetSessionKey(user, addr)
	if err != nil {
		return "", err
	}

	return sessionKey.Account.Address.Hex(), nil
}
