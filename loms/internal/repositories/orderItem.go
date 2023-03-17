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

type OrderItemRepository struct {
	provider transactor.QueryEngineProvider
}

func NewOrderItemRepository(provider transactor.QueryEngineProvider) domain.OrderItemRepository {
	return &OrderItemRepository{
		provider,
	}
}

const (
	orderItemTable = "orders_items"
)

func (o *OrderItemRepository) GetByOrderId(ctx context.Context, orderID int64) ([]domain.Item, error) {
	sqQuery := sq.Select("sku", "count").
		From(orderItemTable).
		Where("order_id = $1", orderID)

	rawQuery, args, err := sqQuery.ToSql()
	if err != nil {
		return nil, err
	}

	db := o.provider.GetQueryEngine(ctx)

	var orderItems []OrderItem
	err = pgxscan.Select(ctx, db, &orderItems, rawQuery, args...)
	if err != nil {
		return nil, err
	}

	result := make([]domain.Item, 0, len(orderItems))
	for _, orderItem := range orderItems {
		result = append(result, domain.Item{
			Sku:   orderItem.Sku,
			Count: orderItem.Count,
		})
	}
	return result, nil
}

func (o *OrderItemRepository) SaveMany(ctx context.Context, orderID int64, items []domain.Item) error {
	sqQuery := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert(orderItemTable).
		Columns("order_id", "sku", "count")

	for _, item := range items {
		sqQuery = sqQuery.Values(orderID, item.Sku, item.Count)
	}

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
