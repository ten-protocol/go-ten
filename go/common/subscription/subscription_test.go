package subscription

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type ThingHandler interface {
	HandleThing()
}

type TestThingHandler struct {
	Count int
}

func (t *TestThingHandler) HandleThing() {
	t.Count++
}

// try subscribing and unsubscribing a couple of handlers to a manager
func TestSubscribeAndUnsubscribe(t *testing.T) {
	manager := NewManager[ThingHandler]()
	subscriber1 := &TestThingHandler{}
	subscriber2 := &TestThingHandler{}

	unsubscribe1 := manager.Subscribe(subscriber1)
	unsubscribe2 := manager.Subscribe(subscriber2)

	// we should have 2 active subscriptions now
	assert.Equal(t, 2, len(manager.Subscribers()))
	for _, s := range manager.Subscribers() {
		s.HandleThing()
	}
	assert.Equal(t, 1, subscriber1.Count)
	assert.Equal(t, 1, subscriber2.Count)

	unsubscribe1()

	// we should just have one subscription now
	assert.Equal(t, 1, len(manager.Subscribers()))
	for _, s := range manager.Subscribers() {
		s.HandleThing()
	}
	assert.Equal(t, 1, subscriber1.Count)
	assert.Equal(t, 2, subscriber2.Count)

	unsubscribe2()

	// both subscriptions should be gone now
	assert.Equal(t, 0, len(manager.Subscribers()))
}
