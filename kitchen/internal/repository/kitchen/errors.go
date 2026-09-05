package kitchen

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-flow-2-lolfolup-0eb52e46/kitchen/internal/errs"
)

func repositoryError(operation, resource string, err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return errs.NotFound(resource)
	}

	var postgresErr *pgconn.PgError
	if errors.As(err, &postgresErr) {
		switch postgresErr.Code {
		case "23505":
			return errs.AlreadyExists(resource)
		case "23503", "23514", "22P02":
			return errs.InvalidArgument("invalid " + resource)
		}
	}

	return errs.Internal(operation, err)
}
