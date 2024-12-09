package types

type GatewayEnum string

const (
	REDE    GatewayEnum = "rede"
	PAGARME GatewayEnum = "pagarme"
)

func (ge GatewayEnum) ToString() string {
	return string(ge)
}
