package customers_usecase_test

import (
	"context"
	"database/sql"
	"errors"
	customer_entity "payment-layer-card-api/entities/customers"
	customers_usecase "payment-layer-card-api/usecases/customers"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCustomerByID_Execute_ValidID(t *testing.T) {
	id := "12345abcde"
	wantErr := false

	mockCustomer := &customer_entity.Customer{ID: id}
	mockRepo := new(customer_entity.CustomerRepositoryMock)
	mockRepo.On("GetByID", id).Return(mockCustomer, nil)

	getCustomerByID := customers_usecase.NewGetCustomerByID(mockRepo)

	ctx := context.Background()
	customer, err := getCustomerByID.Execute(ctx, id)

	if wantErr {
		assert.NotNil(t, err)
		assert.Nil(t, customer)
	} else {
		assert.Nil(t, err)
		assert.NotNil(t, customer)
		assert.Equal(t, id, customer.ID)
	}
}

func TestGetCustomerByID_Execute_NoID(t *testing.T) {
	id := ""
	wantErr := true

	mockRepo := new(customer_entity.CustomerRepositoryMock)

	getCustomerByID := customers_usecase.NewGetCustomerByID(mockRepo)

	ctx := context.Background()
	customer, err := getCustomerByID.Execute(ctx, id)

	if wantErr {
		assert.NotNil(t, err)
		assert.Nil(t, customer)
	} else {
		assert.Nil(t, err)
		assert.NotNil(t, customer)
	}
}

func TestGetCustomerByID_Execute_RepositoryError(t *testing.T) {
	id := "12345abcde"
	wantErr := true
	mockRepoErr := errors.New("database error")

	mockRepo := new(customer_entity.CustomerRepositoryMock)
	mockRepo.On("GetByID", id).Return(nil, mockRepoErr)

	getCustomerByID := customers_usecase.NewGetCustomerByID(mockRepo)

	ctx := context.Background()
	customer, err := getCustomerByID.Execute(ctx, id)

	if wantErr {
		assert.NotNil(t, err)
		assert.Nil(t, customer)
	} else {
		assert.Nil(t, err)
		assert.NotNil(t, customer)
	}
}

func TestGetCustomerByID_Execute_NotFoundInRepository(t *testing.T) {
	id := "12345abcde"

	mockRepoErr := sql.ErrNoRows

	mockRepo := new(customer_entity.CustomerRepositoryMock)
	mockRepo.On("GetByID", id).Return(nil, mockRepoErr)

	getCustomerByID := customers_usecase.NewGetCustomerByID(mockRepo)

	ctx := context.Background()
	customer, err := getCustomerByID.Execute(ctx, id)

	assert.NotNil(t, err)
	assert.Nil(t, customer)
}
