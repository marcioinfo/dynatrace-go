package customer_token_test

import (
	"context"
	customer_token_entity "payment-layer-card-api/entities/customer_token"
	payment_gateway_entity "payment-layer-card-api/entities/payment_gateway"
	customer_token_usecase "payment-layer-card-api/usecases/customer_token"
	customer_token_dto "payment-layer-card-api/usecases/customer_token/dtos"
	customer_dto "payment-layer-card-api/usecases/customers/dtos"
	"testing"

	"github.com/adhfoundation/layer-tools/datetypes"
	layerErrors "github.com/adhfoundation/payment-layer-error-package/pkg/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateCustomerToken_Execute_ValidInput(t *testing.T) {
	birthdate, _ := datetypes.NewDateFromString("14/09/1990")

	input := &customer_dto.CreateCustomerDTOInput{
		FirstName: "John",
		LastName:  "Doe",
		Document:  "71973773000100",
		BirthDate: birthdate,
		Gender:    "M",
		Phone:     "123456789",
		Email:     "johndoe@example.com",
	}
	customerID := "randomID"

	mockCustomerTokenRepository := new(customer_token_entity.CustomerTokenRepositoryMock)
	paymentGatewayMock := new(payment_gateway_entity.PaymentGatewayMock)

	paymentGatewayMock.On("CreateCustomer", mock.Anything, mock.Anything).Return(&customer_token_dto.CustomerTokenOutputDTO{
		Gateway: "someGateway",
		Token:   "someToken",
	}, 200, (*layerErrors.ErrorOutput)(nil))
	mockCustomerTokenRepository.On("GetByCustomerIDAndGateway", customerID, "someGateway").Return(nil, layerErrors.NewError(layerErrors.NotFoundError, nil, "Not Found"))
	paymentGatewayMock.On("GatewayName").Return("someGateway")
	mockCustomerTokenRepository.On("Insert", mock.AnythingOfType("*customer_token.CustomerToken")).Return(nil)
	mockCustomerTokenRepository.On("GetByCustomerIDAndGateway", mock.Anything, mock.Anything, mock.Anything).Return(nil, layerErrors.NewError(layerErrors.NotFoundError, nil, "Not Found"))
	customerRetryCount := 1
	customerRetryDelayInMilliseconds := 500
	createCustomerToken := customer_token_usecase.NewCreateCustomerToken(mockCustomerTokenRepository, paymentGatewayMock, customerRetryCount, customerRetryDelayInMilliseconds)

	ctx := context.Background()
	output, err := createCustomerToken.Execute(ctx, input, customerID)

	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.Token)
	assert.NotEmpty(t, output.Gateway)
}

func TestCreateCustomerToken_Execute_PaymentGatewayErrorWithRetry500(t *testing.T) {
	birthdate, _ := datetypes.NewDateFromString("14/09/1990")

	input := &customer_dto.CreateCustomerDTOInput{
		FirstName: "John",
		LastName:  "Doe",
		Document:  "71973773000100",
		BirthDate: birthdate,
		Gender:    "M",
		Phone:     "123456789",
		Email:     "johndoe@example.com",
	}
	customerID := "randomID"

	mockCustomerTokenRepository := new(customer_token_entity.CustomerTokenRepositoryMock)
	mockCustomerTokenRepository.On("Insert", mock.AnythingOfType("*customer_token.CustomerToken")).Return(nil)
	mockCustomerTokenRepository.On("GetByCustomerIDAndGateway", customerID, "someGateway").Return(nil, layerErrors.NewError(layerErrors.NotFoundError, nil, "Not Found"))
	paymentGatewayMock := new(payment_gateway_entity.PaymentGatewayMock)

	paymentGatewayMock.On("CreateCustomer", mock.Anything, mock.Anything).Return(nil, 500, layerErrors.NewError(layerErrors.InternalServerError, nil, "gateway error"))

	paymentGatewayMock.On("GatewayName").Return("someGateway")

	customerRetryCount := 3
	customerRetryDelayInMilliseconds := 500
	createCustomerToken := customer_token_usecase.NewCreateCustomerToken(mockCustomerTokenRepository, paymentGatewayMock, customerRetryCount, customerRetryDelayInMilliseconds)

	ctx := context.Background()
	output, err := createCustomerToken.Execute(ctx, input, customerID)

	assert.NotNil(t, err)
	assert.Equal(t, layerErrors.InvalidCustomer, layerErrors.PaymentLayerErrorsCode("GT-108"))
	assert.Nil(t, output)
}

