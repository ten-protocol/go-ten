package actions

import (
	"context"
	"time"

	"github.com/obscuronet/go-obscuro/integration/common"
	"github.com/obscuronet/go-obscuro/integration/networktest"
)

// ContextValueAction is a simple action step that just sets a value on the context
func ContextValueAction(key ActionKey, value interface{}) networktest.Action {
	return &contextValueAction{key: key, value: value}
}

type contextValueAction struct {
	key   ActionKey
	value interface{}
}

func (c *contextValueAction) Run(ctx context.Context, network networktest.NetworkConnector) (context.Context, error) {
	return context.WithValue(ctx, c.key, c.value), nil
}

func (c *contextValueAction) Verify(ctx context.Context, network networktest.NetworkConnector) error {
	// nothing to verify
	return nil
}

func RandomSleepAction(minSleep time.Duration, maxSleep time.Duration) networktest.Action {
	return SleepAction(common.RndBtwTime(minSleep, maxSleep))
}

func SleepAction(duration time.Duration) networktest.Action {
	return &sleepAction{sleepDuration: duration}
}

type sleepAction struct {
	sleepDuration time.Duration
}

func (s *sleepAction) Run(ctx context.Context, network networktest.NetworkConnector) (context.Context, error) {
	time.Sleep(s.sleepDuration)
	return ctx, nil
}

func (s *sleepAction) Verify(ctx context.Context, network networktest.NetworkConnector) error {
	// nothing to verify
	return nil
}

type (
	RunFunc    func(ctx context.Context, network networktest.NetworkConnector) (context.Context, error)
	VerifyFunc func(ctx context.Context, network networktest.NetworkConnector) error
)

// NoStateNoVerifyAction allows you to create an action quickly with an in-line function when it doesn't need state or a verify method
func NoStateNoVerifyAction(run RunFunc) networktest.Action {
	return &basicAction{run: run}
}

type basicAction struct {
	run    RunFunc
	verify VerifyFunc
}

func (b *basicAction) Run(ctx context.Context, network networktest.NetworkConnector) (context.Context, error) {
	if b.run == nil {
		return ctx, nil
	}
	return b.run(ctx, network)
}

func (b *basicAction) Verify(ctx context.Context, network networktest.NetworkConnector) error {
	if b.verify == nil {
		return nil
	}
	return b.verify(ctx, network)
}
