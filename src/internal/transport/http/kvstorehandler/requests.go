package kvstorehandler

// SetRequest is an input payload for creating new k/v item.
type SetRequest struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}
