package payment_gateway

import (
	"context"
	"payment-layer-card-api/entities/customers"
	card_token_dto "payment-layer-card-api/usecases/card_token/dtos"
	customer_token_dto "payment-layer-card-api/usecases/customer_token/dtos"
	customer_dto "payment-layer-card-api/usecases/customers/dtos"

	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
)

type PaymentGatewayInterface interface {
	CreateCustomer(ctx context.Context, customers *customer_dto.CreateCustomerDTOInput, id string) (customerTokenOutputDTO *customer_token_dto.CustomerTokenOutputDTO, httpStatusCode int, err *errors.ErrorOutput)
	CreateCard(ctx context.Context, input *card_token_dto.CardTokenInputDTO) (cardTokenOutputDTO *card_token_dto.CardTokenOutputDTO, httpStatusCode int, err *errors.ErrorOutput)
	DeleteCard(ctx context.Context, customerId string, cardId string) (cardTokenOutputDTO *card_token_dto.CardTokenOutputDTO, httpStatusCode int, err *errors.ErrorOutput)
	UpdateCustomer(ctx context.Context, customer *customers.Customer, id string) (CustomerTokenOutputDTO *customer_token_dto.CustomerTokenOutputDTO, httpStatusCode int, err *errors.ErrorOutput)
	GatewayName() string
}
