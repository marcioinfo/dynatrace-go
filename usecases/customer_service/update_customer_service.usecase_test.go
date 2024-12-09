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

func TestUpdateCustomerService_Execute_Success(t *testing.T) {
	mockRepo := mocks.NewCustomerServiceRepository(t)

	customerServiceInput := &customerservice.CustomerService{
		ID:         "existing-id",
		ServiceID:  "updated-service-id",
		CustomerID: "updated-customer-id",
		Name:       "Jane Doe",
		Document:   "98765432109",
		Email:      "jane.doe@example.com",
		Phone:      "987654321",
		Gender:     "female",
		BirthDate:  datetypes.CustomDate(time.Date(1992, 2, 2, 0, 0, 0, 0, time.UTC)),
	}

	mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(cs *customerservice.CustomerService) bool {
		return cs.ID == "existing-id" && cs.CustomerID == "updated-customer-id"
	})).Return(nil)

	uc := customerservice_usecase.NewUpdateCustomerService(mockRepo)

	ctx := context.Background()
	output, err := uc.Execute(ctx, customerServiceInput)

	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, "existing-id", output.ID)
	mockRepo.AssertExpectations(t)
}

func TestUpdateCustomerService_Execute_Failure_ValidationError(t *testing.T) {
	mockRepo := mocks.NewCustomerServiceRepository(t)

	customerServiceInput := &customerservice.CustomerService{
		ServiceID: "updated-service-id",
	}

	uc := customerservice_usecase.NewUpdateCustomerService(mockRepo)

	ctx := context.Background()
	output, err := uc.Execute(ctx, customerServiceInput)
	assert.Nil(t, output)
	assert.NotNil(t, err)
	assert.Equal(t, errors.ParameterIsRequired, err.Code)
	mockRepo.AssertExpectations(t)
}

func TestUpdateCustomerService_Execute_Failure_UpdateError(t *testing.T) {
	mockRepo := mocks.NewCustomerServiceRepository(t)

	customerServiceInput := &customerservice.CustomerService{
		ID:         "existing-id",
		ServiceID:  "updated-service-id",
		CustomerID: "updated-customer-id",
		Name:       "Jane Doe",
		Document:   "98765432109",
		Email:      "jane.doe@example.com",
		Phone:      "987654321",
		Gender:     "female",
		BirthDate:  datetypes.CustomDate(time.Date(1992, 2, 2, 0, 0, 0, 0, time.UTC)),
	}

	mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(cs *customerservice.CustomerService) bool {
		return cs.ID == "existing-id" && cs.CustomerID == "updated-customer-id"
	})).Return(errors.NewError(errors.InternalServerError, nil, "Erro interno do servidor"))

	uc := customerservice_usecase.NewUpdateCustomerService(mockRepo)

	ctx := context.Background()
	output, err := uc.Execute(ctx, customerServiceInput)

	assert.Nil(t, output)
	assert.NotNil(t, err)
	assert.Equal(t, errors.InternalServerError, err.Code)
	mockRepo.AssertExpectations(t)
}
