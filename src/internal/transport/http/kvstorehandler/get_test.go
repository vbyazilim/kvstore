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
	req := httptest.NewRequest(http.MethodDelete, "/key", nil)
	w := httptest.NewRecorder()

	handler.Get(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("wrong status code, want: %d, got: %d", http.StatusMethodNotAllowed, w.Code)
	}

	shouldContain := "method DELETE not allowed"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want: %s, got: %s", shouldContain, w.Body.String())
	}
}

func TestGetQueryParamRequired(t *testing.T) {
	handler := kvstorehandler.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.Get(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("wrong status code, want: %d, got: %d", http.StatusNotFound, w.Code)
	}

	shouldContain := "key query param required"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want: %s, got: %s", shouldContain, w.Body.String())
	}
}

func TestGetQueryParamKeyNotFound(t *testing.T) {
	handler := kvstorehandler.New()
	req := httptest.NewRequest(http.MethodGet, "/?foo=test", nil)
	w := httptest.NewRecorder()

	handler.Get(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("wrong status code, want: %d, got: %d", http.StatusNotFound, w.Code)
	}

	shouldContain := "key not present"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want: %s, got: %s", shouldContain, w.Body.String())
	}
}

func TestGetTimeout(t *testing.T) {
	handler := kvstorehandler.New(
		kvstorehandler.WithContextTimeout(time.Second*-1),
		kvstorehandler.WithService(&mockService{
			getErr: context.DeadlineExceeded,
		}),
	)

	req := httptest.NewRequest(http.MethodGet, "/?key=test", nil)
	w := httptest.NewRecorder()

	handler.Get(w, req)

	if w.Code != http.StatusGatewayTimeout {
		t.Errorf("wrong status code, want: %d, got: %d", http.StatusGatewayTimeout, w.Code)
	}

	shouldContain := "context deadline exceeded"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want: %s, got: %s", shouldContain, w.Body.String())
	}
}

func TestGetErrUnknown(t *testing.T) {
	handler := kvstorehandler.New(
		kvstorehandler.WithService(&mockService{
			getErr: kverror.ErrUnknown,
		}),
		kvstorehandler.WithLogger(logger),
	)

	req := httptest.NewRequest(http.MethodGet, "/?key=test", nil)
	w := httptest.NewRecorder()

	handler.Get(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("wrong status code, want: %d, got: %d", http.StatusInternalServerError, w.Code)
	}

	shouldContain := "unknown error"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want: %s, got: %s", shouldContain, w.Body.String())
	}
}

func TestGetErrKeyNotFound(t *testing.T) {
	_ = kverror.ErrKeyNotFound.AddData("key=test") // ignore error.

	handler := kvstorehandler.New(
		kvstorehandler.WithService(&mockService{
			getErr: kverror.ErrKeyNotFound,
		}),
		kvstorehandler.WithLogger(logger),
	)

	req := httptest.NewRequest(http.MethodGet, "/?key=test", nil)
	w := httptest.NewRecorder()

	handler.Get(w, req)

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

	_ = kverror.ErrKeyNotFound.DestoryData() // ignore error.
}

func TestGetSuccess(t *testing.T) {
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

	req := httptest.NewRequest(http.MethodGet, "/?key=test", nil)
	w := httptest.NewRecorder()

	handler.Get(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("wrong status code, want: %d, got: %d", http.StatusOK, w.Code)
	}

	shouldEqual := `{"key":"test","value":"test"}`
	if w.Body.String() != shouldEqual {
		t.Errorf("wrong body message, want: %s, got: %s", shouldEqual, w.Body.String())
	}
}
