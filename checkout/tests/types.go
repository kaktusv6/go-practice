package tests

import (
	"github.com/gojuno/minimock/v3"
	"route256/checkout/internal/domain"
	"route256/libs/db"
)

type stockRepositoryMock func(mc *minimock.Controller) domain.StockRepository

type orderRepositoryMock func(mc *minimock.Controller) domain.OrderRepository

type cartItemRepositoryMock func(mc *minimock.Controller) domain.CartItemRepository

type productRepositoryMock func(mc *minimock.Controller) domain.ProductRepository

type dbMock func(mc *minimock.Controller) db.DB
