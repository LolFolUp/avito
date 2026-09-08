package v1

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/errs"
)

func grpcError(err error) error {
	if err == nil {
		return nil
	}

	domainError, ok := errs.As(err)
	if !ok {
		return status.Error(codes.Internal, "internal error")
	}

	switch domainError.Kind {
	case errs.KindInvalidArgument:
		return status.Error(codes.InvalidArgument, domainError.Message)
	case errs.KindNotFound:
		return status.Error(codes.NotFound, domainError.Message)
	case errs.KindFailedPrecondition:
		return status.Error(codes.FailedPrecondition, domainError.Message)
	case errs.KindAlreadyExists:
		return status.Error(codes.AlreadyExists, domainError.Message)
	case errs.KindUnavailable:
		return status.Error(codes.Unavailable, domainError.Message)
	default:
		return status.Error(codes.Internal, "internal error")
	}
}

func (a *API) logError(operation string, err error) {
	domainError, ok := errs.As(err)
	if ok && domainError.Kind != errs.KindInternal {
		return
	}
	a.logger.Error(operation, "error", err)
}
