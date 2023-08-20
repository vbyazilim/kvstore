package basehttphandler

import (
	"encoding/json"
	"net/http"
	"time"
)

// Handler respresents common http handler functionality.
type Handler struct {
	CancelTimeout time.Duration
}

// JSON ...
func (h *Handler) JSON(w http.ResponseWriter, status int, d any) {
	j, err := json.Marshal(d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, _ = w.Write(j)
}
