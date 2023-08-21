package kvstorehandler_test

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/vbyazilim/kvstore/src/internal/kverror"
	"github.com/vbyazilim/kvstore/src/internal/service/kvstoreservice"
	"github.com/vbyazilim/kvstore/src/internal/transport/http/kvstorehandler"
)

func TestUpdateInvalidMethod(t *testing.T) {
	handler := kvstorehandler.New()
	req := httptest.NewRequest("DELETE", "/key", nil)
	w := httptest.NewRecorder()

	handler.Update(w, req)

	if w.Code != 405 {
		t.Error("code not equal")
	}

	if !strings.Contains(w.Body.String(), "method DELETE not allowed") {
		t.Error("body not equal")
	}
}

func TestUpdateEmptyBody(t *testing.T) {
	handler := kvstorehandler.New()
	req := httptest.NewRequest("PUT", "/", nil)
	w := httptest.NewRecorder()

	handler.Update(w, req)

	if w.Code != 400 {
		t.Error("code not equal")
	}

	if !strings.Contains(w.Body.String(), "empty body/payload") {
		t.Error("body not equal")
	}
}

func TestUpdateKeyIsEmpty(t *testing.T) {
	handler := kvstorehandler.New()

	req := httptest.NewRequest("PUT", "/", strings.NewReader("{}"))
	w := httptest.NewRecorder()

	handler.Update(w, req)

	if w.Code != 400 {
		t.Error("code not equal")
	}

	if !strings.Contains(w.Body.String(), "key is empty") {
		t.Error("body not equal")
	}
}

func TestUpdateValueIsEmpty(t *testing.T) {
	handler := kvstorehandler.New()

	req := httptest.NewRequest("PUT", "/",
		strings.NewReader(`{"key":"test"}`))
	w := httptest.NewRecorder()

	handler.Update(w, req)

	if w.Code != 400 {
		t.Error("code not equal")
	}

	if !strings.Contains(w.Body.String(), "value is empty") {
		t.Error("body not equal")
	}
}

func TestUpdateTimeout(t *testing.T) {
	handler := kvstorehandler.New(
		kvstorehandler.WithContextTimeout(time.Second*-1),
		kvstorehandler.WithService(&mockService{
			updateErr: context.DeadlineExceeded,
		}),
	)

	req := httptest.NewRequest("PUT", "/?key=test",
		strings.NewReader(`{"key":"test","value":"test"}`))
	w := httptest.NewRecorder()

	handler.Update(w, req)

	if w.Code != 500 {
		t.Error("code not equal")
	}

	if !strings.Contains(w.Body.String(), "context deadline exceeded") {
		t.Error("body not equal")
	}
}

func TestUpdateErrUnknown(t *testing.T) {
	handler := kvstorehandler.New(
		kvstorehandler.WithService(&mockService{
			updateErr: kverror.ErrUnknown,
		}),
		kvstorehandler.WithLogger(logger),
	)

	req := httptest.NewRequest("PUT", "/",
		strings.NewReader(`{"key":"test","value":"test"}`))
	w := httptest.NewRecorder()

	handler.Update(w, req)

	if w.Code != 500 {
		t.Error("code not equal")
	}

	if !strings.Contains(w.Body.String(), "unknown error") {
		t.Error("body not equal")
	}
}

func TestUpdateErrKeyExists(t *testing.T) {
	// ignore return no need
	_ = kverror.ErrKeyNotFound.AddData("key=test")

	handler := kvstorehandler.New(
		kvstorehandler.WithService(&mockService{
			updateErr: kverror.ErrKeyNotFound,
		}),
		kvstorehandler.WithLogger(logger),
	)

	req := httptest.NewRequest("PUT", "/",
		strings.NewReader(`{"key":"test","value":"test"}`))
	w := httptest.NewRecorder()

	handler.Update(w, req)

	if w.Code != 404 {
		t.Error("code not equal")
	}

	if !strings.Contains(w.Body.String(), "key not found") {
		t.Error("body not equal")
	}

	if !strings.Contains(w.Body.String(), "key=test") {
		t.Error("body not equal")
	}

	// ignore error
	_ = kverror.ErrKeyNotFound.DestoryData()
}

func TestUpdate(t *testing.T) {
	handler := kvstorehandler.New(
		kvstorehandler.WithService(&mockService{}),
		kvstorehandler.WithLogger(logger),
		kvstorehandler.WithService(&mockService{
			updateResponse: &kvstoreservice.ItemResponse{
				Key:   "test",
				Value: "test",
			},
		}),
	)

	req := httptest.NewRequest("PUT", "/",
		strings.NewReader(`{"key":"test","value":"test"}`))
	w := httptest.NewRecorder()

	handler.Update(w, req)

	if w.Code != 200 {
		t.Error("code not equal")
	}

	if w.Body.String() != `{"key":"test","value":"test"}` {
		t.Error("body not equal")
	}
}
