package subscription

import (
	"sync"

	"github.com/google/uuid"
)

// Manager is a subscription manager - allows for thread-safe registering/unregistering subscribers of type T
type Manager[T any] struct {
	subscribers map[uuid.UUID]T
	mx          sync.RWMutex
}

// NewManager creates a new subscription manager
func NewManager[T any]() *Manager[T] {
	return &Manager[T]{
		subscribers: make(map[uuid.UUID]T),
		mx:          sync.RWMutex{},
	}
}

// Subscribe adds a new subscriber, returns an unsubscribe function
func (m *Manager[T]) Subscribe(subscriber T) func() {
	m.mx.Lock()
	defer m.mx.Unlock()

	id := uuid.New()
	m.subscribers[id] = subscriber

	return func() {
		m.mx.Lock()
		defer m.mx.Unlock()

		delete(m.subscribers, id)
	}
}

// Subscribers returns a list of all subscribers
func (m *Manager[T]) Subscribers() []T {
	m.mx.RLock()
	defer m.mx.RUnlock()

	// make a new slice of the subscribers
	subscribers := make([]T, 0, len(m.subscribers))
	for _, subscriber := range m.subscribers {
		subscribers = append(subscribers, subscriber)
	}

	return subscribers
}
