package purchase

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

type Response struct {
}

type Handler struct {
	domain domain.PurchaseMaker
}

func New(domain domain.PurchaseMaker) *Handler {
	return &Handler{
		domain,
	}
}

func (h *Handler) Handle(req Request) (Response, error) {
	var res Response

	err := h.domain.Purchase(req.User)

	return res, err
}
