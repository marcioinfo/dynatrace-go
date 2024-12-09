package customers_usecase_test

import (
	"context"
	"errors"
	customer_entity "payment-layer-card-api/entities/customers"
	customers_usecase "payment-layer-card-api/usecases/customers"
	customer_dto "payment-layer-card-api/usecases/customers/dtos"
	"testing"

	"github.com/adhfoundation/layer-tools/datetypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateCustomer_Execute_ValidInput(t *testing.T) {
	birthdate, _ := datetypes.NewDateFromString("1990-09-14")
	input := customer_dto.CreateCustomerDTOInput{
		FirstName: "John",
		LastName:  "Doe",
		Document:  "71973773000100",
		BirthDate: birthdate,
		Gender:    "M",
		Phone:     "123456789",
		Email:     "johndoe@example.com",
	}
	mockRepoErr := error(nil)

	mockRepo := new(customer_entity.CustomerRepositoryMock)
	if mockRepoErr != nil {
		mockRepo.On("GetCustomerByEmail", input.Email).Return(nil, nil)
		mockRepo.On("Insert", mock.AnythingOfType("*customers.Customer")).Return(mockRepoErr)
	} else {
		mockRepo.On("GetCustomerByEmail", input.Email).Return(nil, nil)
		mockRepo.On("Insert", mock.AnythingOfType("*customers.Customer")).Return(nil)
	}

	createCustomer := customers_usecase.NewCreateCustomer(mockRepo)

	ctx := context.Background()
	output, err := createCustomer.Execute(ctx, input)

	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.NotEmpty(t, output.ID)

}

func TestCreateCustomer_Execute_InvalidDocument(t *testing.T) {
	birthdate, _ := datetypes.NewDateFromString("1990-09-14")

	input := customer_dto.CreateCustomerDTOInput{
		FirstName: "John",
		LastName:  "Doe",
		Document:  "invalidDocument",
		BirthDate: birthdate,
		Gender:    "M",
		Phone:     "123456789",
		Email:     "johndoe@example.com",
	}
	mockRepoErr := error(nil)

	mockRepo := new(customer_entity.CustomerRepositoryMock)
	if mockRepoErr != nil {
		mockRepo.On("GetCustomerByEmail", mock.AnythingOfType("*customers.Customer")).Return(nil, nil)
		mockRepo.On("Insert", mock.AnythingOfType("*customers.Customer")).Return(mockRepoErr)
	} else {
		mockRepo.On("GetCustomerByEmail", mock.AnythingOfType("*customers.Customer")).Return(nil, nil)
		mockRepo.On("Insert", mock.AnythingOfType("*customers.Customer")).Return(nil)
	}

	createCustomer := customers_usecase.NewCreateCustomer(mockRepo)

	ctx := context.Background()
	output, err := createCustomer.Execute(ctx, input)

	assert.NotNil(t, err)
	assert.Nil(t, output)
}

func TestCreateCustomer_Execute_InvalidBirthDate(t *testing.T) {
	birthdate, _ := datetypes.NewDateFromString("19980-09-14")

	input := customer_dto.CreateCustomerDTOInput{
		FirstName: "John",
		LastName:  "Doe",
		Document:  "71973773000100",
		BirthDate: birthdate,
		Gender:    "M",
		Phone:     "123456789",
		Email:     "johndoe@example.com",
	}
	mockRepoErr := error(nil)

	mockRepo := new(customer_entity.CustomerRepositoryMock)
	if mockRepoErr != nil {
		mockRepo.On("GetCustomerByEmail", input.Email).Return(nil, nil)
		mockRepo.On("Insert", mock.AnythingOfType("*customers.Customer")).Return(mockRepoErr)
	} else {
		mockRepo.On("GetCustomerByEmail", input.Email).Return(nil, nil)
		mockRepo.On("Insert", mock.AnythingOfType("*customers.Customer")).Return(nil)
	}

	createCustomer := customers_usecase.NewCreateCustomer(mockRepo)

	ctx := context.Background()
	output, err := createCustomer.Execute(ctx, input)

	assert.NotNil(t, err)
	assert.Nil(t, output)
}

func TestCreateCustomer_Execute_InvalidEmail(t *testing.T) {
	birthdate, _ := datetypes.NewDateFromString("1990-09-14")

	input := customer_dto.CreateCustomerDTOInput{
		FirstName: "John",
		LastName:  "Doe",
		Document:  "71973773000100",
		BirthDate: birthdate,
		Gender:    "M",
		Phone:     "123456789",
		Email:     "",
	}
	mockRepoErr := error(nil)

	mockRepo := new(customer_entity.CustomerRepositoryMock)
	if mockRepoErr != nil {
		mockRepo.On("GetCustomerByEmail", input.Email).Return(nil, nil)
		mockRepo.On("Insert", mock.AnythingOfType("*customers.Customer")).Return(mockRepoErr)
	} else {
		mockRepo.On("GetCustomerByEmail", input.Email).Return(nil, nil)
		mockRepo.On("Insert", mock.AnythingOfType("*customers.Customer")).Return(nil)
	}

	createCustomer := customers_usecase.NewCreateCustomer(mockRepo)

	ctx := context.Background()
	output, err := createCustomer.Execute(ctx, input)

	assert.NotNil(t, err)
	assert.Nil(t, output)
}

func TestCreateCustomer_Execute_RepositoryError(t *testing.T) {
	birthdate, _ := datetypes.NewDateFromString("1990-09-14")

	input := customer_dto.CreateCustomerDTOInput{
		FirstName: "John",
		LastName:  "Doe",
		Document:  "71973773000100",
		BirthDate: birthdate,
		Gender:    "M",
		Phone:     "123456789",
		Email:     "johndoe@example.com",
	}
	mockRepoErr := errors.New("database error")

	mockRepo := new(customer_entity.CustomerRepositoryMock)
	if mockRepoErr != nil {
		mockRepo.On("GetCustomerByEmail", input.Email).Return(nil, nil)
		mockRepo.On("Insert", mock.AnythingOfType("*customers.Customer")).Return(mockRepoErr)
	} else {
		mockRepo.On("GetCustomerByEmail", input.Email).Return(nil, nil)
		mockRepo.On("Insert", mock.AnythingOfType("*customers.Customer")).Return(nil)
	}

	createCustomer := customers_usecase.NewCreateCustomer(mockRepo)

	ctx := context.Background()
	output, err := createCustomer.Execute(ctx, input)

	assert.NotNil(t, err)
	assert.Nil(t, output)
}
