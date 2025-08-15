package store

import "sync"

type MemoryStore struct {
	data  map[string]string
	mutex sync.RWMutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{data: make(map[string]string)}
}

func (m *MemoryStore) Set(short, long string) {
	m.mutex.Lock()
	m.data[short] = long
	m.mutex.Unlock()
}

func (m *MemoryStore) Get(short string) (string, bool) {
	m.mutex.RLock()
	long, ok := m.data[short]
	m.mutex.RUnlock()
	return long, ok
}

func (m *MemoryStore) Delete(short string) {
	m.mutex.Lock()
	delete(m.data, short)
	m.mutex.Unlock()
}
