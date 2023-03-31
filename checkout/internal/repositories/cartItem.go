package repositories

import (
	"context"
)

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"route256/checkout/internal/domain"
	dbLib "route256/libs/db"
)

type CartItemRepository struct {
	provider dbLib.QueryEngineProvider
}

func NewOrderItemRepository(provider dbLib.QueryEngineProvider) domain.CartItemRepository {
	return &CartItemRepository{
		provider: provider,
	}
}

const (
	cartItemTable = "cart_items"
)

func (c *CartItemRepository) Create(ctx context.Context, cartItem *domain.CartItem) error {
	sqQuery := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert(cartItemTable).
		Columns("user_id", "sku", "count").
		Values(cartItem.User, cartItem.Sku, cartItem.Count)

	rawQuery, args, err := sqQuery.ToSql()
	if err != nil {
		return err
	}

	db := c.provider.GetQueryEngine(ctx)
	rows, err := db.Query(ctx, rawQuery, args...)
	defer rows.Close()

	if err != nil {
		return err
	}

	return nil
}

func (c *CartItemRepository) GetUserCartItems(ctx context.Context, user int64) ([]*domain.CartItem, error) {
	sqQuery := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("user_id", "sku", "count").
		From(cartItemTable).
		Where("user_id = ?", user)

	rawQuery, args, err := sqQuery.ToSql()
	if err != nil {
		return nil, err
	}

	db := c.provider.GetQueryEngine(ctx)

	var queryCartItems []CartItem
	if err := pgxscan.Select(ctx, db, &queryCartItems, rawQuery, args...); err != nil {
		return nil, err
	}

	cartItems := make([]*domain.CartItem, 0, len(queryCartItems))
	for _, cartItem := range queryCartItems {
		cartItems = append(cartItems, &domain.CartItem{
			User:  cartItem.User,
			Sku:   cartItem.Sku,
			Count: cartItem.Count,
		})
	}

	return cartItems, nil
}

func (c *CartItemRepository) GetOne(ctx context.Context, user int64, sku uint32) (*domain.CartItem, error) {
	sqQuery := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("user_id", "sku", "count").
		From(cartItemTable).
		Where("user_id = ?", user).
		Where("sku = ?", sku)

	rawQuery, args, err := sqQuery.ToSql()
	if err != nil {
		return nil, err
	}

	db := c.provider.GetQueryEngine(ctx)

	var queryCartItem CartItem
	if err := pgxscan.Get(ctx, db, &queryCartItem, rawQuery, args...); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = nil
		}
		return nil, err
	}

	cartItem := &domain.CartItem{
		User:  queryCartItem.User,
		Sku:   queryCartItem.Sku,
		Count: queryCartItem.Count,
	}

	return cartItem, nil
}

func (c *CartItemRepository) Update(ctx context.Context, cartItem *domain.CartItem) error {
	sqQuery := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Update(cartItemTable).
		Set("count", cartItem.Count).
		Where("user_id = ?", cartItem.User).
		Where("sku = ?", cartItem.Sku)

	rawQuery, args, err := sqQuery.ToSql()
	if err != nil {
		return err
	}

	db := c.provider.GetQueryEngine(ctx)
	rows, err := db.Query(ctx, rawQuery, args...)
	defer rows.Close()

	if err != nil {
		return err
	}

	return nil
}

func (c *CartItemRepository) Delete(ctx context.Context, cartItem *domain.CartItem) error {
	sqQuery := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Delete(cartItemTable).
		Where("user_id = ?", cartItem.User).
		Where("sku = ?", cartItem.Sku)

	rawQuery, args, err := sqQuery.ToSql()
	if err != nil {
		return err
	}

	db := c.provider.GetQueryEngine(ctx)
	rows, err := db.Query(ctx, rawQuery, args...)
	defer rows.Close()

	if err != nil {
		return err
	}

	return nil
}
