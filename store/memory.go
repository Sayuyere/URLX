package store

import (
	"sync"
	"urlx/logging"
)

type MemoryStore struct {
	data  map[string]string
	mutex sync.RWMutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{data: make(map[string]string)}
}

func (m *MemoryStore) Set(short, long string) {
	logger := logging.NewLogger()
	logger.Debug("Set short URL (memory store)", "short", short, "long", long)
	m.mutex.Lock()
	m.data[short] = long
	m.mutex.Unlock()
}

func (m *MemoryStore) Get(short string) (string, bool) {
	logger := logging.NewLogger()
	logger.Debug("Get short URL (memory store)", "short", short)
	m.mutex.RLock()
	long, ok := m.data[short]
	m.mutex.RUnlock()
	return long, ok
}

func (m *MemoryStore) Delete(short string) {
	logger := logging.NewLogger()
	logger.Debug("Delete short URL (memory store)", "short", short)
	m.mutex.Lock()
	delete(m.data, short)
	m.mutex.Unlock()
}
