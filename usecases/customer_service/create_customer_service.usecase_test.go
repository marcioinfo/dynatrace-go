package customerservice_usecase_test

import (
	"context"
	"testing"
	"time"

	customerservice "payment-layer-card-api/entities/customer_service"
	"payment-layer-card-api/mocks"
	customerservice_usecase "payment-layer-card-api/usecases/customer_service"

	"github.com/adhfoundation/layer-tools/datetypes"
	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateCustomerService_Execute_Success(t *testing.T) {
	mockRepo := mocks.NewCustomerServiceRepository(t)

	customerServiceInput := &customerservice.CustomerService{
		ServiceID:  "service-id-1",
		CustomerID: "customer-id-1",
		Name:       "John Doe",
		Document:   "96744054023",
		Email:      "john.doe@example.com",
		Phone:      "123456789",
		Gender:     "male",
		BirthDate:  datetypes.CustomDate(time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)),
	}

	mockRepo.On("Insert", mock.Anything, mock.AnythingOfType("*customerservice.CustomerService")).Return(nil).Run(func(args mock.Arguments) {
		customer := args.Get(1).(*customerservice.CustomerService)
		customer.ID = "new-id"
		customer.CreatedAt = datetypes.CustomDateTime(time.Now())
	})

	uc := customerservice_usecase.NewCreateCustomerService(mockRepo)

	ctx := context.Background()
	output, err := uc.Execute(ctx, customerServiceInput)

	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, "new-id", output.ID)
	mockRepo.AssertExpectations(t)
}

func TestCreateCustomerService_Execute_Failure_DocumentInvalid(t *testing.T) {
	mockRepo := mocks.NewCustomerServiceRepository(t)

	customerServiceInput := &customerservice.CustomerService{
		Document: "invalid-document",
	}

	uc := customerservice_usecase.NewCreateCustomerService(mockRepo)

	ctx := context.Background()
	output, err := uc.Execute(ctx, customerServiceInput)
	assert.Nil(t, output)
	assert.NotNil(t, err)
	assert.Equal(t, errors.DocumentIsInvalid, err.Code)
	mockRepo.AssertExpectations(t)
}

func TestCreateCustomerService_Execute_Failure_InsertError(t *testing.T) {
	mockRepo := mocks.NewCustomerServiceRepository(t)

	customerServiceInput := &customerservice.CustomerService{
		ServiceID:  "service-id-1",
		CustomerID: "customer-id-1",
		Name:       "John Doe",
		Document:   "96744054023",
		Email:      "john.doe@example.com",
		Phone:      "123456789",
		Gender:     "male",
		BirthDate:  datetypes.CustomDate(time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)),
	}

	mockRepo.On("Insert", mock.Anything, mock.AnythingOfType("*customerservice.CustomerService")).Return(errors.NewError(errors.InternalServerError, nil, "database error"))

	uc := customerservice_usecase.NewCreateCustomerService(mockRepo)

	ctx := context.Background()
	output, err := uc.Execute(ctx, customerServiceInput)

	assert.Nil(t, output)
	assert.NotNil(t, err)
	assert.Equal(t, errors.InternalServerError, err.Code)
	mockRepo.AssertExpectations(t)
}
