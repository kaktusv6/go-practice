package listCart

import (
	"github.com/pkg/errors"
	"route256/checkout/internal/domain"
)

type Request struct {
	User int64 `json:"user"`
}

var (
	ReqUserIDErrorEmpty    = errors.New("empty UserID")
	ReqUserIDErrorNegative = errors.New("negative UserID")
)

func (r Request) Validate() error {
	switch {
	case r.User == 0:
		return ReqUserIDErrorEmpty
	case r.User < 0:
		return ReqUserIDErrorNegative
	}

	return nil
}

type Handler struct {
	domain domain.CartListItemGetter
}

func New(domain domain.CartListItemGetter) *Handler {
	return &Handler{
		domain,
	}
}

func (h *Handler) Handle(req Request) (domain.Cart, error) {
	return h.domain.GetListItems(req.User)
}
