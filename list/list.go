package list

import (
	"context"
	"strings"

	"github.com/peterhellberg/kop/rpc"
)

type List struct {
	store Store[string]
}

func New(vals ...string) *List {
	return &List{
		store: newSet(fn(vals, strings.ToUpper)...),
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
