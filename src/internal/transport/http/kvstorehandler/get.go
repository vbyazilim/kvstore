package kvstorehandler

import (
	"context"
	"errors"
	"net/http"

	"github.com/vbyazilim/kvstore/src/internal/kverror"
)

func (h *kvstoreHandler) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.JSON(
			w,
			http.StatusMethodNotAllowed,
			map[string]string{"error": "method " + r.Method + " not allowed"},
		)
		return
	}

	if len(r.URL.Query()) == 0 {
		h.JSON(
			w,
			http.StatusNotFound,
			map[string]string{"error": "key query param required"},
		)
		return
	}

	keys, ok := r.URL.Query()["key"]
	if !ok {
		h.JSON(
			w,
			http.StatusNotFound,
			map[string]string{"error": "key not present"},
		)
		return
	}

	key := keys[0]

	ctx, cancel := context.WithTimeout(r.Context(), h.CancelTimeout)
	defer cancel()

	serviceResponse, err := h.service.Get(ctx, key)
	if err != nil {
		var kvErr *kverror.Error

		if errors.As(err, &kvErr) {
			clientMessage := kvErr.Message
			if kvErr.Data != nil {
				data, ok := kvErr.Data.(string)
				if ok {
					clientMessage = clientMessage + ", " + data
				}
			}

			if kvErr.Loggable {
				h.Logger.Error("kvstorehandler Set", "err", clientMessage)
			}

			if kvErr == kverror.ErrKeyNotFound {
				h.JSON(w, http.StatusNotFound, map[string]string{"error": clientMessage})
				return
			}
		}
		h.JSON(
			w,
			http.StatusInternalServerError,
			map[string]string{"error": err.Error()},
		)
		return
	}

	handlerResponse := ItemResponse{
		Key:   serviceResponse.Key,
		Value: serviceResponse.Value,
	}

	h.JSON(
		w,
		http.StatusOK,
		handlerResponse,
	)
}
