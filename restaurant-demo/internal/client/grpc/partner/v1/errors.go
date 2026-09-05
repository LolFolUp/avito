package v1

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/errs"
)

func clientError(err error) error {
	grpcStatus, ok := status.FromError(err)
	if !ok {
		return &errs.Error{Kind: errs.KindInternal, Message: "partner request failed", Cause: err}
	}

	kind := errs.KindInternal
	switch grpcStatus.Code() {
	case codes.InvalidArgument:
		kind = errs.KindInvalidArgument
	case codes.NotFound:
		kind = errs.KindNotFound
	case codes.FailedPrecondition:
		kind = errs.KindFailedPrecondition
	case codes.AlreadyExists:
		kind = errs.KindAlreadyExists
	case codes.Unavailable, codes.DeadlineExceeded:
		kind = errs.KindUnavailable
	}

	return &errs.Error{Kind: kind, Message: grpcStatus.Message(), Cause: err}
}
