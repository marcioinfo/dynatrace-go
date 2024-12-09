package card_token

import (
	"github.com/adhfoundation/layer-tools/datetypes"
	"github.com/google/uuid"
)

type CardToken struct {
	ID        string
	CardID    string
	Gateway   string
	CardToken string
	CreatedAt datetypes.CustomDateTime
	UpdatedAt datetypes.CustomDateTime
}

var CardTokensAttributes = map[string]string{
	"TableName":  "Tokens de Cartão",
	"id":         "Id",
	"card_id":    "Id do Cartão",
	"gateway":    "Gateway do token",
	"card_token": "Token do cartão",
	"created_at": "Data de criação",
	"updated_at": "Data da última alteração",
}

func NewCardToken() *CardToken {
	return &CardToken{}
}

func (c *CardToken) InitID() {
	c.ID = uuid.New().String()
}
