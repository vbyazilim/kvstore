package apiserver

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vbyazilim/kvstore/src/internal/service/kvstoreservice"
	"github.com/vbyazilim/kvstore/src/internal/storage/memory/kvstorage"
	"github.com/vbyazilim/kvstore/src/internal/transport/http/kvstorehandler"
	"github.com/vbyazilim/kvstore/src/releaseinfo"
)

// constants.
const (
	ContextCancelTimeout = 5 * time.Second
	ShutdownTimeout      = 10 * time.Second
	ServerReadTimeout    = 10 * time.Second
	ServerWriteTimeout   = 10 * time.Second
	ServerIdleTimeout    = 60 * time.Second

	apiV1Prefix = "/api/v1"
)

type apiServer struct {
	db        kvstorage.MemoryDB
	logLevel  slog.Level
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

// WithLogLevel sets logLevel option.
func WithLogLevel(level string) Option {
	return func(s *apiServer) {
		var logLevel slog.Level

		switch level {
		case "DEBUG":
			logLevel = slog.LevelDebug
		case "WARN":
			logLevel = slog.LevelWarn
		case "ERROR":
			logLevel = slog.LevelError
		default:
			logLevel = slog.LevelInfo
		}

		s.logLevel = logLevel
	}
}

// New instantiates new server instance.
func New(options ...Option) error {
	apisrvr := &apiServer{
		db:       kvstorage.MemoryDB(make(map[string]any)), // default db
		logLevel: slog.LevelInfo,
	}

	for _, o := range options {
		o(apisrvr)
	}

	// default logging options if logger not present.
	if apisrvr.logger == nil {
		logHandlerOpts := &slog.HandlerOptions{Level: apisrvr.logLevel}
		logHandler := slog.NewJSONHandler(os.Stdout, logHandlerOpts)
		apisrvr.logger = slog.New(logHandler)
	}
	slog.SetDefault(apisrvr.logger)

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

	mux.HandleFunc("/healthz/live/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		j, _ := json.Marshal(map[string]any{
			"server":            apisrvr.serverEnv,
			"version":           releaseinfo.Version,
			"build_information": releaseinfo.BuildInformation,
			"message":           "liveness is OK!, server is ready to accept connections",
		})
		_, _ = w.Write(j)
	})
	mux.HandleFunc("/healthz/ready/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		j, _ := json.Marshal(map[string]any{
			"server":            apisrvr.serverEnv,
			"version":           releaseinfo.Version,
			"build_information": releaseinfo.BuildInformation,
			"message":           "readiness is OK!, server is ready to accept connections",
		})
		_, _ = w.Write(j)
	})

	mux.HandleFunc(apiV1Prefix+"/set/", kvStoreHandler.Set)
	mux.HandleFunc(apiV1Prefix+"/get/", kvStoreHandler.Get)
	mux.HandleFunc(apiV1Prefix+"/update/", kvStoreHandler.Update)
	mux.HandleFunc(apiV1Prefix+"/delete/", kvStoreHandler.Delete)
	mux.HandleFunc(apiV1Prefix+"/list/", kvStoreHandler.List)

	api := &http.Server{
		Addr:         ":8000",
		Handler:      appendSlashMiddleware(httpLoggingMiddleware(logger, mux)),
		ReadTimeout:  ServerReadTimeout,
		WriteTimeout: ServerWriteTimeout,
		IdleTimeout:  ServerIdleTimeout,
	}

	shutdown := make(chan os.Signal, 1)
	apiError := make(chan error, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// var wg sync.WaitGroup
	// wg.Add(1)

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

		// wg.Done()
		// wg.Wait()
	}

	return nil
}
