package service

import "context"

// SetRequest is an input payload for Set behaviour.
type SetRequest struct {
	Key   string
	Value any
}

// UpdateRequest is an input payload for Update behaviour.
type UpdateRequest struct {
	Key   string
	Value any
}

// ItemResponse represents common k/v response element.
type ItemResponse struct {
	Key   string
	Value any
}

// ListResponse is a collection on ItemResponse.
type ListResponse []ItemResponse

// Servicer defines service behaviours.
type Servicer interface {
	Set(context.Context, *SetRequest) (*ItemResponse, error)
	Get(context.Context, string) (*ItemResponse, error)
	Update(context.Context, *UpdateRequest) (*ItemResponse, error)
	Delete(context.Context, string) error
	List(context.Context) (*ListResponse, error)
}