func TestCreateCustomerToken_Execute_PaymentGatewayErrorWithRetry400(t *testing.T) {
	birthdate, _ := datetypes.NewDateFromString("14/09/1990")

	input := &customer_dto.CreateCustomerDTOInput{
		FirstName: "John",
		LastName:  "Doe",
		Document:  "71973773000100",
		BirthDate: birthdate,
		Gender:    "M",
		Phone:     "123456789",
		Email:     "johndoe@example.com",
	}
	customerID := "randomID"

	mockCustomerTokenRepository := new(customer_token_entity.CustomerTokenRepositoryMock)
	mockCustomerTokenRepository.On("Insert", mock.AnythingOfType("*customer_token.CustomerToken")).Return(nil)
	mockCustomerTokenRepository.On("GetByCustomerIDAndGateway", customerID, "someGateway").Return(nil, layerErrors.NewError(layerErrors.NotFoundError, nil, "Not Found"))
	paymentGatewayMock := new(payment_gateway_entity.PaymentGatewayMock)

	paymentGatewayMock.On("CreateCustomer", mock.Anything, mock.Anything).Return(nil, 500, layerErrors.NewError(layerErrors.InternalServerError, nil, "gateway error")).Once()
	paymentGatewayMock.On("CreateCustomer", mock.Anything, mock.Anything).Return(nil, 400, layerErrors.NewError(layerErrors.InternalServerError, nil, "gateway error")).Once()

	paymentGatewayMock.On("GatewayName").Return("someGateway")

	customerRetryCount := 3
	customerRetryDelayInMilliseconds := 500
	createCustomerToken := customer_token_usecase.NewCreateCustomerToken(mockCustomerTokenRepository, paymentGatewayMock, customerRetryCount, customerRetryDelayInMilliseconds)

	ctx := context.Background()
	output, err := createCustomerToken.Execute(ctx, input, customerID)

	assert.NotNil(t, err, err)
	assert.Nil(t, output)
	assert.Equal(t, layerErrors.InvalidCustomer, err.Type())
}

func TestCreateCustomerToken_Execute_PaymentGatewayErrorWithRetry500AndSuccess(t *testing.T) {
	birthdate, _ := datetypes.NewDateFromString("14/09/1990")

	input := &customer_dto.CreateCustomerDTOInput{
		FirstName: "John",
		LastName:  "Doe",
		Document:  "71973773000100",
		BirthDate: birthdate,
		Gender:    "M",
		Phone:     "123456789",
		Email:     "johndoe@example.com",
	}
	customerID := "randomID"

	mockCustomerTokenRepository := new(customer_token_entity.CustomerTokenRepositoryMock)
	mockCustomerTokenRepository.On("Insert", mock.AnythingOfType("*customer_token.CustomerToken")).Return(nil)
	mockCustomerTokenRepository.On("GetByCustomerIDAndGateway", customerID, "someGateway").Return(nil, layerErrors.NewError(layerErrors.NotFoundError, nil, "Not Found"))
	paymentGatewayMock := new(payment_gateway_entity.PaymentGatewayMock)

	paymentGatewayMock.On("CreateCustomer", mock.Anything, mock.Anything).Return(nil, 500, layerErrors.NewError(layerErrors.InternalServerError, nil, "gateway error")).Once()
	gateway := "someGateway"
	token := "someToken"
	paymentGatewayMock.On("CreateCustomer", mock.Anything, mock.Anything).Return(&customer_token_dto.CustomerTokenOutputDTO{
		Gateway: gateway,
		Token:   token,
	}, 200, (*layerErrors.ErrorOutput)(nil)).Once()
	paymentGatewayMock.On("GatewayName").Return("someGateway")

	customerRetryCount := 3
	customerRetryDelayInMilliseconds := 500
	createCustomerToken := customer_token_usecase.NewCreateCustomerToken(mockCustomerTokenRepository, paymentGatewayMock, customerRetryCount, customerRetryDelayInMilliseconds)

	ctx := context.Background()
	output, err := createCustomerToken.Execute(ctx, input, customerID)

	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, token, output.Token)
	assert.Equal(t, gateway, output.Gateway)
}

func TestCreateCustomerToken_Execute_RepositoryError(t *testing.T) {
	birthdate, _ := datetypes.NewDateFromString("14/09/1990")

	input := &customer_dto.CreateCustomerDTOInput{
		FirstName: "John",
		LastName:  "Doe",
		Document:  "71973773000100",
		BirthDate: birthdate,
		Gender:    "M",
		Phone:     "123456789",
		Email:     "johndoe@example.com",
	}
	customerID := "randomID"

	mockCustomerTokenRepository := new(customer_token_entity.CustomerTokenRepositoryMock)
	paymentGatewayMock := new(payment_gateway_entity.PaymentGatewayMock)

	paymentGatewayMock.On("CreateCustomer", mock.Anything, mock.Anything).Return(&customer_token_dto.CustomerTokenOutputDTO{
		Gateway: "someGateway",
		Token:   "someToken",
	}, 200, (*layerErrors.ErrorOutput)(nil))
	mockCustomerTokenRepository.On("GetByCustomerIDAndGateway", customerID, "someGateway").Return(nil, layerErrors.NewError(layerErrors.NotFoundError, nil, "Not Found"))
	paymentGatewayMock.On("GatewayName").Return("someGateway")
	mockCustomerTokenRepository.On("Insert", mock.AnythingOfType("*customer_token.CustomerToken")).Return(layerErrors.NewError(layerErrors.InternalServerError, nil, "repository error"))

	customerRetryCount := 1
	customerRetryDelayInMilliseconds := 500
	createCustomerToken := customer_token_usecase.NewCreateCustomerToken(mockCustomerTokenRepository, paymentGatewayMock, customerRetryCount, customerRetryDelayInMilliseconds)

	ctx := context.Background()
	output, err := createCustomerToken.Execute(ctx, input, customerID)

	assert.NotNil(t, err)
	assert.Nil(t, output)
}
