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

type StockRepository struct {
	provider db.QueryEngineProvider
}

func NewStockRepository(provider db.QueryEngineProvider) domain.StockRepository {
	return &StockRepository{
		provider,
	}
}

const (
	stockTable = "stocks"
)

func (s *StockRepository) GetListBySKU(ctx context.Context, sku uint32) ([]*domain.Stock, error) {
	sqQuery := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("sku", "warehouse_id", "count").
		From(stockTable).
		Where("sku = ?", sku)

	rawQuery, args, err := sqQuery.ToSql()
	if err != nil {
		return nil, err
	}

	db := s.provider.GetQueryEngine(ctx)

	var queryResult []Stock
	if err := pgxscan.Select(ctx, db, &queryResult, rawQuery, args...); err != nil {
		return nil, err
	}

	result := make([]*domain.Stock, 0, len(queryResult))
	for _, stock := range queryResult {
		result = append(result, &domain.Stock{
			Sku:         stock.Sku,
			WarehouseID: stock.WarehouseID,
			Count:       stock.Count,
		})
	}

	return result, nil
}

func (s *StockRepository) UpdateCount(ctx context.Context, stock *domain.Stock) error {
	sqQuery := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Update(stockTable).
		Set("count", stock.Count).
		Where("sku = ?", stock.Sku).
		Where("warehouse_id = ?", stock.WarehouseID)

	rawQuery, args, err := sqQuery.ToSql()
	if err != nil {
		return err
	}

	db := s.provider.GetQueryEngine(ctx)
	rows, err := db.Query(ctx, rawQuery, args...)
	defer rows.Close()

	if err != nil {
		return err
	}

	return nil
}

func (s *StockRepository) GetByWarehouseIDAndSku(ctx context.Context, warehouseID int64, sku uint32) (*domain.Stock, error) {
	sqQuery := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("sku", "warehouse_id", "count").
		From(stockTable).
		Where("warehouse_id = ?", warehouseID).
		Where("sku = ?", sku)

	rawQuery, args, err := sqQuery.ToSql()
	if err != nil {
		return nil, err
	}

	db := s.provider.GetQueryEngine(ctx)

	var queryResult Stock
	if err := pgxscan.Get(ctx, db, &queryResult, rawQuery, args...); err != nil {
		return nil, err
	}

	stock := &domain.Stock{
		Sku:         queryResult.Sku,
		WarehouseID: queryResult.WarehouseID,
		Count:       queryResult.Count,
	}

	return stock, nil
}
