package kvstorehandler_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/vbyazilim/kvstore/src/internal/kverror"
	"github.com/vbyazilim/kvstore/src/internal/service/kvstoreservice"
	"github.com/vbyazilim/kvstore/src/internal/transport/http/kvstorehandler"
)

type errorReader struct{}

func (e *errorReader) Read(_ []byte) (n int, err error) {
	return 0, errors.New("forced error") // nolint
}

func TestSetInvalidMethod(t *testing.T) {
	handler := kvstorehandler.New()
	req := httptest.NewRequest(http.MethodDelete, "/key", nil)
	w := httptest.NewRecorder()

	handler.Set(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("wrong status code, want: %d, got: %d", http.StatusMethodNotAllowed, w.Code)
	}

	shouldContain := "method DELETE not allowed"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want: %s, got: %s", shouldContain, w.Body.String())
	}
}

func TestSetBodyReadError(t *testing.T) {
	handler := kvstorehandler.New()
	req := httptest.NewRequest(http.MethodPost, "/key", &errorReader{})

	w := httptest.NewRecorder()

	handler.Set(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("wrong status code, want: %d, got: %d", http.StatusBadRequest, w.Code)
	}
}

func TestSetBodyUnmarshal(t *testing.T) {
	handler := kvstorehandler.New()
	handlerRequest := bytes.NewBufferString(`{"key": "key", "value": "123}`)
	req := httptest.NewRequest(http.MethodPost, "/key", handlerRequest)
	w := httptest.NewRecorder()

	handler.Set(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("wrong status code, want: %d, got: %d", http.StatusInternalServerError, w.Code)
	}
}

func TestSetEmptyBody(t *testing.T) {
	handler := kvstorehandler.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	w := httptest.NewRecorder()

	handler.Set(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("wrong status code, want: %d, got: %d", http.StatusBadRequest, w.Code)
	}

	shouldContain := "empty body/payload"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want: %s, got: %s", shouldContain, w.Body.String())
	}
}

func TestSetKeyIsEmpty(t *testing.T) {
	handler := kvstorehandler.New()

	payload := strings.NewReader("{}")
	req := httptest.NewRequest(http.MethodPost, "/", payload)
	w := httptest.NewRecorder()

	handler.Set(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("wrong status code, want: %d, got: %d", http.StatusBadRequest, w.Code)
	}

	shouldContain := "key is empty"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want: %s, got: %s", shouldContain, w.Body.String())
	}
}

func TestSetValueIsEmpty(t *testing.T) {
	handler := kvstorehandler.New()

	payload := strings.NewReader(`{"key":"test"}`)
	req := httptest.NewRequest(http.MethodPost, "/", payload)
	w := httptest.NewRecorder()

	handler.Set(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("wrong status code, want: %d, got: %d", http.StatusBadRequest, w.Code)
	}

	shouldContain := "value is empty"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want: %s, got: %s", shouldContain, w.Body.String())
	}
}

func TestSetTimeout(t *testing.T) {
	handler := kvstorehandler.New(
		kvstorehandler.WithContextTimeout(time.Second*-1),
		kvstorehandler.WithService(&mockService{
			getErr: context.DeadlineExceeded,
		}),
	)

	payload := strings.NewReader(`{"key":"test","value":"test"}`)
	req := httptest.NewRequest(http.MethodPost, "/?key=test", payload)
	w := httptest.NewRecorder()

	handler.Set(w, req)

	if w.Code != http.StatusGatewayTimeout {
		t.Errorf("wrong status code, want: %d, got: %d", http.StatusGatewayTimeout, w.Code)
	}

	shouldContain := "context deadline exceeded"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want: %s, got: %s", shouldContain, w.Body.String())
	}
}

func TestSetErrUnknown(t *testing.T) {
	handler := kvstorehandler.New(
		kvstorehandler.WithService(&mockService{
			setErr: kverror.ErrUnknown,
		}),
		kvstorehandler.WithLogger(logger),
	)

	payload := strings.NewReader(`{"key":"test","value":"test"}`)
	req := httptest.NewRequest(http.MethodPost, "/", payload)
	w := httptest.NewRecorder()

	handler.Set(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("wrong status code, want: %d, got: %d", http.StatusInternalServerError, w.Code)
	}

	shouldContain := "unknown error"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want: %s, got: %s", shouldContain, w.Body.String())
	}
}

func TestSetServiceUnknownError(t *testing.T) {
	handler := kvstorehandler.New(
		kvstorehandler.WithService(&mockService{
			getErr: kverror.ErrUnknown.AddData("fake error"),
		}),
		kvstorehandler.WithLogger(logger),
	)

	payload := strings.NewReader(`{"key":"test","value":"test"}`)
	req := httptest.NewRequest(http.MethodPost, "/", payload)
	w := httptest.NewRecorder()

	handler.Set(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("wrong status code, want: %d, got: %d", http.StatusBadRequest, w.Code)
	}
}

func TestSetServiceNilExistingItem(t *testing.T) {
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

	payload := strings.NewReader(`{"key":"test","value":"test"}`)
	req := httptest.NewRequest(http.MethodPost, "/", payload)
	w := httptest.NewRecorder()

	handler.Set(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("wrong status code, want: %d, got: %d", http.StatusConflict, w.Code)
	}
}

func TestSetErrKeyExists(t *testing.T) {
	_ = kverror.ErrKeyExists.AddData("key=test") // ignore error.

	handler := kvstorehandler.New(
		kvstorehandler.WithService(&mockService{
			setErr: kverror.ErrKeyExists,
		}),
		kvstorehandler.WithLogger(logger),
	)

	payload := strings.NewReader(`{"key":"test","value":"test"}`)
	req := httptest.NewRequest(http.MethodPost, "/", payload)
	w := httptest.NewRecorder()

	handler.Set(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("wrong status code, want: %d, got: %d", http.StatusConflict, w.Code)
	}

	shouldContain := "key exist"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want: %s, got: %s", shouldContain, w.Body.String())
	}

	shouldContain = "key=test"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want: %s, got: %s", shouldContain, w.Body.String())
	}

	_ = kverror.ErrKeyExists.DestoryData() // ignore error.
}

func TestSetSuccess(t *testing.T) {
	handler := kvstorehandler.New(
		kvstorehandler.WithService(&mockService{}),
		kvstorehandler.WithLogger(logger),
		kvstorehandler.WithService(&mockService{
			setResponse: &kvstoreservice.ItemResponse{
				Key:   "test",
				Value: "test",
			},
		}),
	)

	payload := strings.NewReader(`{"key":"test","value":"test"}`)
	req := httptest.NewRequest(http.MethodPost, "/", payload)
	w := httptest.NewRecorder()

	handler.Set(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("wrong status code, want: %d, got: %d", http.StatusCreated, w.Code)
	}

	shouldEqual := `{"key":"test","value":"test"}`
	if w.Body.String() != shouldEqual {
		t.Errorf("wrong body message, want: %s, got: %s", shouldEqual, w.Body.String())
	}
}
