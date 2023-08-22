package kvstorehandler

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/vbyazilim/kvstore/src/internal/kverror"
	"github.com/vbyazilim/kvstore/src/internal/service/kvstoreservice"
)

func (h *kvstoreHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		h.JSON(
			w,
			http.StatusMethodNotAllowed,
			map[string]string{"error": "method " + r.Method + " not allowed"},
		)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.JSON(
			w,
			http.StatusBadRequest,
			map[string]string{"error": err.Error()},
		)
		return
	}

	if len(body) == 0 {
		h.JSON(
			w,
			http.StatusBadRequest,
			map[string]string{"error": "empty body/payload"},
		)
		return
	}

	var handlerRequest UpdateRequest
	if err = json.Unmarshal(body, &handlerRequest); err != nil {
		h.JSON(
			w,
			http.StatusInternalServerError,
			map[string]string{"error": err.Error()},
		)
		return
	}

	if handlerRequest.Key == "" {
		h.JSON(
			w,
			http.StatusBadRequest,
			map[string]string{"error": "key is empty"},
		)
		return
	}

	if handlerRequest.Value == nil {
		h.JSON(
			w,
			http.StatusBadRequest,
			map[string]string{"error": "value is empty"},
		)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.CancelTimeout)
	defer cancel()

	serviceRequest := kvstoreservice.UpdateRequest{
		Key:   handlerRequest.Key,
		Value: handlerRequest.Value,
	}

	serviceResponse, err := h.service.Update(ctx, &serviceRequest)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			h.JSON(
				w,
				http.StatusGatewayTimeout,
				map[string]string{"error": err.Error()},
			)
			return
		}

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
				h.Logger.Error("kvstorehandler Update service.Update", "err", clientMessage)
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
