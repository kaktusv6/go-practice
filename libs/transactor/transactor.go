package transactor

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/multierr"
)

type QueryEngine interface {
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

type QueryEngineProvider interface {
	GetQueryEngine(ctx context.Context) QueryEngine // tx/pool
}

type TransactionManager struct {
	pool *pgxpool.Pool
}

func NewQueryEngineProvider(pool *pgxpool.Pool) *TransactionManager {
	return &TransactionManager{
		pool: pool,
	}
}

const ContextKey = "transaction-pool"

func (p *TransactionManager) GetQueryEngine(ctx context.Context) QueryEngine {
	tx, ok := ctx.Value(ContextKey).(*pgxpool.Tx)
	if ok && tx != nil {
		return tx
	}

	return p.pool
}

func (p *TransactionManager) RunRepeatableReade(ctx context.Context, fx func(ctxTX context.Context) error) error {
	tx, err := p.GetQueryEngine(ctx).(*pgxpool.Pool).BeginTx(ctx,
		pgx.TxOptions{
			IsoLevel: pgx.RepeatableRead,
		})

	if err != nil {
		return err
	}

	if err := fx(context.WithValue(ctx, ContextKey, tx)); err != nil {
		return multierr.Combine(err, tx.Rollback(ctx))
	}

	if err := tx.Commit(ctx); err != nil {
		return multierr.Combine(err, tx.Rollback(ctx))
	}

	return nil
}
