package payment_gateway

import (
	"context"
	"payment-layer-card-api/entities/customers"
	card_token_dto "payment-layer-card-api/usecases/card_token/dtos"
	customer_token_dto "payment-layer-card-api/usecases/customer_token/dtos"
	customer_dto "payment-layer-card-api/usecases/customers/dtos"

	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
	"github.com/stretchr/testify/mock"
)

type PaymentGatewayMock struct {
	mock.Mock
}

func (m *PaymentGatewayMock) CreateCustomer(ctx context.Context, input *customer_dto.CreateCustomerDTOInput, id string) (customerTokenOutputDTO *customer_token_dto.CustomerTokenOutputDTO, httpStatusCode int, err *errors.ErrorOutput) {
	args := m.Called(input, id)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2).(*errors.ErrorOutput)
	}
	return args.Get(0).(*customer_token_dto.CustomerTokenOutputDTO), args.Int(1), args.Error(2).(*errors.ErrorOutput)
}

func (m *PaymentGatewayMock) CreateCard(ctx context.Context, input *card_token_dto.CardTokenInputDTO) (*card_token_dto.CardTokenOutputDTO, int, *errors.ErrorOutput) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2).(*errors.ErrorOutput)
	}
	return args.Get(0).(*card_token_dto.CardTokenOutputDTO), args.Int(1), args.Error(2).(*errors.ErrorOutput)
}

func (m *PaymentGatewayMock) DeleteCard(ctx context.Context, customerId string, cardId string) (*card_token_dto.CardTokenOutputDTO, int, *errors.ErrorOutput) {
	args := m.Called(customerId, cardId)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2).(*errors.ErrorOutput)
	}
	return args.Get(0).(*card_token_dto.CardTokenOutputDTO), args.Int(1), args.Error(2).(*errors.ErrorOutput)
}

func (m *PaymentGatewayMock) GatewayName() string {
	args := m.Called()
	return args.String(0)
}

func (m *PaymentGatewayMock) UpdateCustomer(ctx context.Context, customer *customers.Customer, id string) (CustomerTokenOutputDTO *customer_token_dto.CustomerTokenOutputDTO, httpStatusCode int, err *errors.ErrorOutput) {
	args := m.Called(customer, id)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Get(2).(*errors.ErrorOutput)
	}
	return args.Get(0).(*customer_token_dto.CustomerTokenOutputDTO), args.Int(1), args.Get(2).(*errors.ErrorOutput)
}
