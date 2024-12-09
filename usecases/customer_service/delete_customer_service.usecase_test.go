package customerservice_usecase_test

import (
	"context"
	"testing"

	customerservice "payment-layer-card-api/entities/customer_service"
	"payment-layer-card-api/mocks"
	customerservice_usecase "payment-layer-card-api/usecases/customer_service"

	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteCustomerService_Execute_Success(t *testing.T) {
	mockRepo := mocks.NewCustomerServiceRepository(t)

	customerServiceID := "cs-id-1"

	mockRepo.On("Delete", mock.Anything, mock.MatchedBy(func(cs *customerservice.CustomerService) bool {
		return cs.ID == customerServiceID
	})).Return(nil)

	uc := customerservice_usecase.NewDeleteCustomerService(mockRepo)

	ctx := context.Background()
	output, err := uc.Execute(ctx, customerServiceID)

	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, customerServiceID, output.ID)
	mockRepo.AssertExpectations(t)
}

func TestDeleteCustomerService_Execute_Failure_NotFound(t *testing.T) {
	mockRepo := mocks.NewCustomerServiceRepository(t)

	customerServiceID := "non-existent-cs-id"

	mockRepo.On("Delete", mock.Anything, mock.MatchedBy(func(cs *customerservice.CustomerService) bool {
		return cs.ID == customerServiceID
	})).Return(errors.NewError(errors.NotFoundError, nil, "customer service não encontrado"))

	uc := customerservice_usecase.NewDeleteCustomerService(mockRepo)

	ctx := context.Background()
	output, err := uc.Execute(ctx, customerServiceID)

	assert.Nil(t, output)
	assert.NotNil(t, err)
	assert.Equal(t, errors.NotFoundError, err.Code)
	mockRepo.AssertExpectations(t)
}

func TestDeleteCustomerService_Execute_Failure_DatabaseError(t *testing.T) {
	mockRepo := mocks.NewCustomerServiceRepository(t)

	customerServiceID := "cs-id-1"

	mockRepo.On("Delete", mock.Anything, mock.MatchedBy(func(cs *customerservice.CustomerService) bool {
		return cs.ID == customerServiceID
	})).Return(errors.NewError(errors.InternalServerError, nil, "database error"))

	uc := customerservice_usecase.NewDeleteCustomerService(mockRepo)

	ctx := context.Background()
	output, err := uc.Execute(ctx, customerServiceID)

	assert.Nil(t, output)
	assert.NotNil(t, err)
	assert.Equal(t, errors.InternalServerError, err.Code)
	mockRepo.AssertExpectations(t)
}
