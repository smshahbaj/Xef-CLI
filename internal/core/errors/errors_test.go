package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrap(t *testing.T) {
	original := errors.New("original error")
	wrapped := Wrap(original, "context")

	assert.Error(t, wrapped)
	assert.Contains(t, wrapped.Error(), "context")
	assert.Contains(t, wrapped.Error(), "original error")
}

func TestWrapNil(t *testing.T) {
	assert.Nil(t, Wrap(nil, "context"))
}

func TestNew(t *testing.T) {
	err := New("test error")
	assert.Error(t, err)
	assert.Equal(t, "test error", err.Error())
}

func TestWithCode(t *testing.T) {
	err := WithCode(ExitInvalidInput, "bad input")
	assert.Error(t, err)
	assert.Equal(t, ExitInvalidInput, err.(*AppError).Code)
}

func TestAppErrorUnwrap(t *testing.T) {
	cause := errors.New("cause")
	appErr := &AppError{Message: "wrapped", Cause: cause, Code: ExitGeneralError}
	assert.Equal(t, cause, appErr.Unwrap())
}

func TestIs(t *testing.T) {
	err1 := New("test")
	err2 := Wrap(err1, "wrapped")
	assert.True(t, Is(err2, err1))
}
