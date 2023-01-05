package response

type ListResponse[T any] struct {
	Data []*T `json:"data"`
}
