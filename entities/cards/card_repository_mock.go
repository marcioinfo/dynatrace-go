package cards

import (
	"context"

	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
	"github.com/stretchr/testify/mock"
)

type CardRepositoryMock struct {
	mock.Mock
}

func (m *CardRepositoryMock) Insert(ctx context.Context, c *Card) *errors.ErrorOutput {
	args := m.Called(c)
	if err := args.Error(0); err != nil {
		return &errors.ErrorOutput{
			Message: err.Error(),
			Code:    "",
		}
	}
	return nil
}

func (m *CardRepositoryMock) GetCardsByCustomerID(ctx context.Context, customerID string) ([]*Card, *errors.ErrorOutput) {
	args := m.Called(customerID)
	if args.Get(0) == nil {
		return nil, &errors.ErrorOutput{
			Message: args.Error(1).Error(),
			Code:    "",
		}
	}
	return args.Get(0).([]*Card), nil
}

func (m *CardRepositoryMock) GetByID(ctx context.Context, id string) (*Card, *errors.ErrorOutput) {
	args := m.Called(id)
	if args.Get(0) == nil {
		if args.Error(1).Error() == "sql: no rows in result set" {
			return nil, &errors.ErrorOutput{
				Message: args.Error(1).Error(),
				Code:    errors.NotFoundError,
			}
		}
		return nil, &errors.ErrorOutput{
			Message: args.Error(1).Error(),
			Code:    errors.InternalServerError,
		}
	}
	return args.Get(0).(*Card), nil
}

func (m *CardRepositoryMock) GetByIDAndCustomerID(ctx context.Context, id string, customerId string) (*Card, *errors.ErrorOutput) {
	args := m.Called(id, customerId)
	if args.Get(0) == nil {
		return nil, &errors.ErrorOutput{
			Message: args.Error(1).Error(),
			Code:    "",
		}
	}
	return args.Get(0).(*Card), nil
}

func (m *CardRepositoryMock) DeleteByID(ctx context.Context, id string) (string, *errors.ErrorOutput) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return "", &errors.ErrorOutput{
			Message: args.Error(1).Error(),
			Code:    "",
		}
	}
	return args.Get(0).(string), nil
}

func (m *CardRepositoryMock) GetCardByFingerprint(ctx context.Context, fingerprint string) (*Card, *errors.ErrorOutput) {
	args := m.Called(fingerprint)
	if args.Get(0) == nil {
		if args.Error(1).Error() == "sql: no rows in result set" {
			return nil, &errors.ErrorOutput{
				Message: args.Error(1).Error(),
				Code:    errors.NotFoundError,
			}
		}
		return nil, &errors.ErrorOutput{
			Message: args.Error(1).Error(),
			Code:    errors.InternalServerError,
		}
	}

	return args.Get(0).(*Card), nil
}
