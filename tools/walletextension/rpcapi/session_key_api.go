package rpcapi

import (
	"context"
	"fmt"

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

func (api *SessionKeyAPI) Activate(ctx context.Context) (bool, error) {
	user, err := extractUserForRequest(ctx, api.we)
	if err != nil {
		return false, err
	}

	return api.we.SKManager.ActivateSessionKey(user)
}

func (api *SessionKeyAPI) Deactivate(ctx context.Context) (bool, error) {
	user, err := extractUserForRequest(ctx, api.we)
	if err != nil {
		return false, err
	}

	return api.we.SKManager.DeactivateSessionKey(user)
}

func (api *SessionKeyAPI) Delete(ctx context.Context) (bool, error) {
	user, err := extractUserForRequest(ctx, api.we)
	if err != nil {
		return false, err
	}

	return api.we.SKManager.DeleteSessionKey(user)
}
