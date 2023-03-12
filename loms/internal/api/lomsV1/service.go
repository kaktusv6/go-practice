package lomsV1

import (
	"route256/loms/internal/domain"
	desc "route256/loms/pkg/loms_v1"
)

type Implementation struct {
	desc.UnimplementedLomsV1Server

	domain domain.Domain
}

func NewLomsV1(domain domain.Domain) desc.LomsV1Server {
	return &Implementation{
		desc.UnimplementedLomsV1Server{},
		domain,
	}
}
