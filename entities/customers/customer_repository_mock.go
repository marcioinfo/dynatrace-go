package customers

import (
	"context"

	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
	"github.com/stretchr/testify/mock"
)

type CustomerRepositoryMock struct {
	mock.Mock
}

func (m *CustomerRepositoryMock) Insert(ctx context.Context, c *Customer) *errors.ErrorOutput {
	args := m.Called(c)
	if err := args.Error(0); err != nil {
		return &errors.ErrorOutput{Message: err.Error(), Code: ""}
	}
	return nil
}

func (m *CustomerRepositoryMock) GetCustomerByDocument(ctx context.Context, document string, serviceID string) (*Customer, *errors.ErrorOutput) {
	args := m.Called(document)
	if args.Get(0) == nil {
		return nil, &errors.ErrorOutput{Message: args.Error(1).Error(), Code: ""}
	}
	return args.Get(0).(*Customer), nil
}

func (m *CustomerRepositoryMock) GetByID(ctx context.Context, id string) (*Customer, *errors.ErrorOutput) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, &errors.ErrorOutput{Message: args.Error(1).Error(), Code: ""}
	}
	return args.Get(0).(*Customer), nil
}

func (m *CustomerRepositoryMock) GetCustomerByEmail(ctx context.Context, email string) (*Customer, *errors.ErrorOutput) {
	args := m.Called(email)
	if args.Get(0) == nil {
		if args.Error(1) == nil {
			return nil, nil
		}
		return nil, &errors.ErrorOutput{Message: args.Error(1).Error(), Code: ""}
	}
	return args.Get(0).(*Customer), nil
}

func (m *CustomerRepositoryMock) Update(ctx context.Context, c *Customer) *errors.ErrorOutput {
	args := m.Called(c)
	if err := args.Error(0); err != nil {
		return &errors.ErrorOutput{Message: err.Error(), Code: ""}
	}
	return nil
}
