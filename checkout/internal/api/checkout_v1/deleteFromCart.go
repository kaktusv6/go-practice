package checkoutV1

import (
	"context"
)

import (
	"google.golang.org/protobuf/types/known/emptypb"
	"route256/checkout/internal/domain"
	desc "route256/checkout/pkg/checkout_v1"
)

func (i *Implementation) DeleteFromCart(ctx context.Context, req *desc.ItemInfoRequest) (*emptypb.Empty, error) {
	err := i.domain.DeleteFromCart(ctx, domain.ItemInfo{
		User:  req.GetUser(),
		Sku:   req.GetSku(),
		Count: uint16(req.GetCount()),
	})
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
