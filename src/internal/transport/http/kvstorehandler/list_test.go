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

func TestListInvalidMethod(t *testing.T) {
	handler := kvstorehandler.New()
	req := httptest.NewRequest("DELETE", "/key", nil)
	w := httptest.NewRecorder()

	handler.List(w, req)

	if w.Code != 405 && strings.Contains(w.Body.String(), "method not allowed") {
		t.Error("code not equal")
	}
}

func TestListTimeout(t *testing.T) {
	handler := kvstorehandler.New(
		kvstorehandler.WithContextTimeout(time.Second*-1),
		kvstorehandler.WithService(&mockService{
			listErr: context.DeadlineExceeded,
		}),
	)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	handler.List(w, req)

	if w.Code != 500 {
		t.Error("code not equal")
	}

	if !strings.Contains(w.Body.String(), "context deadline exceeded") {
		t.Error("body not equal")
	}
}

func TestListErrUnknown(t *testing.T) {
	handler := kvstorehandler.New(
		kvstorehandler.WithService(&mockService{
			listErr: kverror.ErrUnknown,
		}),
		kvstorehandler.WithLogger(logger),
	)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	handler.List(w, req)

	if w.Code != 500 {
		t.Error("code not equal")
	}

	if !strings.Contains(w.Body.String(), "unknown error") {
		t.Error("body not equal")
	}
}

func TestList(t *testing.T) {
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

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	handler.List(w, req)

	if w.Code != 200 {
		t.Error("code not equal")
	}

	if w.Body.String() != `[{"key":"test","value":"test"}]` {
		t.Error("body not equal")
	}
}
