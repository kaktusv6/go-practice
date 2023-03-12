package cancelOrder

import (
	"github.com/pkg/errors"
	"route256/loms/internal/domain"
)

type Request struct {
	OrderID int64 `json:"orderID"`
}

var (
	ReqOrderIDErrorEmpty    = errors.New("empty OrderID")
	ReqOrderIDErrorNegative = errors.New("negative OrderID")
)

func (r Request) Validate() error {
	switch {
	case r.OrderID == 0:
		return ReqOrderIDErrorEmpty
	case r.OrderID < 0:
		return ReqOrderIDErrorNegative
	}
	return nil
}

type Response struct{}

type Handler struct {
	domain domain.OrderCanceling
}

func New(domain domain.OrderCanceling) *Handler {
	return &Handler{
		domain,
	}
}

func (h *Handler) Handle(req Request) (Response, error) {
	err := h.domain.CancelOrder(req.OrderID)
	return Response{}, err
}
