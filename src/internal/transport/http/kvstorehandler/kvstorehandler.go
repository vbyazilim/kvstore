package kvstorehandler

import (
	"log/slog"
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

// WithContextTimeout sets handler context cancel timeout.
func WithContextTimeout(d time.Duration) StoreHandlerOption {
	return func(s *kvstoreHandler) {
		s.Handler.CancelTimeout = d
	}
}

// WithServerEnv sets handler server env.
func WithServerEnv(env string) StoreHandlerOption {
	return func(s *kvstoreHandler) {
		s.Handler.ServerEnv = env
	}
}

// WithLogger sets handler logger.
func WithLogger(l *slog.Logger) StoreHandlerOption {
	return func(s *kvstoreHandler) {
		s.Handler.Logger = l
	}
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
}
