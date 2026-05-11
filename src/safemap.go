package src

import "sync"

type SafeMap struct {
	mu sync.RWMutex
	m  map[string]string
}

func NewSafeMap() *SafeMap {
	return &SafeMap{m: make(map[string]string)}
}

func (s *SafeMap) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.m[key]
	return v, ok
}

func (s *SafeMap) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[key] = value
}

func (s *SafeMap) GetOrDefault(key, def string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if v, ok := s.m[key]; ok {
		return v
	}
	return def
}

type SafeIntMap struct {
	mu sync.RWMutex
	m  map[int]string
}

func NewSafeIntMap() *SafeIntMap {
	return &SafeIntMap{m: make(map[int]string)}
}

func (s *SafeIntMap) Get(key int) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.m[key]
	return v, ok
}

func (s *SafeIntMap) Set(key int, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[key] = value
}
