package db

import (
	"context"

	generated "github.com/AmithSAI007/prj-flexdesk-user-service/internal/db/generated"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	generated.Querier
	ExecTx(ctx context.Context, fn func(generated.Querier) error) error
}

type SQLStore struct {
	connPool *pgxpool.Pool
	*generated.Queries
}

func NewStore(connPool *pgxpool.Pool) Store {
	return &SQLStore{
		connPool: connPool,
		Queries:  generated.New(connPool),
	}
}

var _ Store = (*SQLStore)(nil)

func (store *SQLStore) ExecTx(ctx context.Context, fn func(generated.Querier) error) error {
	tx, err := store.connPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	q := generated.New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return rbErr
		}
		return err
	}

	return tx.Commit(ctx)
}
