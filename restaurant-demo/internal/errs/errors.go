package errs

import "fmt"

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

func (e *Error) Error() string {
	if e.Cause == nil {
		return e.Message
	}
	return fmt.Sprintf("%s: %v", e.Message, e.Cause)
}

func (e *Error) Unwrap() error {
	return e.Cause
}
