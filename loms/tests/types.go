package tests

import (
	"github.com/gojuno/minimock/v3"
	"route256/libs/db"
	"route256/loms/internal/domain"
)

type stockRepositoryMock func(mc *minimock.Controller) domain.StockRepository

type orderRepositoryMock func(mc *minimock.Controller) domain.OrderRepository

type orderItemRepositoryMock func(mc *minimock.Controller) domain.OrderItemRepository

type orderItemStockRepositoryMock func(mc *minimock.Controller) domain.OrderItemStockRepository

type dbMock func(mc *minimock.Controller) db.DB

type orderStatusNotifierMock func(mc *minimock.Controller) domain.OrderStatusNotifier
