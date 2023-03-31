package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

// Client ...
type Client interface {
	Close() error
	DB() DB
}

type client struct {
	db DB
}

func NewClient(ctx context.Context, config *Config) (Client, error) {
	clientImpl := &client{}
	err := clientImpl.bootstrap(ctx, config)
	if err != nil {
		return nil, err
	}
	return clientImpl, nil
}

func (c *client) Close() error {
	if c != nil {
		if c.db != nil {
			c.db.Close()
		}
	}

	return nil
}

func (c *client) DB() DB {
	return c.db
}

func (c *client) bootstrap(ctx context.Context, config *Config) error {
	psqlConn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.Name,
	)

	pool, err := pgxpool.Connect(ctx, psqlConn)
	if err != nil {
		return err
	}

	if err := pool.Ping(ctx); err != nil {
		return err
	}

	return nil
}
