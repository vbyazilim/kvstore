package kverror_test

import (
	"errors"
	"testing"

	"github.com/vbyazilim/kvstore/src/internal/kverror"
)

func TestAddData(t *testing.T) {
	var errorData = "some data"

	// Create a new error
	err := kverror.New("some error", true)

	// Add data to error
	kvError := err.AddData(errorData)

	// Check if data is added
	if kvError.(*kverror.Error).Data != errorData {
		t.Error("data not added")
	}

	// Destroy data
	kvError.DestoryData()
}

func TestWrap(t *testing.T) {

	// Create a new error
	err := kverror.New("some error", true)
	wrappedErr := errors.New("wrapped error")

	// Wrap error
	kvError := err.Wrap(wrappedErr)

	// Check if error is wrapped
	if kvError.(*kverror.Error).Err != wrappedErr {
		t.Error("error not wrapped")
	}
}

func TestUnwrap(t *testing.T) {

	// Create a new error
	err := kverror.New("some error", true)
	wrappedErr := errors.New("wrapped error")

	// Wrap error
	kvError := err.Wrap(wrappedErr)

	// Check if error is wrapped
	if kvError.Unwrap() != wrappedErr {
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
	if kvError.(*kverror.Error).Data != errorData {
		t.Error("data not added")
	}

	// Destroy data
	kvError.DestoryData()

	// Check if data is destroyed
	if kvError.(*kverror.Error).Data != nil {
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
	if kvError.Error() != "some error" {
		t.Error("error message is not correct")
	}

	// Wrap error
	wrappedErr := errors.New("wrapped error")
	kvError = kvError.Wrap(wrappedErr)

	// Check if error message is correct
	if kvError.Error() != "wrapped error, some error" {
		t.Error("error message is not correct")
	}
}
