package restaurant

import (
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func New(pgx *pgxpool.Pool, getter *trmpgx.CtxGetter) *Repository {
	if getter == nil {
		getter = trmpgx.DefaultCtxGetter
	}
	return &Repository{
		pool:   pgx,
		getter: getter,
	}
}
