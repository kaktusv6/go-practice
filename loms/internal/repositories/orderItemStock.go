package repositories

import (
	"context"
)

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"route256/libs/db"
	"route256/loms/internal/domain"
)

type OrderItemStockRepository struct {
	provider db.QueryEngineProvider
}

func NewOrderItemStockRepository(provider db.QueryEngineProvider) domain.OrderItemStockRepository {
	return &OrderItemStockRepository{
		provider,
	}
}

const orderItemStockTable = "orders_items_stocks"

func (o *OrderItemStockRepository) Save(ctx context.Context, orderItemStock *domain.OrderItemStock) error {
	sqQuery := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert(orderItemStockTable).
		Columns("order_id", "sku", "count", "warehouse_id").
		Values(
			orderItemStock.OrderId,
			orderItemStock.Sku,
			orderItemStock.Count,
			orderItemStock.WarehouseID,
		)

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

func (o *OrderItemStockRepository) GetListByOrderID(ctx context.Context, orderID int64) ([]*domain.OrderItemStock, error) {
	sqQuery := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("order_id", "sku", "count", "warehouse_id").
		From(orderItemStockTable).
		Where("order_id = ?", orderID)

	rawQuery, args, err := sqQuery.ToSql()
	if err != nil {
		return nil, err
	}

	db := o.provider.GetQueryEngine(ctx)

	var queryResult []OrderItemStock
	if err := pgxscan.Select(ctx, db, &queryResult, rawQuery, args...); err != nil {
		return nil, err
	}

	result := make([]*domain.OrderItemStock, 0, len(queryResult))
	for _, orderItemStock := range queryResult {
		result = append(result, &domain.OrderItemStock{
			OrderId:     orderItemStock.OrderId,
			Sku:         orderItemStock.Sku,
			WarehouseID: orderItemStock.WarehouseID,
			Count:       orderItemStock.Count,
		})
	}

	return result, nil
}

func (o *OrderItemStockRepository) Delete(ctx context.Context, orderItemStock *domain.OrderItemStock) error {
	sqQuery := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Delete(orderItemStockTable).
		Where("order_id = ?", orderItemStock.OrderId).
		Where("sku = ?", orderItemStock.Sku).
		Where("warehouse_id = ?", orderItemStock.WarehouseID)

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
