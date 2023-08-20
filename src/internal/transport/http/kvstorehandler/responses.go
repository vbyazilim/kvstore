package kvstorehandler

// ItemResponse represents k/v item.
type ItemResponse struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}
