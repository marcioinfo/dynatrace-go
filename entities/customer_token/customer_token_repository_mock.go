package customer_token

import (
	"context"

	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
	"github.com/stretchr/testify/mock"
)

type CustomerTokenRepositoryMock struct {
	mock.Mock
}

func (m *CustomerTokenRepositoryMock) Insert(ctx context.Context, customerToken *CustomerToken) *errors.ErrorOutput {
	args := m.Called(customerToken)

	if err := args.Error(0); err != nil {
		errorOutput := &errors.ErrorOutput{
			Message: err.Error(),
			Code:    "",
		}
		return errorOutput
	}

	return nil
}

func (m *CustomerTokenRepositoryMock) GetByCustomerID(ctx context.Context, customerID string) ([]*CustomerToken, *errors.ErrorOutput) {
	args := m.Called(customerID)
	if args.Get(0) == nil {
		err := args.Error(1)
		errorOutput := &errors.ErrorOutput{
			Message: err.Error(),
			Code:    errors.NotFoundError,
		}
		return nil, errorOutput
	}
	return args.Get(0).([]*CustomerToken), nil
}

func (m *CustomerTokenRepositoryMock) GetByCustomerIDAndGateway(ctx context.Context, customerID string, gateway string) (*CustomerToken, *errors.ErrorOutput) {
	args := m.Called(customerID, gateway)
	if args.Get(0) == nil {
		err := args.Error(1)
		errorOutput := &errors.ErrorOutput{
			Message: err.Error(),
			Code:    errors.NotFoundError,
		}
		return nil, errorOutput
	}
	return args.Get(0).(*CustomerToken), nil
}
