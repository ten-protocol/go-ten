package rpcapi

import (
	"context"
	"fmt"

	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/ten-protocol/go-ten/tools/walletextension/services"
)

type SessionKeyAPI struct {
	we *services.Services
}

func NewSessionKeyAPI(we *services.Services) *SessionKeyAPI {
	return &SessionKeyAPI{we}
}

func (api *SessionKeyAPI) Create(ctx context.Context) (gethcommon.Address, error) {
	user, err := extractUserForRequest(ctx, api.we)
	if err != nil {
		return gethcommon.Address{}, err
	}

	sk, err := api.we.SKManager.CreateSessionKey(user)
	if err != nil {
		return gethcommon.Address{}, fmt.Errorf("unable to create session key: %w", err)
	}
	return *sk.Account.Address, nil
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

	return api.we.SKManager.ActivateSessionKey(user)
}

func (api *SessionKeyAPI) Delete(ctx context.Context) (bool, error) {
	user, err := extractUserForRequest(ctx, api.we)
	if err != nil {
		return false, err
	}

	return api.we.SKManager.DeleteSessionKey(user)
}
