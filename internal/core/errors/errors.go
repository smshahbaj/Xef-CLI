// Package errors provides custom error types and utilities for XefCLI.
// All errors are wrapped with context to aid debugging.
package errors

import (
	"errors"
	"fmt"
)

// ExitCode represents a process exit code.
type ExitCode int

const (
	ExitSuccess       ExitCode = 0
	ExitGeneralError  ExitCode = 1
	ExitInvalidInput  ExitCode = 2
	ExitNotFound      ExitCode = 3
	ExitPermission    ExitCode = 4
	ExitTimeout       ExitCode = 5
	ExitInterrupted   ExitCode = 130
)

// AppError is the base application error type.
type AppError struct {
	Code    ExitCode
	Message string
	Cause   error
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Cause
}

// Wrap wraps an error with a message.
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	return &AppError{Message: msg, Cause: err, Code: ExitGeneralError}
}

// Wrapf wraps an error with a formatted message.
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return &AppError{Message: fmt.Sprintf(format, args...), Cause: err, Code: ExitGeneralError}
}

// New creates a new error.
func New(msg string) error {
	return &AppError{Message: msg, Code: ExitGeneralError}
}

// Newf creates a new formatted error.
func Newf(format string, args ...interface{}) error {
	return &AppError{Message: fmt.Sprintf(format, args...), Code: ExitGeneralError}
}

// WithCode creates an error with a specific exit code.
func WithCode(code ExitCode, msg string) error {
	return &AppError{Message: msg, Code: code}
}

// Is reports whether any error in err's tree matches target.
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As finds the first error in err's tree that matches target.
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// Common errors.
var (
	ErrNotFound     = New("not found")
	ErrInvalidInput = WithCode(ExitInvalidInput, "invalid input")
	ErrPermission   = WithCode(ExitPermission, "permission denied")
	ErrTimeout      = WithCode(ExitTimeout, "operation timed out")
	ErrCanceled     = WithCode(ExitInterrupted, "operation canceled")
)
