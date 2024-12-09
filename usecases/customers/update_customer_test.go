package customers_usecase_test

import (
	"context"
	"payment-layer-card-api/entities/customers"
	"payment-layer-card-api/mocks"
	customers_usecase "payment-layer-card-api/usecases/customers"
	customer_dto "payment-layer-card-api/usecases/customers/dtos"
	"testing"

	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func stringPtr(s string) *string {
	return &s
}

func TestUpdateCustomerUsecase_Execute(t *testing.T) {
	ctx := context.TODO()

	tests := []struct {
		name           string
		setupMocks     func(customerRepo *mocks.CustomerRepositoryInterface)
		input          *customer_dto.UpdateCustomerDTO
		id             string
		expectedError  *errors.ErrorOutput
		expectedResult *customers.Customer
	}{
		{
			name: "successful update",
			setupMocks: func(customerRepo *mocks.CustomerRepositoryInterface) {
				customerRepo.On("GetByID", mock.Anything, "1").Return(&customers.Customer{
					ID:    "1",
					Email: "oldemail@example.com",
				}, nil)
				customerRepo.On("GetCustomerByEmail", mock.Anything, "newemail@example.com").Return(nil, nil)
				customerRepo.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
			input: &customer_dto.UpdateCustomerDTO{
				Email: stringPtr("newemail@example.com"),
			},
			id:            "1",
			expectedError: nil,
			expectedResult: &customers.Customer{
				ID:    "1",
				Email: "newemail@example.com",
			},
		},
		{
			name: "email already exists",
			setupMocks: func(customerRepo *mocks.CustomerRepositoryInterface) {
				customerRepo.On("GetByID", mock.Anything, "1").Return(&customers.Customer{
					ID:    "1",
					Email: "oldemail@example.com",
				}, nil)
				customerRepo.On("GetCustomerByEmail", mock.Anything, "existingemail@example.com").Return(&customers.Customer{
					ID: "2",
				}, nil)
			},
			input: &customer_dto.UpdateCustomerDTO{
				Email: stringPtr("existingemail@example.com"),
			},
			id:             "1",
			expectedError:  errors.NewPaymentLayerError(errors.AlredyExists, "Email already exists"),
			expectedResult: nil,
		},
		{
			name: "customer não encontrado",
			setupMocks: func(customerRepo *mocks.CustomerRepositoryInterface) {
				customerRepo.On("GetByID", mock.Anything, "1").Return(nil, errors.NewPaymentLayerError(errors.NotFoundError, "Customer não encontrado"))
			},
			input: &customer_dto.UpdateCustomerDTO{
				Email: stringPtr("newemail@example.com"),
			},
			id:             "1",
			expectedError:  errors.NewPaymentLayerError(errors.NotFoundError, "Customer não encontrado"),
			expectedResult: nil,
		},
		{
			name: "repository update error",
			setupMocks: func(customerRepo *mocks.CustomerRepositoryInterface) {
				customerRepo.On("GetByID", mock.Anything, "1").Return(&customers.Customer{
					ID:    "1",
					Email: "oldemail@example.com",
				}, nil)
				customerRepo.On("GetCustomerByEmail", mock.Anything, "newemail@example.com").Return(nil, nil)
				customerRepo.On("Update", mock.Anything, mock.Anything).Return(errors.NewPaymentLayerError(errors.InternalServerError, "Internal server error"))
			},
			input: &customer_dto.UpdateCustomerDTO{
				Email: stringPtr("newemail@example.com"),
			},
			id:             "1",
			expectedError:  errors.NewPaymentLayerError(errors.InternalServerError, "Internal server error"),
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			customerRepo := new(mocks.CustomerRepositoryInterface)
			tt.setupMocks(customerRepo)

			uc := customers_usecase.NewUpdateCustomerUsecase(customerRepo)

			result, err := uc.Execute(ctx, tt.input, tt.id)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Code, err.Code)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}

			customerRepo.AssertExpectations(t)
		})
	}
}
