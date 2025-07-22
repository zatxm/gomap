package gomap

import "sync"

type MapKV struct {
	mu   sync.RWMutex
	data *DataKv
}

func NewMap() *MapKV {
	return &MapKV{data: &DataKv{}}
}

func (m *MapKV) Set(key, val any) {
	m.mu.Lock()
	m.data.Set(key, val)
	m.mu.Unlock()
}

func (m *MapKV) Get(key any) any {
	m.mu.Lock()
	val := m.data.Get(key)
	m.mu.Unlock()
	return val
}

func (m *MapKV) Reset() {
	m.mu.Lock()
	m.data.Reset()
	m.mu.Unlock()
}

func (m *MapKV) Remove(key any) {
	m.mu.Lock()
	m.data.Remove(key)
	m.mu.Unlock()
}

func (m *MapKV) Len() int {
	return len(*m.data)
}

func (m *MapKV) All(fn func(k, v any) bool) {
	m.mu.Lock()
	data := *m.data
	for i := range data {
		k := data[i].key
		v := data[i].val
		if fn(k, v) {
			m.mu.Unlock()
			return
		}
	}
	m.mu.Unlock()
}
