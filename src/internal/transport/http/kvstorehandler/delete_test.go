package kvstorehandler_test

import (
	"context"
	"net/http"
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

	if w.Code != http.StatusMethodNotAllowed && strings.Contains(w.Body.String(), "method not allowed") {
		t.Error("code not equal")
	}
}

func TestDeleteQueryParamRequired(t *testing.T) {
	handler := kvstorehandler.New()
	req := httptest.NewRequest("DELETE", "/", nil)
	w := httptest.NewRecorder()

	handler.Delete(w, req)

	if w.Code != http.StatusNotFound && strings.Contains(w.Body.String(), "key query param required") {
		t.Error("code not equal")
	}
}

func TestDeleteQueryParamKeyNotFound(t *testing.T) {
	handler := kvstorehandler.New()
	req := httptest.NewRequest("DELETE", "/?foo=test", nil)
	w := httptest.NewRecorder()

	handler.Delete(w, req)

	if w.Code != http.StatusNotFound {
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

	if w.Code != http.StatusInternalServerError {
		t.Error("code not equal")
	}

	if !strings.Contains(w.Body.String(), "context deadline exceeded") {
		t.Error("body not equal")
	}
}

func TestDeleteErrUnknown(t *testing.T) {
	handler := kvstorehandler.New(
		kvstorehandler.WithService(&mockService{
			deleteErr: kverror.ErrUnknown,
		}),
		kvstorehandler.WithLogger(logger),
	)

	req := httptest.NewRequest("DELETE", "/?key=test", nil)
	w := httptest.NewRecorder()

	handler.Delete(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Error("code not equal")
	}

	if !strings.Contains(w.Body.String(), "unknown error") {
		t.Error("body not equal")
	}
}

func TestDeleteErrKeyNotFound(t *testing.T) {
	// ignore error.
	_ = kverror.ErrKeyNotFound.AddData("key=test")

	handler := kvstorehandler.New(
		kvstorehandler.WithService(&mockService{
			deleteErr: kverror.ErrKeyNotFound,
		}),
		kvstorehandler.WithLogger(logger),
	)

	req := httptest.NewRequest("DELETE", "/?key=test", nil)
	w := httptest.NewRecorder()

	handler.Delete(w, req)

	if w.Code != http.StatusNotFound {
		t.Error("code not equal")
	}

	if !strings.Contains(w.Body.String(), "key not found") {
		t.Error("body not equal")
	}

	if !strings.Contains(w.Body.String(), "key=test") {
		t.Error("body not equal")
	}

	// ignore error.
	_ = kverror.ErrKeyNotFound.DestoryData()
}

func TestDelete(t *testing.T) {
	handler := kvstorehandler.New(
		kvstorehandler.WithService(&mockService{}),
		kvstorehandler.WithLogger(logger),
	)

	req := httptest.NewRequest("DELETE", "/?key=test", nil)
	w := httptest.NewRecorder()

	handler.Delete(w, req)

	if w.Code != http.StatusNoContent {
		t.Error("code not equal")
	}

	if w.Body.String() != "" {
		t.Error("body not equal")
	}
}
