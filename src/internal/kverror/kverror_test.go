package kverror_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/vbyazilim/kvstore/src/internal/kverror"
)

var errWrapped = errors.New("wrapped error")
var err = kverror.New("some error", true)

func TestAddData(t *testing.T) {
	var errorData = "some data"

	// Add data to error
	kvError := err.AddData(errorData)

	// Check if data is added
	if !strings.Contains(kvError.Error(), errorData) == false {
		t.Error("data not added")
	}

	// Destroy data
	kvError.DestoryData()
}

func TestWrap(t *testing.T) {
	// Wrap error
	kvError := err.Wrap(errWrapped)

	// Check if error is wrapped
	if kvError.Error() != "wrapped error, some error" {
		t.Error("error not wrapped")
	}
}

func TestUnwrap(t *testing.T) {

	// Wrap error
	kvError := err.Wrap(errWrapped)

	// Check if error is wrapped
	if !errors.Is(kvError.Unwrap(), errWrapped) {
		t.Error("error not unwrapped")
	}
}

func TestDestoryData(t *testing.T) {
	var errorData = "some data"

	// Create a new error
	err := kverror.New("some error", true)

	// Add data to error
	kvError := err.AddData(errorData)

	// Check if data is added
	if kvError.GetData() != errorData {
		t.Error("data not added")
	}

	// Destroy data
	kvError.DestoryData()

	// Check if data is destroyed
	if kvError.GetData() != nil {
		t.Error("data not destroyed")
	}
}

func TestError(t *testing.T) {
	var errorData = "some data"

	// Create a new error
	err := kverror.New("some error", true)

	// Add data to error
	kvError := err.AddData(errorData)

	// Check if error message is correct
	if kvError.GetData() != errorData {
		t.Error("error message is not correct")
	}

	// Wrap error
	kvError = kvError.Wrap(errWrapped)

	// Check if error message is correct
	if kvError.Error() != "wrapped error, some error" {
		t.Error("error message is not correct")
	}
}
