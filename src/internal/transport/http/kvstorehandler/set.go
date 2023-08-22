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

func (h *kvstoreHandler) Set(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var handlerRequest SetRequest
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

	existingItem, err := h.service.Get(ctx, handlerRequest.Key)
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
				h.Logger.Error("kvstorehandler Set service.Get", "err", clientMessage)
			}

			if kvErr != kverror.ErrKeyNotFound {
				h.JSON(
					w,
					http.StatusBadRequest,
					map[string]string{"error": clientMessage},
				)
				return
			}
		}
	}

	// this should be nil. means, key does not exist
	if existingItem != nil {
		h.JSON(
			w,
			http.StatusConflict,
			map[string]string{"error": "can not set, '" + handlerRequest.Key + "' already exists"},
		)
		return
	}

	serviceRequest := kvstoreservice.SetRequest{
		Key:   handlerRequest.Key,
		Value: handlerRequest.Value,
	}

	serviceResponse, err := h.service.Set(ctx, &serviceRequest)
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
				h.Logger.Error("kvstorehandler Set service.Set", "err", clientMessage)
			}

			if kvErr == kverror.ErrKeyExists {
				h.JSON(w, http.StatusConflict, map[string]string{"error": clientMessage})
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
		http.StatusCreated,
		handlerResponse,
	)
}
