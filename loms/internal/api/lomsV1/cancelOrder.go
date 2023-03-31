package lomsV1

import (
	"context"
)

import (
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/protobuf/types/known/emptypb"
	desc "route256/loms/pkg/loms_v1"
)

func (i *Implementation) CancelOrder(ctx context.Context, req *desc.OrderID) (*emptypb.Empty, error) {
	order, err := i.domain.GetListOrder(ctx, req.OrderId)
	if err != nil {
		return nil, err
	}

	err = i.domain.CancelOrder(ctx, order)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}
