package kvstorehandler_test

import (
	"context"
	"log/slog"
	"os"

	"github.com/vbyazilim/kvstore/src/internal/service/kvstoreservice"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

type mockService struct {
	deleteErr      error
	getErr         error
	getResponse    *kvstoreservice.ItemResponse
	listErr        error
	listResponse   *kvstoreservice.ListResponse
	setErr         error
	setResponse    *kvstoreservice.ItemResponse
	updateErr      error
	updateResponse *kvstoreservice.ItemResponse
}

func (m *mockService) Delete(_ context.Context, _ string) error {
	return m.deleteErr
}

func (m *mockService) Get(_ context.Context, _ string) (*kvstoreservice.ItemResponse, error) {
	return m.getResponse, m.getErr
}

func (m *mockService) List(_ context.Context) (*kvstoreservice.ListResponse, error) {
	return m.listResponse, m.listErr
}

func (m *mockService) Set(_ context.Context, _ *kvstoreservice.SetRequest) (*kvstoreservice.ItemResponse, error) {
	return m.setResponse, m.setErr
}

func (m *mockService) Update(_ context.Context, _ *kvstoreservice.UpdateRequest) (*kvstoreservice.ItemResponse, error) {
	return m.updateResponse, m.updateErr
}
