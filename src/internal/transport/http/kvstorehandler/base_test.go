package kvstorehandler_test

import (
	"context"
	"errors"

	"github.com/vbyazilim/kvstore/src/internal/service/kvstoreservice"
)

var (
	errServiceDelete = errors.New("service.Delete")
	errServiceGet    = errors.New("service.Get")
	errServiceList   = errors.New("service.List")
	errServiceSet    = errors.New("service.Set")
	errServiceUpdate = errors.New("service.Update")
)

type mockService struct {
	deleteErr error
	getErr    error
	listErr   error
	setErr    error
	updateErr error
}

func (m *mockService) Delete(_ context.Context, _ string) error {
	return m.deleteErr
}

func (m *mockService) Get(_ context.Context, _ string) (*kvstoreservice.ItemResponse, error) {
	return nil, m.getErr
}

func (m *mockService) List(_ context.Context) (*kvstoreservice.ListResponse, error) {
	return nil, m.listErr
}

func (m *mockService) Set(_ context.Context, _ *kvstoreservice.SetRequest) (*kvstoreservice.ItemResponse, error) {
	return nil, m.setErr
}

func (m *mockService) Update(_ context.Context, _ *kvstoreservice.UpdateRequest) (*kvstoreservice.ItemResponse, error) {
	return nil, m.updateErr
}

type mockLogger struct {
	logErr error
}
