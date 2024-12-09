package rede

import (
	"payment-layer-card-api/common/types"
)

type RedeAdapter struct {
	url          string
	merchant_key string
	merchant_id  string
}

func NewRedeAdapter(url string, merchant_key string, merchant_id string) *RedeAdapter {
	return &RedeAdapter{
		url,
		merchant_key,
		merchant_id,
	}
}

func (r *RedeAdapter) GatewayName() string {
	return types.REDE.ToString()
}
