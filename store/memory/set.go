package memory

import (
	"cmp"
	"slices"
	"sync"
)

type set[E cmp.Ordered] struct {
	mu sync.RWMutex
	m  map[E]struct{}
}

func newSet[E cmp.Ordered](vals ...E) *set[E] {
	s := &set[E]{
		m: map[E]struct{}{},
	}

	for _, v := range vals {
		s.m[v] = struct{}{}
	}

	return s
}

func (s *set[E]) Add(vals ...E) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, v := range vals {
		s.m[v] = struct{}{}
	}
}

func (s *set[E]) Remove(vals ...E) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, v := range vals {
		delete(s.m, v)
	}
}

func (s *set[E]) Contains(v E) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.m[v]
	return ok
}

func (s *set[E]) Members() []E {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]E, 0, len(s.m))

	for v := range s.m {
		result = append(result, v)
	}

	slices.Sort(result)

	return result
}

func (s *set[E]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	clear(s.m)
}
