package v1

import (
	"errors"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/errs"
)

func TestGRPCErrorMapping(t *testing.T) {
	tests := []struct {
		name string
		err  error
		code codes.Code
	}{
		{"invalid argument", errs.InvalidArgument("invalid"), codes.InvalidArgument},
		{"not found", errs.NotFound("order"), codes.NotFound},
		{"failed precondition", errs.FailedPrecondition(errs.CodeMenuChanged, "changed", nil), codes.FailedPrecondition},
		{"already exists", errs.AlreadyExists("order"), codes.AlreadyExists},
		{"domain internal", errs.Internal("query", errors.New("sql details")), codes.Internal},
		{"unknown internal", errors.New("unexpected"), codes.Internal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := status.Code(grpcError(tt.err)); got != tt.code {
				t.Fatalf("code = %v, want %v", got, tt.code)
			}
		})
	}
}

func TestGRPCErrorDoesNotExposeInternalCause(t *testing.T) {
	err := grpcError(errs.Internal("query order", errors.New("password=secret")))
	if got := status.Convert(err).Message(); got != "internal server error" {
		t.Fatalf("internal message leaked: %q", got)
	}
}
