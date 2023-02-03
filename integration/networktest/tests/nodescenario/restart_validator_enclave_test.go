package nodescenario

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/obscuronet/go-obscuro/integration/networktest/traffic"

	"github.com/obscuronet/go-obscuro/go/common/retry"
	"github.com/obscuronet/go-obscuro/integration/networktest"
	"github.com/obscuronet/go-obscuro/integration/networktest/env"
)

func TestRestartValidatorEnclave(t *testing.T) {
	networktest.TestOnlyRunsInIDE(t)
	networktest.Run(t, env.LocalDevNetwork(), restartEnclaveTest(1))
}

func restartEnclaveTest(validatorIdxToRestart int) networktest.NetworkTest {
	return &restartValidatorEnclaveTest{
		validatorIdx: validatorIdxToRestart,
	}
}

type restartValidatorEnclaveTest struct {
	validatorIdx int
}

func (r *restartValidatorEnclaveTest) Name() string {
	return "restart-enclave"
}

func (r *restartValidatorEnclaveTest) Run(network networktest.NetworkConnector) error {
	// run basic traffic test for 30secs
	test := traffic.DurationTest(traffic.NativeFundsTransfers(), 30*time.Second)
	err := test.Run(network)
	if err != nil {
		return err
	}

	validator := network.GetValidatorNode(r.validatorIdx)
	validator.StopEnclave()

	// poll the validator's health check until it stops succeeding (enclave gone down)
	err = retry.Do(func() error {
		hcErr := networktest.NodeHealthCheck(validator.HostRPCAddress())
		if hcErr == nil {
			return errors.New("health check succeeded but waiting for it to fail")
		}
		return nil
	}, retry.NewTimeoutStrategy(30*time.Second, 1*time.Second))
	if err != nil {
		return err
	}
	// give it a few seconds after the health check starts failing to tear down
	time.Sleep(3 * time.Second)

	validator.StartEnclave()

	// poll the validator's health check until it's healthy again
	err = retry.Do(func() error {
		return networktest.NodeHealthCheck(validator.HostRPCAddress())
	}, retry.NewTimeoutStrategy(30*time.Second, 1*time.Second))
	if err != nil {
		return fmt.Errorf("restarted enclave did not recover - %w", err)
	}

	postTest := traffic.DurationTest(traffic.NativeFundsTransfers(), 30*time.Second)
	err = postTest.Run(network)
	if err != nil {
		return err
	}
	return nil
}
