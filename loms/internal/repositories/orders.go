package repositories

import (
	"context"
)

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"route256/libs/transactor"
	"route256/loms/internal/domain"
)

type OrderRepository struct {
	provider transactor.QueryEngineProvider
}

func NewOrderRepository(provider transactor.QueryEngineProvider) domain.OrderRepository {
	return &OrderRepository{
		provider,
	}
}

const (
	orderTable = "orders"
)

func (o *OrderRepository) GetById(ctx context.Context, orderID int64) (*domain.Order, error) {
	sqQuery := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("id", "user_id", "status").
		From(orderTable).
		Where("id = ?", orderID)

	rawQuery, args, err := sqQuery.ToSql()
	if err != nil {
		return nil, err
	}

	db := o.provider.GetQueryEngine(ctx)

	var order Order
	if err := pgxscan.Get(ctx, db, &order, rawQuery, args...); err != nil {
		return nil, err
	}

	return &domain.Order{
		ID:     order.ID,
		Status: order.Status,
		User:   order.User,
	}, nil
}

func (o *OrderRepository) Save(ctx context.Context, order *domain.Order) error {
	sqQuery := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).Insert(orderTable).
		Columns("user_id", "status").
		Values(order.User, order.Status).
		Suffix("returning id")

	rawQuery, args, err := sqQuery.ToSql()
	if err != nil {
		return err
	}

	db := o.provider.GetQueryEngine(ctx)
	row := db.QueryRow(ctx, rawQuery, args...)

	err = row.Scan(&order.ID)
	if err != nil {
		return err
	}

	return nil
}

func (o *OrderRepository) Update(ctx context.Context, order *domain.Order) error {
	sqQuery := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Update(orderTable).
		Set("status", order.Status).
		Set("user_id", order.User).
		Where("id = ?", order.ID)

	rawQuery, args, err := sqQuery.ToSql()
	if err != nil {
		return err
	}

	db := o.provider.GetQueryEngine(ctx)
	rows, err := db.Query(ctx, rawQuery, args...)
	defer rows.Close()

	if err != nil {
		return err
	}

	return nil
}
