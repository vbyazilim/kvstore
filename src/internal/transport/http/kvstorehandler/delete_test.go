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
	req := httptest.NewRequest(http.MethodGet, "/key", nil)
	w := httptest.NewRecorder()

	handler.Delete(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("wrong status code, want: %d, got: %d", http.StatusMethodNotAllowed, w.Code)
	}

	shouldContain := "method GET not allowed"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want: %s, got: %s", shouldContain, w.Body.String())
	}
}

func TestDeleteQueryParamRequired(t *testing.T) {
	handler := kvstorehandler.New()
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	w := httptest.NewRecorder()

	handler.Delete(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("wrong status code, want: %d, got: %d", http.StatusNotFound, w.Code)
	}

	shouldContain := "key query param required"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want: %s, got: %s", shouldContain, w.Body.String())
	}
}

func TestDeleteQueryParamKeyNotFound(t *testing.T) {
	handler := kvstorehandler.New()
	req := httptest.NewRequest(http.MethodDelete, "/?foo=test", nil)
	w := httptest.NewRecorder()

	handler.Delete(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("wrong status code, want: %d, got: %d", http.StatusNotFound, w.Code)
	}

	shouldContain := "key not present"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want: %s, got: %s", shouldContain, w.Body.String())
	}
}

func TestDeleteTimeout(t *testing.T) {
	handler := kvstorehandler.New(
		kvstorehandler.WithContextTimeout(time.Second*-1),
		kvstorehandler.WithService(&mockService{
			deleteErr: context.DeadlineExceeded,
		}),
	)

	req := httptest.NewRequest(http.MethodDelete, "/?key=test", nil)
	w := httptest.NewRecorder()

	handler.Delete(w, req)

	if w.Code != http.StatusGatewayTimeout {
		t.Errorf("wrong status code, want: %d, got: %d", http.StatusGatewayTimeout, w.Code)
	}

	shouldContain := "context deadline exceeded"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want: %s, got: %s", shouldContain, w.Body.String())
	}
}

func TestDeleteErrUnknown(t *testing.T) {
	handler := kvstorehandler.New(
		kvstorehandler.WithService(&mockService{
			deleteErr: kverror.ErrUnknown,
		}),
		kvstorehandler.WithLogger(logger),
	)

	req := httptest.NewRequest(http.MethodDelete, "/?key=test", nil)
	w := httptest.NewRecorder()

	handler.Delete(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("wrong status code, want: %d, got: %d", http.StatusInternalServerError, w.Code)
	}

	shouldContain := "unknown error"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want: %s, got: %s", shouldContain, w.Body.String())
	}
}

func TestDeleteErrKeyNotFound(t *testing.T) {
	_ = kverror.ErrKeyNotFound.AddData("key=test") // ignore error.

	handler := kvstorehandler.New(
		kvstorehandler.WithService(&mockService{
			deleteErr: kverror.ErrKeyNotFound,
		}),
		kvstorehandler.WithLogger(logger),
	)

	req := httptest.NewRequest(http.MethodDelete, "/?key=test", nil)
	w := httptest.NewRecorder()

	handler.Delete(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("wrong status code, want: %d, got: %d", http.StatusNotFound, w.Code)
	}

	if !strings.Contains(w.Body.String(), "key not found") {
		t.Error("body not equal")
	}

	shouldContain := "key=test"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want: %s, got: %s", shouldContain, w.Body.String())
	}

	_ = kverror.ErrKeyNotFound.DestoryData() // ignore error.
}

func TestDeleteSuccess(t *testing.T) {
	handler := kvstorehandler.New(
		kvstorehandler.WithService(&mockService{}),
		kvstorehandler.WithLogger(logger),
	)

	req := httptest.NewRequest(http.MethodDelete, "/?key=test", nil)
	w := httptest.NewRecorder()

	handler.Delete(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("wrong status code, want: %d, got: %d", http.StatusNoContent, w.Code)
	}

	if w.Body.Len() != 0 {
		t.Errorf("wrong body size, want: 0, got: %d", w.Body.Len())
	}
}
