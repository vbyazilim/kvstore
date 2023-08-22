package kvstorehandler_test

import (
	"bytes"
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

func TestUpdateInvalidMethod(t *testing.T) {
	handler := kvstorehandler.New()
	req := httptest.NewRequest(http.MethodDelete, "/key", nil)
	w := httptest.NewRecorder()

	handler.Update(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("wrong status code, want: %d, got: %d", http.StatusMethodNotAllowed, w.Code)
	}

	shouldContain := "method DELETE not allowed"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want: %s, got: %s", shouldContain, w.Body.String())
	}
}

func TestUpdateBodyReadError(t *testing.T) {
	handler := kvstorehandler.New()
	req := httptest.NewRequest(http.MethodPut, "/", &errorReader{})

	w := httptest.NewRecorder()

	handler.Update(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("wrong status code, want: %d, got: %d", http.StatusBadRequest, w.Code)
	}
}

func TestUpdateEmptyBody(t *testing.T) {
	handler := kvstorehandler.New()
	req := httptest.NewRequest(http.MethodPut, "/", nil)
	w := httptest.NewRecorder()

	handler.Update(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("wrong status code, want: %d, got: %d", http.StatusBadRequest, w.Code)
	}

	shouldContain := "empty body/payload"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want: %s, got: %s", shouldContain, w.Body.String())
	}
}

func TestUpdateBodyUnmarshal(t *testing.T) {
	handler := kvstorehandler.New()
	handlerRequest := bytes.NewBufferString(`{"key": "key", "value": "123}`)
	req := httptest.NewRequest(http.MethodPut, "/", handlerRequest)
	w := httptest.NewRecorder()

	handler.Update(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("wrong status code, want: %d, got: %d", http.StatusInternalServerError, w.Code)
	}
}

func TestUpdateKeyIsEmpty(t *testing.T) {
	handler := kvstorehandler.New()

	payload := strings.NewReader("{}")
	req := httptest.NewRequest(http.MethodPut, "/", payload)
	w := httptest.NewRecorder()

	handler.Update(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("wrong status code, want: %d, got: %d", http.StatusBadRequest, w.Code)
	}

	shouldContain := "key is empty"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want: %s, got: %s", shouldContain, w.Body.String())
	}
}

func TestUpdateValueIsEmpty(t *testing.T) {
	handler := kvstorehandler.New()

	payload := strings.NewReader(`{"key":"test"}`)
	req := httptest.NewRequest(http.MethodPut, "/", payload)
	w := httptest.NewRecorder()

	handler.Update(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("wrong status code, want: %d, got: %d", http.StatusBadRequest, w.Code)
	}

	shouldContain := "value is empty"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want: %s, got: %s", shouldContain, w.Body.String())
	}
}

func TestUpdateTimeout(t *testing.T) {
	handler := kvstorehandler.New(
		kvstorehandler.WithContextTimeout(time.Second*-1),
		kvstorehandler.WithService(&mockService{
			updateErr: context.DeadlineExceeded,
		}),
	)

	payload := strings.NewReader(`{"key":"test","value":"test"}`)
	req := httptest.NewRequest(http.MethodPut, "/?key=test", payload)
	w := httptest.NewRecorder()

	handler.Update(w, req)

	if w.Code != http.StatusGatewayTimeout {
		t.Errorf("wrong status code, want: %d, got: %d", http.StatusGatewayTimeout, w.Code)
	}

	shouldContain := "context deadline exceeded"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want: %s, got: %s", shouldContain, w.Body.String())
	}
}

func TestUpdateErrUnknown(t *testing.T) {
	handler := kvstorehandler.New(
		kvstorehandler.WithService(&mockService{
			updateErr: kverror.ErrUnknown,
		}),
		kvstorehandler.WithLogger(logger),
	)

	payload := strings.NewReader(`{"key":"test","value":"test"}`)
	req := httptest.NewRequest(http.MethodPut, "/", payload)
	w := httptest.NewRecorder()

	handler.Update(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("wrong status code, want: %d, got: %d", http.StatusInternalServerError, w.Code)
	}

	shouldContain := "unknown error"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want: %s, got: %s", shouldContain, w.Body.String())
	}
}

func TestUpdateErrKeyExists(t *testing.T) {
	_ = kverror.ErrKeyNotFound.AddData("key=test") // ignore return no need

	handler := kvstorehandler.New(
		kvstorehandler.WithService(&mockService{
			updateErr: kverror.ErrKeyNotFound,
		}),
		kvstorehandler.WithLogger(logger),
	)

	payload := strings.NewReader(`{"key":"test","value":"test"}`)
	req := httptest.NewRequest(http.MethodPut, "/", payload)
	w := httptest.NewRecorder()

	handler.Update(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("wrong status code, want: %d, got: %d", http.StatusNotFound, w.Code)
	}

	shouldContain := "key not found"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want: %s, got: %s", shouldContain, w.Body.String())
	}

	shouldContain = "key=test"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want: %s, got: %s", shouldContain, w.Body.String())
	}

	_ = kverror.ErrKeyNotFound.DestoryData() // ignore error
}

func TestUpdateSuccess(t *testing.T) {
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

	payload := strings.NewReader(`{"key":"test","value":"test"}`)
	req := httptest.NewRequest(http.MethodPut, "/", payload)
	w := httptest.NewRecorder()

	handler.Update(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("wrong status code, want: %d, got: %d", http.StatusOK, w.Code)
	}

	shouldEqual := `{"key":"test","value":"test"}`
	if w.Body.String() != shouldEqual {
		t.Errorf("wrong body message, want: %s, got: %s", shouldEqual, w.Body.String())
	}
}
