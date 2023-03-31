package db

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type key string

// TxKey ключ подключения к базе в контексте
const (
	TxKey key = "tx"
)

type Transactor interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

type Pinger interface {
	Ping(ctx context.Context) error
}

type QueryEngine interface {
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

type QueryEngineProvider interface {
	GetQueryEngine(ctx context.Context) QueryEngine // tx/pool
}

type DB interface {
	QueryEngine
	Transactor
	Pinger
	Close()
}

type db struct {
	pool *pgxpool.Pool
}

func (d *db) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return d.pool.BeginTx(ctx, txOptions)
}

func (d *db) Ping(ctx context.Context) error {
	return d.pool.Ping(ctx)
}

func (d *db) Close() {
	d.pool.Close()
}

func (d *db) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	return d.pool.Query(ctx, query, args)
}

func (d *db) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return d.pool.QueryRow(ctx, sql, args)
}

func GetContextTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, TxKey, tx)
}
