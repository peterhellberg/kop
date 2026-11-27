// Package defs is the fixture that package rpctest is generated from.
package defs

//go:generate go run ../generate ../gen.go

type List interface {
	Add(AddRequest) AddResponse
	Clear(ClearRequest) ClearResponse
	Items(ItemsRequest) ItemsResponse
}

type AddRequest struct {
	Items []string
}

type AddResponse struct {
	Items []string
}

type ClearRequest struct{}

type ClearResponse struct{}

type ItemsRequest struct {
	Thing  *Thing
	ByName map[string][]Thing `json:"by_name,omitempty"`
}

type ItemsResponse struct {
	Items []string
}

type Thing struct {
	Name       string
	Children   []Thing
	ID         string
	HTTPServer string
	UserID     string
}
