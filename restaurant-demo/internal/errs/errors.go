package errs

import (
	"errors"
	"fmt"
)

type Kind string

const (
	KindInvalidArgument    Kind = "invalid_argument"
	KindNotFound           Kind = "not_found"
	KindFailedPrecondition Kind = "failed_precondition"
	KindAlreadyExists      Kind = "already_exists"
	KindUnavailable        Kind = "unavailable"
	KindInternal           Kind = "internal"
)

type Error struct {
	Kind    Kind
	Message string
	Cause   error
}

func InvalidArgument(message string) error {
	return &Error{Kind: KindInvalidArgument, Message: message}
}

func NotFound(message string) error {
	return &Error{Kind: KindNotFound, Message: message}
}

func FailedPrecondition(message string) error {
	return &Error{Kind: KindFailedPrecondition, Message: message}
}

func AlreadyExists(message string) error {
	return &Error{Kind: KindAlreadyExists, Message: message}
}

func Unavailable(message string, cause error) error {
	return &Error{Kind: KindUnavailable, Message: message, Cause: cause}
}

func Internal(message string, cause error) error {
	return &Error{Kind: KindInternal, Message: message, Cause: cause}
}

func As(err error) (*Error, bool) {
	var domainError *Error
	if !errors.As(err, &domainError) {
		return nil, false
	}
	return domainError, true
}

func (e *Error) Error() string {
	if e.Cause == nil {
		return e.Message
	}
	return fmt.Sprintf("%s: %v", e.Message, e.Cause)
}

func (e *Error) Unwrap() error {
	return e.Cause
}
