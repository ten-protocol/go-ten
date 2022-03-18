package db

import "sync"

type MemDB struct {
	lock sync.RWMutex
	db   map[string][]byte
}

func (m *MemDB) Get(id []byte) ([]byte, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.db[string(id)], nil
}

func (m *MemDB) Store(id []byte, val []byte) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	return nil
}

func (m *MemDB) Delete(id []byte) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.db, string(id))
	return nil
}

func NewMemDB() Database {
	return &MemDB{
		db: map[string][]byte{},
	}
}
