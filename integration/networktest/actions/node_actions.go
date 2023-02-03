package actions

import (
	"context"
	"time"

	"github.com/obscuronet/go-obscuro/go/common/retry"
	"github.com/obscuronet/go-obscuro/integration/networktest"
)

func StartValidatorEnclave(validatorIdx int) networktest.Action {
	return &startValidatorEnclaveAction{validatorIdx: validatorIdx}
}

type startValidatorEnclaveAction struct {
	validatorIdx int
}

func (s *startValidatorEnclaveAction) Run(ctx context.Context, network networktest.NetworkConnector) (context.Context, error) {
	validator := network.GetValidatorNode(s.validatorIdx)
	err := validator.StartEnclave()
	if err != nil {
		return nil, err
	}
	return ctx, nil
}

func (s *startValidatorEnclaveAction) Verify(_ context.Context, _ networktest.NetworkConnector) error {
	return nil
}

func StopValidatorEnclave(validatorIdx int) networktest.Action {
	return &stopValidatorEnclaveAction{validatorIdx: validatorIdx}
}

type stopValidatorEnclaveAction struct {
	validatorIdx int
}

func (s *stopValidatorEnclaveAction) Run(ctx context.Context, network networktest.NetworkConnector) (context.Context, error) {
	validator := network.GetValidatorNode(s.validatorIdx)
	err := validator.StopEnclave()
	if err != nil {
		return nil, err
	}
	return ctx, nil
}

func (s *stopValidatorEnclaveAction) Verify(_ context.Context, _ networktest.NetworkConnector) error {
	return nil
}

func WaitForValidatorHealthCheck(validatorIdx int, maxWait time.Duration) networktest.Action {
	return &waitForValidatorHealthCheckAction{
		validatorIdx: validatorIdx,
		maxWait:      maxWait,
	}
}

type waitForValidatorHealthCheckAction struct {
	validatorIdx int
	maxWait      time.Duration
}

func (w *waitForValidatorHealthCheckAction) Run(ctx context.Context, network networktest.NetworkConnector) (context.Context, error) {
	validator := network.GetValidatorNode(w.validatorIdx)
	// poll the health check until success or timeout
	err := retry.Do(func() error {
		return networktest.NodeHealthCheck(validator.HostRPCAddress())
	}, retry.NewTimeoutStrategy(30*time.Second, 1*time.Second))
	if err != nil {
		return nil, err
	}
	return ctx, nil
}

func (w *waitForValidatorHealthCheckAction) Verify(_ context.Context, _ networktest.NetworkConnector) error {
	return nil
}
