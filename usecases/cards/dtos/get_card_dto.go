package dtos

import "github.com/adhfoundation/layer-tools/datetypes"

type GetCardWithTokensDTO struct {
	ID          string                   `json:"id" example:"835524dd-d897-45ed-a68e-2e1a8bf522d7"`          // Identificador do cartão
	CustomerID  string                   `json:"customer_id" example:"835524dd-d897-45ed-a68e-2e1a8bf52314"` // Identificador do cliente
	Holder      string                   `json:"holder" example:"Joao Silva"`                                // Nome do portador do cartão
	Brand       string                   `json:"brand" example:"visa"`                                       // Bandeira do cartão
	FirstDigits string                   `json:"first_digits" example:"411111"`                              // Primeiros dígitos do cartão
	LastDigits  string                   `json:"last_digits" example:"1111"`                                 // Últimos dígitos do cartão
	ExpMonth    string                   `json:"exp_month" example:"12"`                                     // Mês de expiração do cartão
	ExpYear     string                   `json:"exp_year" example:"2025"`                                    // Ano de expiração do cartão
	Tokens      []*CardToken             `json:"tokens" `                                                    // Tokens do cartão
	CreatedAt   datetypes.CustomDateTime `json:"created_at" example:"2006-01-02T15:04:05-07:00"`             // Data de criação do cartão
	UpdatedAt   datetypes.CustomDateTime `json:"updated_at" example:"2006-01-02T15:04:05-07:00"`             // Data de atualização do cartão
}

type CardToken struct {
	Token   string `json:"token" example:"9KwUaB1kwz8="` // Token do cartão
	Gateway string `json:"gateway" example:"rede"`       // Gateway do token
}

type GetCardWithoutTokensDTO struct {
	ID          string                   `json:"id" example:"835524dd-d897-45ed-a68e-2e1a8bf522d7"`          // Identificador do cartão
	CustomerID  string                   `json:"customer_id" example:"835524dd-d897-45ed-a68e-2e1a8bf52314"` // Identificador do cliente
	Holder      string                   `json:"holder" example:"Joao Silva"`                                // Nome do portador do cartão
	Brand       string                   `json:"brand" example:"visa"`                                       // Bandeira do cartão
	FirstDigits string                   `json:"first_digits" example:"411111"`                              // Primeiros dígitos do cartão
	LastDigits  string                   `json:"last_digits" example:"1111"`                                 // Últimos dígitos do cartão
	ExpMonth    string                   `json:"exp_month" example:"12"`                                     // Mês de expiração do cartão
	ExpYear     string                   `json:"exp_year" example:"2025"`                                    // Ano de expiração do cartão
	CreatedAt   datetypes.CustomDateTime `json:"created_at" example:"2006-01-02T15:04:05-07:00"`             // Data de criação do cartão
	UpdatedAt   datetypes.CustomDateTime `json:"updated_at" example:"2006-01-02T15:04:05-07:00"`             // Data de atualização do cartão
}
