package actions

import (
	"context"
	"fmt"
	"time"

	"github.com/ten-protocol/go-ten/go/common/retry"
	"github.com/ten-protocol/go-ten/integration/networktest"
)

func StartValidatorEnclave(validatorIdx int) networktest.Action {
	return &startValidatorEnclaveAction{validatorIdx: validatorIdx}
}

type startValidatorEnclaveAction struct {
	validatorIdx int
}

func (s *startValidatorEnclaveAction) Run(ctx context.Context, network networktest.NetworkConnector) (context.Context, error) {
	fmt.Printf("Validator %d: starting enclave\n", s.validatorIdx)
	validator := network.GetValidatorNode(s.validatorIdx)
	// note: these actions are assuming single-enclave setups
	err := validator.StartEnclave(0)
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
	fmt.Printf("Validator %d: stopping enclave\n", s.validatorIdx)
	validator := network.GetValidatorNode(s.validatorIdx)
	err := validator.StopEnclave(0)
	if err != nil {
		return nil, err
	}
	return ctx, nil
}

func (s *stopValidatorEnclaveAction) Verify(_ context.Context, _ networktest.NetworkConnector) error {
	return nil
}

func StopValidatorHost(validatorIdx int) networktest.Action {
	return RunOnlyAction(func(ctx context.Context, network networktest.NetworkConnector) (context.Context, error) {
		fmt.Printf("Validator %d: stopping host\n", validatorIdx)
		validator := network.GetValidatorNode(validatorIdx)
		err := validator.StopHost()
		if err != nil {
			return nil, err
		}
		return ctx, nil
	})
}

func StartValidatorHost(validatorIdx int) networktest.Action {
	return RunOnlyAction(func(ctx context.Context, network networktest.NetworkConnector) (context.Context, error) {
		fmt.Printf("Validator %d: starting host\n", validatorIdx)
		validator := network.GetValidatorNode(validatorIdx)
		err := validator.StartHost()
		if err != nil {
			return nil, err
		}
		return ctx, nil
	})
}

func StopSequencerHost() networktest.Action {
	return RunOnlyAction(func(ctx context.Context, network networktest.NetworkConnector) (context.Context, error) {
		fmt.Println("Sequencer: stopping host")
		sequencer := network.GetSequencerNode()
		err := sequencer.StopHost()
		if err != nil {
			return nil, err
		}
		return ctx, nil
	})
}

func StartSequencerHost() networktest.Action {
	return RunOnlyAction(func(ctx context.Context, network networktest.NetworkConnector) (context.Context, error) {
		fmt.Println("Sequencer: starting host")
		sequencer := network.GetSequencerNode()
		err := sequencer.StartHost()
		if err != nil {
			return nil, err
		}
		return ctx, nil
	})
}

func StopSequencerEnclave(enclaveIdx int) networktest.Action {
	return RunOnlyAction(func(ctx context.Context, network networktest.NetworkConnector) (context.Context, error) {
		fmt.Printf("Sequencer: stopping enclave %d\n", enclaveIdx)
		sequencer := network.GetSequencerNode()
		err := sequencer.StopEnclave(enclaveIdx)
		if err != nil {
			return nil, err
		}
		return ctx, nil
	})
}

func StartSequencerEnclave(enclaveIdx int) networktest.Action {
	return RunOnlyAction(func(ctx context.Context, network networktest.NetworkConnector) (context.Context, error) {
		fmt.Printf("Sequencer: starting enclave %d\n", enclaveIdx)
		sequencer := network.GetSequencerNode()
		err := sequencer.StartEnclave(enclaveIdx)
		if err != nil {
			return nil, err
		}
		return ctx, nil
	})
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
		return networktest.NodeHealthCheck(validator.HostRPCWSAddress())
	}, retry.NewTimeoutStrategy(w.maxWait, 1*time.Second))
	if err != nil {
		return nil, err
	}
	return ctx, nil
}

func (w *waitForValidatorHealthCheckAction) Verify(_ context.Context, _ networktest.NetworkConnector) error {
	return nil
}

func WaitForSequencerHealthCheck(maxWait time.Duration) networktest.Action {
	return RunOnlyAction(func(ctx context.Context, network networktest.NetworkConnector) (context.Context, error) {
		sequencer := network.GetSequencerNode()
		// poll the health check until success or timeout
		err := retry.Do(func() error {
			return networktest.NodeHealthCheck(sequencer.HostRPCWSAddress())
		}, retry.NewTimeoutStrategy(maxWait, 1*time.Second))
		if err != nil {
			return nil, err
		}
		return ctx, nil
	})
}
