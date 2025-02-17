package actions

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ten-protocol/go-ten/integration/networktest"
	"golang.org/x/sync/errgroup"
)

type MultiAction struct {
	actions    []networktest.Action
	isParallel bool
	name       string
	start      time.Time
	end        time.Time
}

func (m *MultiAction) String() string {
	return fmt.Sprintf("%s (%d actions in %s)", m.name, len(m.actions), m.executionType())
}

func (m *MultiAction) Run(ctx context.Context, network networktest.NetworkConnector) (context.Context, error) {
	m.recordStart()
	var err error
	if m.isParallel {
		ctx, err = m.runParallel(ctx, network)
	} else {
		ctx, err = m.runSeries(ctx, network)
	}
	if err != nil {
		return nil, err
	}
	m.recordEnd()
	return ctx, nil
}

func (m *MultiAction) runSeries(ctx context.Context, network networktest.NetworkConnector) (context.Context, error) {
	var err error
	for _, a := range m.actions {
		start := time.Now()
		ctx, err = a.Run(ctx, network)
		if err != nil {
			fmt.Printf("error %s (%s)\n", err, time.Since(start).Round(time.Millisecond))
			return ctx, err
		}
	}
	return ctx, nil
}

func (m *MultiAction) runParallel(ctx context.Context, network networktest.NetworkConnector) (context.Context, error) {
	grp, _ := errgroup.WithContext(ctx)
	var err error
	for _, a := range m.actions {
		action := a
		grp.Go(func() error {
			// note: we cannot easily merge contexts so ctx is not modified by a parallel execution
			if _, err = action.Run(ctx, network); err != nil {
				return err
			}
			return nil
		})
	}
	if err = grp.Wait(); err != nil {
		return nil, err
	}

	return ctx, nil
}

func (m *MultiAction) Verify(ctx context.Context, network networktest.NetworkConnector) error {
	var actionFailures []error
	mu := sync.Mutex{} // mutex for modifying the action failures

	grp, _ := errgroup.WithContext(ctx)
	var err error
	for _, a := range m.actions {
		action := a
		grp.Go(func() error {
			if err = action.Verify(ctx, network); err != nil {
				actionErr := fmt.Errorf("%s failed - %w", action, err)
				fmt.Println(actionErr)

				mu.Lock()
				actionFailures = append(actionFailures, actionErr)
				mu.Unlock()

				return err
			}
			return nil
		})
	}
	if err = grp.Wait(); err != nil {
		return fmt.Errorf("series failed, %d / %d failed - %s", len(actionFailures), len(m.actions), actionFailures)
	}

	return nil
}

func (m *MultiAction) recordStart() {
	if m.name != "" {
		fmt.Printf("START :: %s\n", m)
	}
	m.start = time.Now()
}

func (m *MultiAction) recordEnd() {
	m.end = time.Now()
	if m.name != "" {
		fmt.Printf("END :: %s  [ %s ]\n", m, m.end.Sub(m.start))
	}
}

func (m *MultiAction) executionType() string {
	if m.isParallel {
		return "parallel"
	}
	return "series"
}

func NamedSeries(name string, actions ...networktest.Action) *MultiAction {
	return &MultiAction{
		actions: actions,
		name:    name,
	}
}

func Series(actions ...networktest.Action) *MultiAction {
	return NamedSeries("", actions...)
}

func NamedParallel(name string, actions ...networktest.Action) *MultiAction {
	return &MultiAction{actions: actions, name: name, isParallel: true}
}

func Parallel(actions ...networktest.Action) *MultiAction {
	return NamedParallel("", actions...)
}
