package customer_dto

import (
	"github.com/adhfoundation/layer-tools/datetypes"
	"github.com/adhfoundation/layer-tools/middlewares"
)

type SendCustomerDTO struct {
	ID                string                  `json:"id"`
	Name              string                  `json:"name"`
	Document          string                  `json:"document"`
	Birthdate         datetypes.CustomDate    `json:"birthdate"`
	Email             string                  `json:"email"`
	Phone             string                  `json:"phone"`
	Gender            string                  `json:"gender"`
	ServiceID         string                  `json:"service_id"`
	CustomerServiceID string                  `json:"customer_service_id"`
	Tokens            []Token                 `json:"tokens"`
	ApmLink           middlewares.ApmInfoSend `json:"apmLink"`
}

type Token struct {
	Token   string `json:"token"`
	Gateway string `json:"gateway"`
}
