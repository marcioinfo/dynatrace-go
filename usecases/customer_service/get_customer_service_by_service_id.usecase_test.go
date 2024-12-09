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

func TestGetCustomerServiceByServiceID_Execute_Success(t *testing.T) {
	mockRepo := mocks.NewCustomerServiceRepository(t)

	serviceID := "service-id-1"
	expectedCustomerServices := []*customerservice.CustomerService{
		{
			ID:         "cs-id-1",
			ServiceID:  serviceID,
			CustomerID: "customer-id-1",
			Name:       "John Doe",
			Document:   "12345678901",
			Email:      "john.doe@example.com",
			Phone:      "123456789",
			Gender:     "male",
			BirthDate:  datetypes.CustomDate(time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)),
			CreatedAt:  datetypes.CustomDateTime(time.Now()),
			UpdatedAt:  datetypes.CustomDateTime(time.Now()),
		},
		{
			ID:         "cs-id-2",
			ServiceID:  serviceID,
			CustomerID: "customer-id-2",
			Name:       "Jane Doe",
			Document:   "10987654321",
			Email:      "jane.doe@example.com",
			Phone:      "987654321",
			Gender:     "female",
			BirthDate:  datetypes.CustomDate(time.Date(1992, 2, 2, 0, 0, 0, 0, time.UTC)),
			CreatedAt:  datetypes.CustomDateTime(time.Now()),
			UpdatedAt:  datetypes.CustomDateTime(time.Now()),
		},
	}

	mockRepo.On("GetByServiceID", mock.Anything, serviceID).Return(expectedCustomerServices, nil)

	uc := customerservice_usecase.NewGetCustomerServiceByServiceID(mockRepo)

	ctx := context.Background()
	output, err := uc.Execute(ctx, serviceID)

	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.Len(t, output, 2)
	assert.Equal(t, expectedCustomerServices, output)
	mockRepo.AssertExpectations(t)
}

func TestGetCustomerServiceByServiceID_Execute_Failure_NotFound(t *testing.T) {
	mockRepo := mocks.NewCustomerServiceRepository(t)

	serviceID := "non-existent-service-id"

	mockRepo.On("GetByServiceID", mock.Anything, serviceID).Return(nil, errors.NewError(errors.NotFoundError, nil, "customer service não encontrado"))

	uc := customerservice_usecase.NewGetCustomerServiceByServiceID(mockRepo)

	ctx := context.Background()
	output, err := uc.Execute(ctx, serviceID)

	assert.Nil(t, output)
	assert.NotNil(t, err)
	assert.Equal(t, errors.NotFoundError, err.Code)
	mockRepo.AssertExpectations(t)
}

func TestGetCustomerServiceByServiceID_Execute_Failure_DatabaseError(t *testing.T) {
	mockRepo := mocks.NewCustomerServiceRepository(t)

	serviceID := "service-id-1"

	mockRepo.On("GetByServiceID", mock.Anything, serviceID).Return(nil, errors.NewError(errors.InternalServerError, nil, "database error"))

	uc := customerservice_usecase.NewGetCustomerServiceByServiceID(mockRepo)

	ctx := context.Background()
	output, err := uc.Execute(ctx, serviceID)

	assert.Nil(t, output)
	assert.NotNil(t, err)
	assert.Equal(t, errors.InternalServerError, err.Code)
	mockRepo.AssertExpectations(t)
}
