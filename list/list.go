package list

import (
	"context"
	"strings"

	"github.com/peterhellberg/kop/rpc"
)

type Store interface {
	Add(...string)
	Remove(...string)
	Contains(string) bool
	Members() []string
	Clear()
}

type List struct {
	store Store
}

func New(store Store) *List {
	return &List{
		store: store,
	}
}

func (l *List) Add(ctx context.Context, r rpc.AddRequest) (*rpc.AddResponse, error) {
	l.store.Add(fn(r.Items, strings.ToUpper)...)

	return &rpc.AddResponse{Items: l.store.Members()}, nil
}

func (l *List) Remove(ctx context.Context, r rpc.RemoveRequest) (*rpc.RemoveResponse, error) {
	l.store.Remove(fn(r.Items, strings.ToUpper)...)

	return &rpc.RemoveResponse{Items: l.store.Members()}, nil
}

func (l *List) Clear(ctx context.Context, r rpc.ClearRequest) (*rpc.ClearResponse, error) {
	l.store.Clear()

	return &rpc.ClearResponse{}, nil
}

func (l *List) Items(ctx context.Context, r rpc.ItemsRequest) (*rpc.ItemsResponse, error) {
	return &rpc.ItemsResponse{Items: l.store.Members()}, nil
}

func fn[T, V any](ts []T, fn func(T) V) []V {
	result := make([]V, len(ts))

	for i, t := range ts {
		result[i] = fn(t)
	}

	return result
}
