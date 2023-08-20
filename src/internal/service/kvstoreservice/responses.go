package kvstoreservice

// ItemResponse represents common k/v response element.
type ItemResponse struct {
	Key   string
	Value any
}

// ListResponse is a collection on ItemResponse.
type ListResponse []ItemResponse
