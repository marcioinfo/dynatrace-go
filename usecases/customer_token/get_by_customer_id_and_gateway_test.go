package customer_token_test

import (
	"context"
	"database/sql"
	"errors"
	customer_token_entity "payment-layer-card-api/entities/customer_token"
	customer_token_usecase "payment-layer-card-api/usecases/customer_token"
	"testing"
	"time"

	"github.com/adhfoundation/layer-tools/datetypes"
	layerErrors "github.com/adhfoundation/payment-layer-error-package/pkg/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetByCustomerIDAndGateway_Execute_ErrorOnEmptyCustomerId(t *testing.T) {

	customerTokenRepositoryMock := new(customer_token_entity.CustomerTokenRepositoryMock)

	getByCustomerIDAndGatewayUseCase := customer_token_usecase.NewGetByCustomerIDAndGateway(customerTokenRepositoryMock)
	ctx := context.Background()
	customerTokens, err := getByCustomerIDAndGatewayUseCase.Execute(ctx, "", "")

	assert.NotNil(t, err)
	assert.Equal(t, layerErrors.ParameterIsRequired, err.Type())
	assert.Nil(t, customerTokens)

}

func TestGetByCustomerIDAndGateway_Execute_ErrorOnEmptyGateway(t *testing.T) {

	customerTokenRepositoryMock := new(customer_token_entity.CustomerTokenRepositoryMock)

	getByCustomerIDAndGatewayUseCase := customer_token_usecase.NewGetByCustomerIDAndGateway(customerTokenRepositoryMock)
	ctx := context.Background()
	customerTokens, err := getByCustomerIDAndGatewayUseCase.Execute(ctx, "CustomerID", "")

	assert.NotNil(t, err)
	assert.Equal(t, layerErrors.WithoutGateway, err.Type())
	assert.Nil(t, customerTokens)
}

func TestGetByCustomerIDAndGateway_Execute_ErrorOnRepositoryDatabaseError(t *testing.T) {

	customerTokenRepositoryMock := new(customer_token_entity.CustomerTokenRepositoryMock)
	mockError := errors.New("error on get information from database")
	customerTokenRepositoryMock.On("GetByCustomerIDAndGateway", mock.Anything, mock.Anything).Return(nil, mockError)

	getByCustomerIDAndGatewayUseCase := customer_token_usecase.NewGetByCustomerIDAndGateway(customerTokenRepositoryMock)
	ctx := context.Background()
	customerTokens, err := getByCustomerIDAndGatewayUseCase.Execute(ctx, "CustomerID", "rede")

	assert.NotNil(t, err)
	assert.Nil(t, customerTokens)
	assert.Equal(t, mockError.Error(), err.Error())

}

func TestGetByCustomerIDAndGateway_Execute_ErrorOnRepositoryNotFoundError(t *testing.T) {

	customerTokenRepositoryMock := new(customer_token_entity.CustomerTokenRepositoryMock)
	mockError := sql.ErrNoRows
	customerTokenRepositoryMock.On("GetByCustomerIDAndGateway", mock.Anything, mock.Anything).Return(nil, mockError)

	getByCustomerIDAndGatewayUseCase := customer_token_usecase.NewGetByCustomerIDAndGateway(customerTokenRepositoryMock)
	ctx := context.Background()
	customerTokens, err := getByCustomerIDAndGatewayUseCase.Execute(ctx, "CustomerID", "rede")

	assert.NotNil(t, err)
	assert.Nil(t, customerTokens)
	assert.Equal(t, "sql: no rows in result set", err.Error())

}

func TestGetByCustomerIDAndGateway_Execute_Sucess(t *testing.T) {

	customerTokenRepositoryMock := new(customer_token_entity.CustomerTokenRepositoryMock)

	customerTokenMock := customer_token_entity.CustomerToken{
		ID:            "ID",
		CustomerId:    "CustomerId",
		CustomerToken: "CustomerToken",
		Gateway:       "rede",
		CreatedAt:     datetypes.CustomDateTime(time.Now()),
		UpdatedAt:     datetypes.CustomDateTime(time.Now()),
	}

	customerTokenRepositoryMock.On("GetByCustomerIDAndGateway", mock.Anything, mock.Anything).Return(&customerTokenMock, nil)

	getByCustomerIDAndGatewayUseCase := customer_token_usecase.NewGetByCustomerIDAndGateway(customerTokenRepositoryMock)
	ctx := context.Background()
	customerToken, err := getByCustomerIDAndGatewayUseCase.Execute(ctx, "CustomerID", "rede")

	assert.Nil(t, err)
	assert.NotNil(t, customerToken)
	assert.Equal(t, customerToken, &customerTokenMock)
}
