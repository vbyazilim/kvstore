package kvstorehandler_test

import (
	"context"
	"fmt"
	"log/slog"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/vbyazilim/kvstore/src/internal/kverror"
	"github.com/vbyazilim/kvstore/src/internal/transport/http/kvstorehandler"
)

func TestDeleteInvalidMethod(t *testing.T) {
	handler := kvstorehandler.New()
	req := httptest.NewRequest("GET", "/key", nil)
	w := httptest.NewRecorder()

	handler.Delete(w, req)

	if w.Code != 405 && strings.Contains(w.Body.String(), "method not allowed") {
		t.Error("code not equal")
	}
}

func TestDeleteQueryParamRequired(t *testing.T) {
	handler := kvstorehandler.New()
	req := httptest.NewRequest("DELETE", "/", nil)
	w := httptest.NewRecorder()

	handler.Delete(w, req)

	if w.Code != 404 && strings.Contains(w.Body.String(), "key query param required") {
		t.Error("code not equal")
	}
}

func TestDeleteQueryParamKeyNotFound(t *testing.T) {
	handler := kvstorehandler.New()
	req := httptest.NewRequest("DELETE", "/?foo=test", nil)
	w := httptest.NewRecorder()

	handler.Delete(w, req)

	if w.Code != 404 {
		t.Error("code not equal")
	}

	if !strings.Contains(w.Body.String(), "key not present") {
		t.Error("body not equal")
	}
}

func TestDeleteTimeout(t *testing.T) {
	handler := kvstorehandler.New(
		kvstorehandler.WithContextTimeout(time.Second*-1),
		kvstorehandler.WithService(&mockService{
			deleteErr: context.DeadlineExceeded,
		}),
	)

	req := httptest.NewRequest("DELETE", "/?key=test", nil)
	w := httptest.NewRecorder()

	handler.Delete(w, req)

	if w.Code != 500 {
		t.Error("code not equal")
	}

	if !strings.Contains(w.Body.String(), "context deadline exceeded") {
		t.Error("body not equal")
	}
}

func TestDeleteErrUnknown(t *testing.T) {
	logger := slog.Default()
	handler := kvstorehandler.New(
		kvstorehandler.WithService(&mockService{
			deleteErr: kverror.ErrUnknown,
		}),
		kvstorehandler.WithLogger(logger),
	)

	req := httptest.NewRequest("DELETE", "/?key=test", nil)
	w := httptest.NewRecorder()

	handler.Delete(w, req)

	if w.Code != 500 {
		t.Error("code not equal")
	}

	fmt.Print("body: ", w.Body.String(), "\n")

	if !strings.Contains(w.Body.String(), "unknown error") {
		t.Error("body not equal")
	}
}

func TestDeleteErrKeyNotFound(t *testing.T) {
	logger := slog.Default()
	kverror.ErrKeyNotFound.AddData("key=test")

	handler := kvstorehandler.New(
		kvstorehandler.WithService(&mockService{
			deleteErr: kverror.ErrKeyNotFound,
		}),
		kvstorehandler.WithLogger(logger),
	)

	req := httptest.NewRequest("DELETE", "/?key=test", nil)
	w := httptest.NewRecorder()

	handler.Delete(w, req)

	if w.Code != 404 {
		t.Error("code not equal")
	}

	if !strings.Contains(w.Body.String(), "key not found") {
		t.Error("body not equal")
	}

	if !strings.Contains(w.Body.String(), "key=test") {
		t.Error("body not equal")
	}

	kverror.ErrKeyNotFound.DestoryData()
}

func TestDelete(t *testing.T) {
	logger := slog.Default()

	handler := kvstorehandler.New(
		kvstorehandler.WithService(&mockService{}),
		kvstorehandler.WithLogger(logger),
	)

	req := httptest.NewRequest("DELETE", "/?key=test", nil)
	w := httptest.NewRecorder()

	handler.Delete(w, req)

	if w.Code != 204 {
		t.Error("code not equal")
	}

	if w.Body.String() != "" {
		t.Error("body not equal")
	}
}
