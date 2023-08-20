package kvstorehandler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/vbyazilim/kvstore/src/internal/service"
	"github.com/vbyazilim/kvstore/src/internal/transport"
	"github.com/vbyazilim/kvstore/src/internal/transport/http/basehttphandler"
)

var _ transport.KVStoreHTTPHandler = (*kvstoreHandler)(nil) // compile time proof

type kvstoreHandler struct {
	basehttphandler.Handler

	service service.Servicer
}

// StoreHandlerOption represents store handler option type.
type StoreHandlerOption func(*kvstoreHandler)

// WithService sets service option.
func WithService(srvc service.Servicer) StoreHandlerOption {
	return func(s *kvstoreHandler) {
		s.service = srvc
	}
}

func (h *kvstoreHandler) Set(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errMessage := fmt.Sprintf("method %s not allowed", r.Method)
		h.JSON(
			w,
			http.StatusMethodNotAllowed,
			map[string]string{"error": errMessage},
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

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	serviceRequest := service.SetRequest{
		Key:   handlerRequest.Key,
		Value: handlerRequest.Value,
	}

	serviceResponse, err := h.service.Set(ctx, &serviceRequest)
	if err != nil {
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

// New instantiates new kvstoreHandler instance.
func New(options ...StoreHandlerOption) transport.KVStoreHTTPHandler {
	kvsh := &kvstoreHandler{
		Handler: basehttphandler.Handler{},
	}

	for _, o := range options {
		o(kvsh)
	}

	return kvsh
	//
	// return &kvstoreHandler{
	// 	Handler: basehttphandler.Handler{},
	// 	service: srv,
	// }
}
