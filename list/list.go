package list

import (
	"cmp"
	"context"
	"slices"
	"strings"
	"sync"

	"github.com/peterhellberg/kop/rpc"
)

type List struct {
	set *set[string]
}

func New(vals ...string) *List {
	return &List{
		set: newSet(fn(vals, strings.ToUpper)...),
	}
}

func (l *List) Add(ctx context.Context, r rpc.AddRequest) (*rpc.AddResponse, error) {
	l.set.Add(fn(r.Items, strings.ToUpper)...)

	return &rpc.AddResponse{Items: l.set.Members()}, nil
}

func (l *List) Remove(ctx context.Context, r rpc.RemoveRequest) (*rpc.RemoveResponse, error) {
	l.set.Remove(fn(r.Items, strings.ToUpper)...)

	return &rpc.RemoveResponse{Items: l.set.Members()}, nil
}

func (l *List) Clear(ctx context.Context, r rpc.ClearRequest) (*rpc.ClearResponse, error) {
	l.set.Clear()

	return &rpc.ClearResponse{}, nil
}

func (l *List) Items(ctx context.Context, r rpc.ItemsRequest) (*rpc.ItemsResponse, error) {
	return &rpc.ItemsResponse{Items: l.set.Members()}, nil
}

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

func fn[T, V any](ts []T, fn func(T) V) []V {
	result := make([]V, len(ts))

	for i, t := range ts {
		result[i] = fn(t)
	}

	return result
}
