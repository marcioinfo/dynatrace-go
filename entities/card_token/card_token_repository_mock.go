package card_token

import (
	"context"

	"github.com/adhfoundation/payment-layer-error-package/pkg/errors"
	"github.com/stretchr/testify/mock"
)

type CardTokenRepositoryMock struct {
	mock.Mock
}

func (m *CardTokenRepositoryMock) Insert(ctx context.Context, cardToken *CardToken) *errors.ErrorOutput {
	args := m.Called(cardToken)
	if err := args.Error(0); err != nil {
		return err.(*errors.ErrorOutput)
	}
	return nil
}

func (m *CardTokenRepositoryMock) GetByCardID(ctx context.Context, cardID string) ([]*CardToken, *errors.ErrorOutput) {
	args := m.Called(cardID)
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
	return args.Get(0).([]*CardToken), nil
}

func (m *CardTokenRepositoryMock) DeleteByCardToken(ctx context.Context, cardID string) *errors.ErrorOutput {
	args := m.Called(cardID)
	if err := args.Error(0); err != nil {
		return &errors.ErrorOutput{
			Message: err.Error(),
			Code:    "",
		}
	}
	return nil
}

func (m *CardTokenRepositoryMock) GetByCardIDAndGateway(ctx context.Context, cardID string, gateway string) (*CardToken, *errors.ErrorOutput) {
	args := m.Called(cardID, gateway)
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
	return args.Get(0).(*CardToken), nil
}
