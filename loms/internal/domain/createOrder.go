package domain

import "context"

func (d *domain) CreateOrder(ctx context.Context, user int64, items []Item) (int64, error) {
	// Fixture
	return 123, nil
}
