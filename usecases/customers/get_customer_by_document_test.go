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

func TestGetCustomerByDocument_Execute_ValidDocument(t *testing.T) {
	document := "71973773000100"
	wantErr := false

	mockCustomer := &customer_entity.Customer{Document: document}
	mockRepo := new(customer_entity.CustomerRepositoryMock)
	mockRepo.On("GetCustomerByDocument", document).Return(mockCustomer, nil)

	getCustomerByDocument := customers_usecase.NewGetCustomerByDocument(mockRepo)

	ctx := context.Background()
	customer, err := getCustomerByDocument.Execute(ctx, document, "1")

	if wantErr {
		assert.NotNil(t, err)
		assert.Nil(t, customer)
	} else {
		assert.Nil(t, err)
		assert.NotNil(t, customer)
		assert.Equal(t, document, customer.Document)
	}
}

func TestGetCustomerByDocument_Execute_NoDocument(t *testing.T) {
	document := ""
	wantErr := true

	mockRepo := new(customer_entity.CustomerRepositoryMock)

	getCustomerByDocument := customers_usecase.NewGetCustomerByDocument(mockRepo)

	ctx := context.Background()
	customer, err := getCustomerByDocument.Execute(ctx, document, "1")

	if wantErr {
		assert.NotNil(t, err)
		assert.Nil(t, customer)
	} else {
		assert.Nil(t, err)
		assert.NotNil(t, customer)
	}
}

func TestGetCustomerByDocument_Execute_RepositoryError(t *testing.T) {
	document := "71973773000100"
	wantErr := true
	mockRepoErr := errors.New("database error")

	mockRepo := new(customer_entity.CustomerRepositoryMock)
	mockRepo.On("GetCustomerByDocument", document).Return(nil, mockRepoErr)

	getCustomerByDocument := customers_usecase.NewGetCustomerByDocument(mockRepo)

	ctx := context.Background()
	customer, err := getCustomerByDocument.Execute(ctx, document, "1")

	if wantErr {
		assert.NotNil(t, err)
		assert.Nil(t, customer)
	} else {
		assert.Nil(t, err)
		assert.NotNil(t, customer)
	}
}

func TestGetCustomerByDocument_Execute_NotFoundInRepository(t *testing.T) {
	document := "71973773000100"
	mockRepoErr := sql.ErrNoRows

	mockRepo := new(customer_entity.CustomerRepositoryMock)
	mockRepo.On("GetCustomerByDocument", document).Return(nil, mockRepoErr)

	getCustomerByDocument := customers_usecase.NewGetCustomerByDocument(mockRepo)

	ctx := context.Background()
	customer, err := getCustomerByDocument.Execute(ctx, document, "1")

	assert.NotNil(t, err)
	assert.Nil(t, customer)
}
