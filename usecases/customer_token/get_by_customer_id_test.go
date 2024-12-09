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
)

func TestGetByCustomerID_Execute_ValidCustomerID(t *testing.T) {
	customerID := "123456"

	mockToken := &customer_token_entity.CustomerToken{
		ID:            "1",
		CustomerId:    customerID,
		CustomerToken: "sample-token",
		Gateway:       "sample-gateway",
		CreatedAt:     datetypes.CustomDateTime(time.Now()),
		UpdatedAt:     datetypes.CustomDateTime(time.Now()),
	}

	mockRepo := new(customer_token_entity.CustomerTokenRepositoryMock)
	mockRepo.On("GetByCustomerID", customerID).Return([]*customer_token_entity.CustomerToken{mockToken}, nil)

	getByCustomerID := customer_token_usecase.NewGetByCustomerID(mockRepo)

	ctx := context.Background()
	tokens, err := getByCustomerID.Execute(ctx, customerID)

	assert.Nil(t, err)
	assert.NotNil(t, tokens)
	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, customerID, tokens[0].CustomerId)
}

func TestGetByCustomerID_Execute_NoCustomerID(t *testing.T) {
	customerID := ""

	mockRepo := new(customer_token_entity.CustomerTokenRepositoryMock)

	getByCustomerID := customer_token_usecase.NewGetByCustomerID(mockRepo)

	ctx := context.Background()
	tokens, err := getByCustomerID.Execute(ctx, customerID)

	assert.NotNil(t, err)
	assert.Nil(t, tokens)
}

func TestGetByCustomerID_Execute_RepositoryError(t *testing.T) {
	customerID := "123456"

	mockRepoErr := errors.New("database error")

	mockRepo := new(customer_token_entity.CustomerTokenRepositoryMock)
	mockRepo.On("GetByCustomerID", customerID).Return(nil, mockRepoErr)

	getByCustomerID := customer_token_usecase.NewGetByCustomerID(mockRepo)

	ctx := context.Background()
	tokens, err := getByCustomerID.Execute(ctx, customerID)

	assert.NotNil(t, err)
	assert.Nil(t, tokens)
}

func TestGetByCustomerID_Execute_NotFoundInRepository(t *testing.T) {
	customerID := "123456"

	mockRepoErr := sql.ErrNoRows

	mockRepo := new(customer_token_entity.CustomerTokenRepositoryMock)
	mockRepo.On("GetByCustomerID", customerID).Return(nil, mockRepoErr)

	getByCustomerID := customer_token_usecase.NewGetByCustomerID(mockRepo)

	ctx := context.Background()
	tokens, err := getByCustomerID.Execute(ctx, customerID)

	assert.NotNil(t, err)
	assert.Nil(t, tokens)
	assert.Equal(t, layerErrors.NotFoundError, err.Type())
}
