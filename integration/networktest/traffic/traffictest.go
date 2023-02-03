package traffic

import (
	"fmt"
	"time"

	"github.com/obscuronet/go-obscuro/integration/networktest"
)

// TestOpt functions allow overriding of durationTrafficTest defaults. E.g. WithNumUsers(20), WithActionsPerSec(50)
type TestOpt func(test *durationTrafficTest)

// DurationTest returns a NetworkTest that runs a TrafficRunner for a specified duration and then verifies it
func DurationTest(runner Runner, duration time.Duration, opts ...TestOpt) networktest.NetworkTest {
	t := &durationTrafficTest{
		runner:        runner,
		duration:      duration,
		numSimUsers:   5,
		actionsPerSec: 2,
		verifier:      SuccessesVerifier,
	}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

type RunVerifier func(RunData, networktest.NetworkConnector) error

// DurationTest is an implementation of NetworkTest that uses a traffic runner for a duration
type durationTrafficTest struct {
	runner        Runner
	duration      time.Duration
	verifier      RunVerifier
	numSimUsers   int
	actionsPerSec int
}

func (t *durationTrafficTest) Run(network networktest.NetworkConnector) error {
	err := t.runner.Start(network)
	if err != nil {
		return err
	}
	time.Sleep(t.duration)
	t.runner.Stop()
	return t.verifier(t.runner.RunData(), network)
}

func (t *durationTrafficTest) Name() string {
	return fmt.Sprintf("traffic-%s", t.runner.Name())
}

// WithVerifier is a construction Option for the DurationTest that allows you to override the ver
func WithVerifiers(verifiers []RunVerifier) TestOpt {
	return func(test *durationTrafficTest) {
	}
}
