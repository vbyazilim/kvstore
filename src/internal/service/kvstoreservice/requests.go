package kvstoreservice

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
