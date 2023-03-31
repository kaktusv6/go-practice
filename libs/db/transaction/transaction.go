package transaction

import (
	"context"
	"route256/libs/db"

	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

type Handler func(ctx context.Context) error

type Manager interface {
	db.QueryEngineProvider
	RepeatableRead(ctx context.Context, f Handler) error
}

type manager struct {
	db db.DB
}

func NewTransactionManager(db db.DB) Manager {
	return &manager{
		db: db,
	}
}

func (m *manager) GetQueryEngine(ctx context.Context) db.QueryEngine {
	return m.db
}

func (m *manager) transaction(ctx context.Context, opts pgx.TxOptions, fn Handler) (err error) {
	tx, ok := ctx.Value(db.TxKey).(pgx.Tx)
	if ok {
		return fn(ctx)
	}

	tx, err = m.db.BeginTx(ctx, opts)
	if err != nil {
		return errors.Wrap(err, "can't begin transaction")
	}

	ctx = db.GetContextTx(ctx, tx)

	defer func() {
		// recover from panic
		if r := recover(); r != nil {
			err = errors.Errorf("panic recovered: %v", r)
		}

		if nil == err {
			err = tx.Commit(ctx)
			if err != nil {
				err = errors.Wrap(err, "tx commit failed")
			}
		}

		if err != nil {
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				err = errors.Wrapf(err, "errRollback: %v", errRollback)
			}
		}
	}()

	if err = fn(ctx); err != nil {
		err = errors.Wrap(err, "failed executing code inside transaction")
	}

	return err
}

func (m *manager) RepeatableRead(ctx context.Context, f Handler) error {
	opts := pgx.TxOptions{IsoLevel: pgx.RepeatableRead}
	return m.transaction(ctx, opts, f)
}
