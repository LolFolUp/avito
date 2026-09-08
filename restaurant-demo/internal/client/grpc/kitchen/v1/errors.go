package v1

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/restaurant-demo/internal/errs"
)

func clientError(operation string, err error) error {
	grpcStatus, ok := status.FromError(err)
	if !ok {
		return errs.Internal(operation, err)
	}

	switch grpcStatus.Code() {
	case codes.InvalidArgument:
		return errs.InvalidArgument(grpcStatus.Message())
	case codes.NotFound:
		return errs.NotFound(grpcStatus.Message())
	case codes.FailedPrecondition:
		return errs.FailedPrecondition(grpcStatus.Message())
	case codes.AlreadyExists:
		return errs.AlreadyExists(grpcStatus.Message())
	case codes.Unavailable, codes.DeadlineExceeded:
		return errs.Unavailable(operation+" is unavailable", err)
	default:
		return errs.Internal(operation, err)
	}
}
