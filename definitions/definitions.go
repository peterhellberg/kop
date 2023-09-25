package definitions

type List interface {
	Items(ItemsRequest) ItemsResponse
	Add(AddRequest) AddResponse
	Remove(RemoveRequest) RemoveResponse
	Clear(ClearRequest) ClearResponse
}

type ItemsRequest struct{}

type ItemsResponse struct {
	Items []string
}

type AddRequest struct {
	Items []string
}

type AddResponse struct {
	Items []string
}

type RemoveRequest struct {
	Items []string
}

type RemoveResponse struct {
	Items []string
}

type ClearRequest struct{}

type ClearResponse struct{}
