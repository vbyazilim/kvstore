package kvstorehandler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/vbyazilim/kvstore/src/internal/kverror"
	"github.com/vbyazilim/kvstore/src/internal/service/kvstoreservice"
	"github.com/vbyazilim/kvstore/src/internal/transport/http/kvstorehandler"
)

func TestGetInvalidMethod(t *testing.T) {
	handler := kvstorehandler.New()
	req := httptest.NewRequest("DELETE", "/key", nil)
	w := httptest.NewRecorder()

	handler.Get(w, req)

	if w.Code != http.StatusMethodNotAllowed && strings.Contains(w.Body.String(), "method not allowed") {
		t.Error("code not equal")
	}
}

func TestGetQueryParamRequired(t *testing.T) {
	handler := kvstorehandler.New()
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	handler.Get(w, req)

	if w.Code != http.StatusNotFound && strings.Contains(w.Body.String(), "key query param required") {
		t.Error("code not equal")
	}
}

func TestGetQueryParamKeyNotFound(t *testing.T) {
	handler := kvstorehandler.New()
	req := httptest.NewRequest("GET", "/?foo=test", nil)
	w := httptest.NewRecorder()

	handler.Get(w, req)

	if w.Code != http.StatusNotFound {
		t.Error("code not equal")
	}

	if !strings.Contains(w.Body.String(), "key not present") {
		t.Error("body not equal")
	}
}

func TestGetTimeout(t *testing.T) {
	handler := kvstorehandler.New(
		kvstorehandler.WithContextTimeout(time.Second*-1),
		kvstorehandler.WithService(&mockService{
			getErr: context.DeadlineExceeded,
		}),
	)

	req := httptest.NewRequest("GET", "/?key=test", nil)
	w := httptest.NewRecorder()

	handler.Get(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Error("code not equal")
	}

	if !strings.Contains(w.Body.String(), "context deadline exceeded") {
		t.Error("body not equal")
	}
}

func TestGetErrUnknown(t *testing.T) {
	handler := kvstorehandler.New(
		kvstorehandler.WithService(&mockService{
			getErr: kverror.ErrUnknown,
		}),
		kvstorehandler.WithLogger(logger),
	)

	req := httptest.NewRequest("GET", "/?key=test", nil)
	w := httptest.NewRecorder()

	handler.Get(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Error("code not equal")
	}

	if !strings.Contains(w.Body.String(), "unknown error") {
		t.Error("body not equal")
	}
}

func TestGetErrKeyNotFound(t *testing.T) {
	// ignore error.
	_ = kverror.ErrKeyNotFound.AddData("key=test")

	handler := kvstorehandler.New(
		kvstorehandler.WithService(&mockService{
			getErr: kverror.ErrKeyNotFound,
		}),
		kvstorehandler.WithLogger(logger),
	)

	req := httptest.NewRequest("GET", "/?key=test", nil)
	w := httptest.NewRecorder()

	handler.Get(w, req)

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

func TestGet(t *testing.T) {
	handler := kvstorehandler.New(
		kvstorehandler.WithService(&mockService{}),
		kvstorehandler.WithLogger(logger),
		kvstorehandler.WithService(&mockService{
			getResponse: &kvstoreservice.ItemResponse{
				Key:   "test",
				Value: "test",
			},
		}),
	)

	req := httptest.NewRequest("GET", "/?key=test", nil)
	w := httptest.NewRecorder()

	handler.Get(w, req)

	if w.Code != http.StatusOK {
		t.Error("code not equal")
	}

	if w.Body.String() != `{"key":"test","value":"test"}` {
		t.Error("body not equal")
	}
}
