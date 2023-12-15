package store

import (
	"distapp/models"
	"sync"
)

type InMemory struct {
	m       sync.RWMutex
	storage map[string]models.Value
}

func NewInMemory() *InMemory {
	return &InMemory{
		storage: map[string]models.Value{},
	}
}

func (m *InMemory) SetValue(newVal string, key string) {
	m.m.Lock()
	defer m.m.Unlock()
	val, ok := m.storage[key]
	gen := 1
	if ok {
		gen = val.Gen + 1
	}
	m.storage[key] = models.Value{Val: newVal, Gen: gen}
}
func (m *InMemory) GetValue(key string) (string, int) {
	m.m.RLock()
	defer m.m.RUnlock()
	val, ok := m.storage[key]
	if ok {
		return val.Val, val.Gen
	}
	return "", 0
}

func (m *InMemory) NotifyValue(curVal string, key string, curGeneration int) bool {
	m.m.Lock()
	defer m.m.Unlock()
	val, ok := m.storage[key]
	if !ok {
		m.storage[key] = models.Value{Gen: 1, Val: curVal}
		return true
	}
	if ok && curGeneration > val.Gen {
		m.storage[key] = models.Value{Gen: curGeneration, Val: curVal}
		return true
	}
	return false
}
