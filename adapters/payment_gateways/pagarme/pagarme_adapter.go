package pagarme

import (
	"payment-layer-card-api/common/types"
)

type PagarmeAdapter struct {
	url           string
	secretKey     string
	authorization string
}

func NewPagarmeAdapter(url string, secretKey string, authorization string) *PagarmeAdapter {
	return &PagarmeAdapter{
		url:           url,
		secretKey:     secretKey,
		authorization: authorization,
	}
}

func (p *PagarmeAdapter) GatewayName() string {
	return types.PAGARME.ToString()
}
