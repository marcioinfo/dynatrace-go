package customer_dto

import (
	"github.com/adhfoundation/layer-tools/datetypes"
)

type CustomerWithTokensOutputDTO struct {
	ID        string                   `json:"id" example:"835524dd-d897-45ed-a68e-2e1a8bf522d7"`
	Name      string                   `json:"name" example:"João Silva"`
	Email     string                   `json:"email" example:"teste@gmail.com"`
	Document  string                   `json:"document" example:"12345678901"`
	BirthDate datetypes.CustomDate     `json:"birth_date" example:"2000-02-11"`
	Phone     string                   `json:"phone" example:"5511999999999"`
	Gender    string                   `json:"gender" example:"male || female"`
	Tokens    []*CustomerTokens        `json:"tokens"`
	CreatedAt datetypes.CustomDateTime `json:"created_at" example:"2006-01-02T15:04:05-07:00"`
	UpdatedAt datetypes.CustomDateTime `json:"updated_at" example:"2006-01-02T15:04:05-07:00"`
}

type CustomerTokens struct {
	Token   string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"`
	Gateway string `json:"gateway" example:"pagarme"`
}

type CustomerOutputDTO struct {
	ID        string                   `json:"id" example:"835524dd-d897-45ed-a68e-2e1a8bf522d7"`
	Name      string                   `json:"name" example:"João Silva"`
	Email     string                   `json:"email" example:"teste@gmail.com"`
	Document  string                   `json:"document" example:"12345678901"`
	BirthDate datetypes.CustomDate     `json:"birth_date" example:"2000-02-30"`
	Phone     string                   `json:"phone" example:"5511999999999"`
	Gender    string                   `json:"gender" example:"male || female"`
	CreatedAt datetypes.CustomDateTime `json:"created_at" example:"2006-01-02T15:04:05-07:00"`
	UpdatedAt datetypes.CustomDateTime `json:"updated_at" example:"2006-01-02T15:04:05-07:00"`
}
