package customer_token

import (
	"github.com/adhfoundation/layer-tools/datetypes"
	"github.com/google/uuid"
)

type CustomerToken struct {
	ID            string
	CustomerId    string
	CustomerToken string
	Gateway       string
	CreatedAt     datetypes.CustomDateTime
	UpdatedAt     datetypes.CustomDateTime
}

var CustomerTokensAttributes = map[string]string{
	"TableName":      "Tokens de Cliente",
	"id":             "Id",
	"customer_id":    "Id do Cliente",
	"gateway":        "Gateway do token",
	"customer_token": "Token do Cliente",
	"created_at":     "Data de criação",
	"updated_at":     "Data da última alteração",
}

func NewCustomerToken() *CustomerToken {
	return &CustomerToken{}
}

func (c *CustomerToken) InitID() {
	c.ID = uuid.New().String()
}
