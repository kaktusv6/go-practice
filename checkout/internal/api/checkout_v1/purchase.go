package checkoutV1

import (
	"context"
)

import (
	"google.golang.org/protobuf/types/known/emptypb"
	desc "route256/checkout/pkg/checkout_v1"
)

func (i *Implementation) Purchase(ctx context.Context, userInfo *desc.UserInfo) (*emptypb.Empty, error) {
	err := i.domain.Purchase(ctx, userInfo.GetUser())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
