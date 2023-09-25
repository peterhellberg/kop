package list

import (
	"context"
	"strings"

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
