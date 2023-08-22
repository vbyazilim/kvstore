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

func TestListInvalidMethod(t *testing.T) {
	handler := kvstorehandler.New()
	req := httptest.NewRequest(http.MethodDelete, "/key", nil)
	w := httptest.NewRecorder()

	handler.List(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("wrong status code, want: %d, got: %d", http.StatusMethodNotAllowed, w.Code)
	}

	shouldContain := "method DELETE not allowed"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want: %s, got: %s", shouldContain, w.Body.String())
	}
}

func TestListTimeout(t *testing.T) {
	handler := kvstorehandler.New(
		kvstorehandler.WithContextTimeout(time.Second*-1),
		kvstorehandler.WithService(&mockService{
			listErr: context.DeadlineExceeded,
		}),
	)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.List(w, req)

	if w.Code != http.StatusGatewayTimeout {
		t.Errorf("wrong status code, want: %d, got: %d", http.StatusGatewayTimeout, w.Code)
	}

	shouldContain := "context deadline exceeded"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want: %s, got: %s", shouldContain, w.Body.String())
	}
}

func TestListErrUnknown(t *testing.T) {
	handler := kvstorehandler.New(
		kvstorehandler.WithService(&mockService{
			listErr: kverror.ErrUnknown.AddData("fake error"),
		}),
		kvstorehandler.WithLogger(logger),
	)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.List(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("wrong status code, want: %d, got: %d", http.StatusInternalServerError, w.Code)
	}

	shouldContain := "unknown error"
	if !strings.Contains(w.Body.String(), shouldContain) {
		t.Errorf("wrong body message, want: %s, got: %s", shouldContain, w.Body.String())
	}
}

func TestEmptyList(t *testing.T) {
	handler := kvstorehandler.New(
		kvstorehandler.WithService(&mockService{}),
		kvstorehandler.WithLogger(logger),
		kvstorehandler.WithService(&mockService{
			listResponse: &kvstoreservice.ListResponse{},
		}),
	)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.List(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("wrong status code, want: %d, got: %d", http.StatusNotFound, w.Code)
	}
}

func TestListSuccess(t *testing.T) {
	handler := kvstorehandler.New(
		kvstorehandler.WithService(&mockService{}),
		kvstorehandler.WithLogger(logger),
		kvstorehandler.WithService(&mockService{
			listResponse: &kvstoreservice.ListResponse{
				{
					Key:   "test",
					Value: "test",
				},
			},
		}),
	)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.List(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("wrong status code, want: %d, got: %d", http.StatusOK, w.Code)
	}

	shouldEqual := `[{"key":"test","value":"test"}]`
	if w.Body.String() != shouldEqual {
		t.Errorf("wrong body message, want: %s, got: %s", shouldEqual, w.Body.String())
	}
}
