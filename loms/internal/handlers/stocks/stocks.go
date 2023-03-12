package stocks

import (
	"github.com/pkg/errors"
	"route256/loms/internal/domain"
)

type Request struct {
	SKU uint32 `json:"sku"`
}

var (
	ReqSKUErrorEmpty = errors.New("empty SKU")
)

func (r Request) Validate() error {
	if r.SKU == 0 {
		return ReqSKUErrorEmpty
	}
	return nil
}

type Response struct {
	Stocks []domain.Stock `json:"stocks"`
}

type Handler struct {
	domain domain.StocksGetter
}

func New(domain domain.StocksGetter) *Handler {
	return &Handler{domain}
}

func (h *Handler) Handle(req Request) (Response, error) {
	response := Response{}
	stocks, err := h.domain.GetStocksBySKU(req.SKU)
	if err != nil {
		return response, err
	}
	response.Stocks = stocks
	return response, nil
}
