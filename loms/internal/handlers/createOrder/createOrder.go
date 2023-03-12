package createOrder

import (
	"github.com/pkg/errors"
	"route256/loms/internal/domain"
)

type Request struct {
	User  int64         `json:"user"`
	Items []domain.Item `json:"items"`
}

var (
	ReqUserErrorEmpty    = errors.New("empty user")
	ReqUserErrorNegative = errors.New("negative user")
	ReqItemsErrorEmpty   = errors.New("empty items")
)

func (r Request) Validate() error {
	switch {
	case r.User == 0:
		return ReqUserErrorEmpty
	case r.User < 0:
		return ReqUserErrorNegative
	case len(r.Items) == 0:
		return ReqItemsErrorEmpty
	}
	return nil
}

type ResOrderID struct {
	OrderID int64 `json:"orderID"`
}

type Handler struct {
	domain domain.OrderCreator
}

func New(domain domain.OrderCreator) *Handler {
	return &Handler{
		domain,
	}
}

func (h *Handler) Handle(request Request) (ResOrderID, error) {
	response := ResOrderID{}

	orderID, err := h.domain.CreateOrder(request.User, request.Items)
	if err != nil {
		return response, err
	}
	response.OrderID = orderID

	return response, nil
}
