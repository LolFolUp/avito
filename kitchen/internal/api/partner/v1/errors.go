package v1

import (
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/errs"
)

func grpcError(err error) error {
	domainErr, ok := errs.As(err)
	if !ok {
		return status.Error(codes.Internal, "internal server error")
	}

	code := codes.Internal
	switch domainErr.Kind {
	case errs.KindInvalidArgument:
		code = codes.InvalidArgument
	case errs.KindNotFound:
		code = codes.NotFound
	case errs.KindFailedPrecondition:
		code = codes.FailedPrecondition
	case errs.KindAlreadyExists:
		code = codes.AlreadyExists
	}

	message := domainErr.Message
	if code == codes.Internal {
		message = "internal server error"
	}
	return status.Error(code, message)
}

func (a *API) logError(operation string, err error) {
	domainErr, ok := errs.As(err)
	if ok && domainErr.Kind != errs.KindInternal {
		return
	}
	a.logger.Error(operation, slog.Any("error", err))
}
