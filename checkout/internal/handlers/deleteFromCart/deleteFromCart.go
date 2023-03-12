package deleteFromCart

import (
	"github.com/pkg/errors"
	"route256/checkout/internal/domain"
)

type Request struct {
	User  int64  `json:"user"`
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
}

var (
	ReqUserIDErrorEmpty    = errors.New("empty UserID")
	ReqUserIDErrorNegative = errors.New("negative UserID")
	ReqSKUErrorEmpty       = errors.New("empty SKU")
	ReqCountErrorEmpty     = errors.New("empty count")
)

func (r Request) Validate() error {
	switch {
	case r.User == 0:
		return ReqUserIDErrorEmpty
	case r.User < 0:
		return ReqUserIDErrorNegative
	case r.SKU == 0:
		return ReqSKUErrorEmpty
	case r.Count == 0:
		return ReqCountErrorEmpty
	}

	return nil
}

type Response struct {
}

type Handler struct {
	domain domain.CartItemDeleting
}

func New(domain domain.CartItemDeleting) *Handler {
	return &Handler{
		domain,
	}
}

func (h *Handler) Handle(req Request) (Response, error) {
	var res Response

	err := h.domain.DeleteFromCart(req.User, req.SKU, req.Count)
	if err != nil {
		return res, err
	}

	return res, nil
}
