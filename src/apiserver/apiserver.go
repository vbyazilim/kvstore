package apiserver

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/vbyazilim/kvstore/src/internal/service/kvstoreservice"
	"github.com/vbyazilim/kvstore/src/internal/storage"
	"github.com/vbyazilim/kvstore/src/internal/storage/memory/kvstorage"
	"github.com/vbyazilim/kvstore/src/internal/transport/http/kvstorehandler"
)

// constants.
const (
	ContextCancelTimeout = 5 * time.Second
	ShutdownTimeout      = 10 * time.Second
	ServerReadTimeout    = 10 * time.Second
	ServerWriteTimeout   = 10 * time.Second
	ServerIdleTimeout    = 60 * time.Second
)

type apiServer struct {
	db        storage.MemoryDB
	logger    *slog.Logger
	serverEnv string
}

// Option represents api server option type.
type Option func(*apiServer)

// WithLogger sets logger option.
func WithLogger(l *slog.Logger) Option {
	return func(s *apiServer) {
		s.logger = l
	}
}

// WithServerEnv sets serverEnv option.
func WithServerEnv(env string) Option {
	return func(s *apiServer) {
		s.serverEnv = env
	}
}

// New instantiates new server instance.
func New(options ...Option) error {
	apisrvr := &apiServer{
		db:     storage.MemoryDB(make(map[string]any)),        // default db
		logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)), // default logger
	}

	for _, o := range options {
		o(apisrvr)
	}

	if apisrvr.serverEnv == "" {
		apisrvr.serverEnv = "production" // default server environment
	}

	logger := apisrvr.logger

	storage := kvstorage.New(
		kvstorage.WithMemoryDB(apisrvr.db),
	)
	service := kvstoreservice.New(
		kvstoreservice.WithStorage(storage),
	)
	kvStoreHandler := kvstorehandler.New(
		kvstorehandler.WithService(service),
		kvstorehandler.WithContextTimeout(ContextCancelTimeout),
		kvstorehandler.WithServerEnv(apisrvr.serverEnv),
		kvstorehandler.WithLogger(logger),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/set", kvStoreHandler.Set)
	mux.HandleFunc("/get", kvStoreHandler.Get)
	mux.HandleFunc("/list", kvStoreHandler.List)

	api := &http.Server{
		Addr:         ":8000",
		Handler:      mux,
		ReadTimeout:  ServerReadTimeout,
		WriteTimeout: ServerWriteTimeout,
		IdleTimeout:  ServerIdleTimeout,
	}

	shutdown := make(chan os.Signal, 1)
	apiError := make(chan error, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		logger.Info("starting api server", "listening", api.Addr, "env", apisrvr.serverEnv)
		apiError <- api.ListenAndServe()
	}()

	select {
	case err := <-apiError:
		return fmt.Errorf("listen and server err: %w", err)
	case sig := <-shutdown:
		logger.Info("starting shutdown", "pid", sig)
		defer logger.Info("shutdown completed", "pid", sig)

		ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			if errr := api.Close(); errr != nil {
				logger.Error("api close", "err", errr)
			}
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}

		wg.Done()
		wg.Wait()
	}

	return nil
}
