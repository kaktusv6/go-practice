package orderPayed

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

type Response struct{}

type Handler struct {
	domain domain.OrderPayedMarker
}

func New(domain domain.OrderPayedMarker) *Handler {
	return &Handler{
		domain,
	}
}

func (h *Handler) Handle(req Request) (Response, error) {
	err := h.domain.OrderPayedMark(req.OrderID)
	return Response{}, err
}
