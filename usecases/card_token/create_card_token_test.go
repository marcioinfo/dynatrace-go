package card_token_test

import (
	"context"
	card_token_entity "payment-layer-card-api/entities/card_token"
	payment_gateway_entity "payment-layer-card-api/entities/payment_gateway"
	card_token_usecase "payment-layer-card-api/usecases/card_token"
	card_token_dto "payment-layer-card-api/usecases/card_token/dtos"
	"testing"

	layerErrors "github.com/adhfoundation/payment-layer-error-package/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var cardTokenInputDTOValidSample = &card_token_dto.CardTokenInputDTO{
	CardID:        "sampleCardId",
	CustomerID:    "sampleCustomerId",
	CustomerToken: "sampleCustomerToken",
	CustomerEmail: "sampleCustomerEmail",
	CustomerAddress: &card_token_dto.HolderAddress{
		Line1:      "Rua",
		Line2:      "Bairro",
		State:      "Estado",
		City:       "Cidade",
		Country:    "BR",
		District:   "Distrito",
		PostalCode: "12345678",
	},
	CustomerPhone:  "123456789",
	Number:         "1234123412341234",
	Holder:         "John Doe",
	HolderDocument: "12313412",
	ExpMonth:       "09",
	ExpYear:        "2025",
	CVV:            "123",
	Brand:          "Visa",
}

func TestCardTokenUseCase_Execute_Error400InGateway(t *testing.T) {
	input := cardTokenInputDTOValidSample

	cardTokenRepositoryMock := new(card_token_entity.CardTokenRepositoryMock)

	paymentGatewayMock := new(payment_gateway_entity.PaymentGatewayMock)
	paymentGatewayMock.On("GatewayName").Return("pagarme")
	paymentGatewayMock.On("CreateCard", input).Return(nil, 400, layerErrors.NewGatewayError(layerErrors.InvalidCard, "gateway error"))

	cardRetryCount := 1
	cardRetryDelayInMilliseconds := 500

	usecase := card_token_usecase.NewCreateCardToken(
		cardTokenRepositoryMock,
		paymentGatewayMock,
		cardRetryCount,
		cardRetryDelayInMilliseconds,
	)

	ctx := context.Background()
	cardTokenDTOOutput, err := usecase.Execute(ctx, input.CardID, input)

	assert.Nil(t, cardTokenDTOOutput)
	assert.NotNil(t, err)
	assert.Equal(t, layerErrors.InvalidCard, err.Type())
}

func TestCardTokenUseCase_Execute_Error500InGateway(t *testing.T) {
	input := cardTokenInputDTOValidSample

	cardTokenRepositoryMock := new(card_token_entity.CardTokenRepositoryMock)

	paymentGatewayMock := new(payment_gateway_entity.PaymentGatewayMock)
	paymentGatewayMock.On("GatewayName").Return("pagarme")
	paymentGatewayMock.On("CreateCard", input).Return(nil, 500, layerErrors.NewGatewayError(layerErrors.InternalServerError, "gateway error"))

	cardRetryCount := 1
	cardRetryDelayInMilliseconds := 500

	usecase := card_token_usecase.NewCreateCardToken(
		cardTokenRepositoryMock,
		paymentGatewayMock,
		cardRetryCount,
		cardRetryDelayInMilliseconds,
	)

	ctx := context.Background()
	cardTokenDTOOutput, err := usecase.Execute(ctx, input.CardID, input)

	assert.Nil(t, cardTokenDTOOutput)
	assert.NotNil(t, err)
	assert.Equal(t, layerErrors.InternalServerError, err.Code)
}

