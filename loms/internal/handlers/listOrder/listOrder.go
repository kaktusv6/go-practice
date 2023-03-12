package listOrder

import (
	"github.com/pkg/errors"
	"route256/loms/internal/domain"
)

type Request struct {
	OrderID int64 `json:"orderID"`
}

var (
	ReqOrderIDErrorEmpty    = errors.New("empty orderID")
	ReqOrderIDErrorNegative = errors.New("negative orderID")
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

type Handler struct {
	domain domain.ListOrderGetter
}

func New(domain domain.ListOrderGetter) *Handler {
	return &Handler{
		domain,
	}
}

func (h *Handler) Handle(req Request) (domain.Order, error) {
	// Fixture
	return h.domain.GetListOrder(req.OrderID)
}
