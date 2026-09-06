package errs

import (
	"errors"
	"fmt"
)

type Error struct {
	Kind    Kind
	Code    string
	Message string
	Details any
	cause   error
}

func (e *Error) Error() string {
	if e.cause == nil {
		return e.Message
	}

	return fmt.Sprintf("%s: %v", e.Message, e.cause)
}

func (e *Error) Unwrap() error {
	return e.cause
}

func New(kind Kind, code, message string, details any) error {
	return &Error{Kind: kind, Code: code, Message: message, Details: details}
}

func Wrap(kind Kind, code, message string, cause error) error {
	return &Error{Kind: kind, Code: code, Message: message, cause: cause}
}

func InvalidArgument(message string) error {
	return New(KindInvalidArgument, CodeInvalidArgument, message, nil)
}

func NotFound(resource string) error {
	return New(KindNotFound, CodeNotFound, resource+" not found", nil)
}

func FailedPrecondition(code, message string, details any) error {
	return New(KindFailedPrecondition, code, message, details)
}

func AlreadyExists(resource string) error {
	return New(KindAlreadyExists, CodeAlreadyExists, resource+" already exists", nil)
}

func Unavailable(operation string, cause error) error {
	return Wrap(KindUnavailable, CodeUnavailable, operation+" is unavailable", cause)
}

func Internal(operation string, cause error) error {
	return Wrap(KindInternal, CodeInternal, operation+" failed", cause)
}

func As(err error) (*Error, bool) {
	var domainErr *Error
	if !errors.As(err, &domainErr) {
		return nil, false
	}

	return domainErr, true
}