func TestCardTokenUseCase_Execute_ErrorOnInsertCardInDatabase(t *testing.T) {
	input := cardTokenInputDTOValidSample
	gatewayName := "pagarme"

	cardTokenRepositoryMock := new(card_token_entity.CardTokenRepositoryMock)
	cardTokenRepositoryMock.On("Insert", mock.Anything).Return(layerErrors.NewError(layerErrors.InternalServerError, nil, "database error"))

	cardTokenOutputDTO := &card_token_dto.CardTokenOutputDTO{
		Gateway: gatewayName,
		Token:   "sampleToken",
	}
	paymentGatewayMock := new(payment_gateway_entity.PaymentGatewayMock)
	paymentGatewayMock.On("GatewayName").Return(gatewayName)
	paymentGatewayMock.On("CreateCard", input).Return(cardTokenOutputDTO, 200, (*layerErrors.ErrorOutput)(nil))

	cardRetryCount := 1
	cardRetryDelayInMilliseconds := 500

	usecase := card_token_usecase.NewCreateCardToken(
		cardTokenRepositoryMock,
		paymentGatewayMock,
		cardRetryCount,
		cardRetryDelayInMilliseconds,
	)

	ctx := context.Background()
	cardTokenDTOOutput, err := usecase.Execute(ctx, input.CardID, input)

	assert.Nil(t, cardTokenDTOOutput)
	assert.NotNil(t, err)
	assert.Equal(t, layerErrors.InternalServerError, err.Code)
}

func TestCardTokenUseCase_Execute_Success(t *testing.T) {
	input := cardTokenInputDTOValidSample
	gatewayName := "pagarme"

	cardTokenRepositoryMock := new(card_token_entity.CardTokenRepositoryMock)
	cardTokenRepositoryMock.On("Insert", mock.Anything).Return(nil)

	cardTokenOutputDTO := &card_token_dto.CardTokenOutputDTO{
		Gateway: gatewayName,
		Token:   "sampleToken",
	}
	paymentGatewayMock := new(payment_gateway_entity.PaymentGatewayMock)
	paymentGatewayMock.On("GatewayName").Return(gatewayName)
	paymentGatewayMock.On("CreateCard", input).Return(cardTokenOutputDTO, 200, (*layerErrors.ErrorOutput)(nil))

	expected := cardTokenOutputDTO

	cardRetryCount := 1
	cardRetryDelayInMilliseconds := 500

	usecase := card_token_usecase.NewCreateCardToken(
		cardTokenRepositoryMock,
		paymentGatewayMock,
		cardRetryCount,
		cardRetryDelayInMilliseconds,
	)

	ctx := context.Background()
	resultCardTokenOutputDTO, err := usecase.Execute(ctx, input.CardID, input)

	assert.NotNil(t, resultCardTokenOutputDTO)
	assert.Nil(t, err)
	assert.Equal(t, expected, resultCardTokenOutputDTO)
}

func TestCardTokenUseCase_Execute_SuccessWithRetry(t *testing.T) {
	input := cardTokenInputDTOValidSample
	gatewayName := "pagarme"

	cardTokenRepositoryMock := new(card_token_entity.CardTokenRepositoryMock)
	cardTokenRepositoryMock.On("Insert", mock.Anything).Return(nil)

	cardTokenOutputDTO := &card_token_dto.CardTokenOutputDTO{
		Gateway: gatewayName,
		Token:   "sampleToken",
	}
	paymentGatewayMock := new(payment_gateway_entity.PaymentGatewayMock)
	paymentGatewayMock.On("GatewayName").Return(gatewayName)
	paymentGatewayMock.On("CreateCard", input).Return((*card_token_dto.CardTokenOutputDTO)(nil), 500, layerErrors.NewGatewayError(layerErrors.ExternalError, "gateway error")).Once()
	paymentGatewayMock.On("CreateCard", input).Return(cardTokenOutputDTO, 200, (*layerErrors.ErrorOutput)(nil)).Once()

	expected := cardTokenOutputDTO

	cardRetryCount := 2
	cardRetryDelayInMilliseconds := 500

	usecase := card_token_usecase.NewCreateCardToken(
		cardTokenRepositoryMock,
		paymentGatewayMock,
		cardRetryCount,
		cardRetryDelayInMilliseconds,
	)

	ctx := context.Background()
	resultCardTokenOutputDTO, err := usecase.Execute(ctx, input.CardID, input)

	assert.NotNil(t, resultCardTokenOutputDTO)
	assert.Nil(t, err)
	assert.Equal(t, expected, resultCardTokenOutputDTO)
}
