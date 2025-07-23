package gomap

import (
	"slices"
	"sync"
)

type Set []any

func (s *Set) Set(v any) {
	if v == nil {
		return
	}

	i := slices.Index(*s, v)
	if i > -1 {
		(*s)[i] = v
		return
	}

	c := cap(*s)
	n := len(*s)
	if c > n {
		*s = (*s)[:n+1]
		(*s)[n] = v
		return
	}

	*s = append(*s, v)
}

func (s *Set) Contains(v any) bool {
	return slices.Contains(*s, v)
}

func (s *Set) Len() int {
	return len(*s)
}

func (s *Set) IsEmpty() bool {
	return s.Len() == 0
}

func (s *Set) Remove(v any) {
	*s = slices.DeleteFunc(*s, func(val any) bool {
		if val == v {
			return true
		}
		return false
	})
}

func (s *Set) Reset() {
	n := len(*s)
	for i := range n {
		(*s)[i] = nil
	}
	*s = (*s)[:0]
}

func (s *Set) Peek(fn func(v any) bool) {
	data := *s
	if slices.ContainsFunc(data, fn) {
		return
	}
}

type Sets struct {
	mu   sync.RWMutex
	data *Set
}

func NewSets() *Sets {
	return &Sets{data: new(Set)}
}

func (s *Sets) Set(v any) {
	s.mu.Lock()
	s.data.Set(v)
	s.mu.Unlock()
}

func (s *Sets) Contains(v any) bool {
	return slices.Contains(*s.data, v)
}

func (s *Sets) Len() int {
	return len(*s.data)
}

func (s *Sets) IsEmpty() bool {
	return s.Len() == 0
}

func (s *Sets) Remove(v any) {
	s.mu.Lock()
	s.data.Remove(v)
	s.mu.Unlock()
}

func (s *Sets) Reset() {
	s.mu.Lock()
	s.data.Reset()
	s.mu.Unlock()
}

func (s *Sets) Peek(fn func(v any) bool) {
	s.mu.Lock()
	data := *s.data
	if slices.ContainsFunc(data, fn) {
		s.mu.Unlock()
		return
	}
	s.mu.Unlock()
}
