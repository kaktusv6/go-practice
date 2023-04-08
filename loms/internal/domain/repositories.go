package domain

import "context"

//go:generate bash -c "mkdir -p mocks"
//go:generate bash -c "rm -rf ./mocks/order_item_repository_minimock.go"
//go:generate bash -c "rm -rf ./mocks/order_item_stock_repository_minimock.go"
//go:generate bash -c "rm -rf ./mocks/order_repository_minimock.go"
//go:generate bash -c "rm -rf ./mocks/stock_repository_minimock.go"
//go:generate minimock -i StockRepository -o ./mocks/ -s "_minimock.go"
//go:generate minimock -i OrderRepository -o ./mocks/ -s "_minimock.go"
//go:generate minimock -i OrderItemRepository -o ./mocks/ -s "_minimock.go"
//go:generate minimock -i OrderItemStockRepository -o ./mocks/ -s "_minimock.go"

type StockRepository interface {
	GetListBySKU(ctx context.Context, sku uint32) ([]*Stock, error)
	UpdateCount(ctx context.Context, stock *Stock) error
	GetByWarehouseIDAndSku(ctx context.Context, warehouseID int64, sku uint32) (*Stock, error)
}

type OrderRepository interface {
	GetById(ctx context.Context, orderID int64) (*Order, error)
	Save(ctx context.Context, order *Order) error
	Update(ctx context.Context, order *Order) error
	GetAll(ctx context.Context) ([]*Order, error)
}

type OrderItemRepository interface {
	GetByOrderId(ctx context.Context, orderID int64) ([]*Item, error)
	SaveMany(ctx context.Context, orderID int64, items []*Item) error
}

type OrderItemStockRepository interface {
	Save(ctx context.Context, orderItemStock *OrderItemStock) error
	GetListByOrderID(ctx context.Context, orderID int64) ([]*OrderItemStock, error)
	Delete(ctx context.Context, orderItemStock *OrderItemStock) error
}
