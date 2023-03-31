package repositories

import (
	"context"
	"time"
)

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"route256/libs/db"
	"route256/loms/internal/domain"
)

type OrderRepository struct {
	provider db.QueryEngineProvider
}

func NewOrderRepository(provider db.QueryEngineProvider) domain.OrderRepository {
	return &OrderRepository{
		provider,
	}
}

const (
	orderTable = "orders"
)

func (o *OrderRepository) GetById(ctx context.Context, orderID int64) (*domain.Order, error) {
	sqQuery := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("id", "user_id", "status", "created_at", "updated_at").
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

	return o.bindingTo(&order), nil
}

func (o *OrderRepository) Save(ctx context.Context, order *domain.Order) error {
	now := time.Now()

	sqQuery := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert(orderTable).
		Columns("user_id", "status", "created_at", "updated_at").
		Values(order.User, order.Status, now, now).
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
		Set("updated_at", time.Now()).
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

func (o *OrderRepository) GetAll(ctx context.Context) ([]*domain.Order, error) {
	sqQuery := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("id", "user_id", "status", "created_at", "updated_at").
		From(orderTable)

	rawQuery, args, err := sqQuery.ToSql()
	if err != nil {
		return nil, err
	}

	db := o.provider.GetQueryEngine(ctx)

	var orders []Order
	if err := pgxscan.Select(ctx, db, &orders, rawQuery, args...); err != nil {
		return nil, err
	}

	result := make([]*domain.Order, 0, len(orders))
	for _, order := range orders {
		result = append(result, o.bindingTo(&order))
	}

	return result, nil
}

func (o *OrderRepository) bindingTo(order *Order) *domain.Order {
	return &domain.Order{
		ID:        order.ID,
		Status:    order.Status,
		User:      order.User,
		CreatedAt: order.CreatedAt.Time,
		UpdatedAt: order.UpdatedAt.Time,
	}
}
