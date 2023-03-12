package checkoutV1

import (
	"route256/checkout/internal/domain"
	desc "route256/checkout/pkg/checkout_v1"
)

type Implementation struct {
	desc.UnimplementedCheckoutV1Server
	domain domain.Domain
}

func New(domain domain.Domain) desc.CheckoutV1Server {
	return &Implementation{
		desc.UnimplementedCheckoutV1Server{},
		domain,
	}
}
