package apiserver

import (
	"context"
	"fmt"
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
	ShutdownTimeout    = 10 * time.Second
	ServerReadTimeout  = 10 * time.Second
	ServerWriteTimeout = 10 * time.Second
	ServerIdleTimeout  = 60 * time.Second
)

// New instantiates new server instance.
func New() error {
	db := storage.MemoryDB(make(map[string]any))
	storage := kvstorage.New(
		kvstorage.WithMemoryDB(db),
	)
	service := kvstoreservice.New(
		kvstoreservice.WithStorage(storage),
	)
	kvStoreHandler := kvstorehandler.New(
		kvstorehandler.WithService(service),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/set", kvStoreHandler.Set)

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
		fmt.Println("listening", api.Addr)
		apiError <- api.ListenAndServe()
	}()

	select {
	case err := <-apiError:
		return fmt.Errorf("listen and server err: %w", err)
	case sig := <-shutdown:
		fmt.Println("shut down", sig)

		ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			if errr := api.Close(); errr != nil {
				fmt.Println("[err] api close", errr)
			}
			return fmt.Errorf("[err] could not stop server gracefully: %w", err)
		}

		wg.Done()
		wg.Wait()
	}

	return nil
}
