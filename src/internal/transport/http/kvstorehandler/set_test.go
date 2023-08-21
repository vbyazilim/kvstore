package kvstorehandler_test

import (
	"context"
	"log/slog"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/vbyazilim/kvstore/src/internal/kverror"
	"github.com/vbyazilim/kvstore/src/internal/service/kvstoreservice"
	"github.com/vbyazilim/kvstore/src/internal/transport/http/kvstorehandler"
)

func TestSetInvalidMethod(t *testing.T) {
	handler := kvstorehandler.New()
	req := httptest.NewRequest("DELETE", "/key", nil)
	w := httptest.NewRecorder()

	handler.Set(w, req)

	if w.Code != 405 {
		t.Error("code not equal")
	}

	if !strings.Contains(w.Body.String(), "method DELETE not allowed") {
		t.Error("body not equal")
	}
}

func TestSetEmptyBody(t *testing.T) {
	handler := kvstorehandler.New()
	req := httptest.NewRequest("POST", "/", nil)
	w := httptest.NewRecorder()

	handler.Set(w, req)

	if w.Code != 400 {
		t.Error("code not equal")
	}

	if !strings.Contains(w.Body.String(), "empty body/payload") {
		t.Error("body not equal")
	}
}

func TestSetKeyIsEmpty(t *testing.T) {
	handler := kvstorehandler.New()

	req := httptest.NewRequest("POST", "/", strings.NewReader("{}"))
	w := httptest.NewRecorder()

	handler.Set(w, req)

	if w.Code != 400 {
		t.Error("code not equal")
	}

	if !strings.Contains(w.Body.String(), "key is empty") {
		t.Error("body not equal")
	}
}

func TestSetValueIsEmpty(t *testing.T) {
	handler := kvstorehandler.New()

	req := httptest.NewRequest("POST", "/",
		strings.NewReader(`{"key":"test"}`))
	w := httptest.NewRecorder()

	handler.Set(w, req)

	if w.Code != 400 {
		t.Error("code not equal")
	}

	if !strings.Contains(w.Body.String(), "value is empty") {
		t.Error("body not equal")
	}
}

func TestSetTimeout(t *testing.T) {
	handler := kvstorehandler.New(
		kvstorehandler.WithContextTimeout(time.Second*-1),
		kvstorehandler.WithService(&mockService{
			getErr: context.DeadlineExceeded,
		}),
	)

	req := httptest.NewRequest("POST", "/?key=test",
		strings.NewReader(`{"key":"test","value":"test"}`))
	w := httptest.NewRecorder()

	handler.Set(w, req)

	if w.Code != 409 {
		t.Error("code not equal")
	}

	if !strings.Contains(w.Body.String(), "can not set, 'test' already exists") {
		t.Error("body not equal")
	}
}

func TestSetErrUnknown(t *testing.T) {
	logger := slog.Default()
	handler := kvstorehandler.New(
		kvstorehandler.WithService(&mockService{
			setErr: kverror.ErrUnknown,
		}),
		kvstorehandler.WithLogger(logger),
	)

	req := httptest.NewRequest("POST", "/",
		strings.NewReader(`{"key":"test","value":"test"}`))
	w := httptest.NewRecorder()

	handler.Set(w, req)

	if w.Code != 500 {
		t.Error("code not equal")
	}

	if !strings.Contains(w.Body.String(), "unknown error") {
		t.Error("body not equal")
	}
}

func TestSetErrKeyExists(t *testing.T) {
	logger := slog.Default()

	// ignore error.
	_ = kverror.ErrKeyExists.AddData("key=test")

	handler := kvstorehandler.New(
		kvstorehandler.WithService(&mockService{
			setErr: kverror.ErrKeyExists,
		}),
		kvstorehandler.WithLogger(logger),
	)

	req := httptest.NewRequest("POST", "/",
		strings.NewReader(`{"key":"test","value":"test"}`))
	w := httptest.NewRecorder()

	handler.Set(w, req)

	if w.Code != 409 {
		t.Error("code not equal")
	}

	if !strings.Contains(w.Body.String(), "key exist") {
		t.Error("body not equal")
	}

	if !strings.Contains(w.Body.String(), "key=test") {
		t.Error("body not equal")
	}

	_ = kverror.ErrKeyExists.DestoryData()
}

func TestSet(t *testing.T) {
	logger := slog.Default()

	handler := kvstorehandler.New(
		kvstorehandler.WithService(&mockService{}),
		kvstorehandler.WithLogger(logger),
		kvstorehandler.WithService(&mockService{
			getResponse: &kvstoreservice.ItemResponse{
				Key:   "test",
				Value: "test",
			},
			setResponse: &kvstoreservice.ItemResponse{
				Key:   "test",
				Value: "test",
			},
		}),
	)

	req := httptest.NewRequest("POST", "/",
		strings.NewReader(`{"key":"test","value":"test"}`))
	w := httptest.NewRecorder()

	handler.Set(w, req)

	if w.Code != 201 {
		t.Error("code not equal")
	}

	if w.Body.String() != `{"key":"test","value":"test"}` {
		t.Error("body not equal")
	}
}
