// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"
	card_token "payment-layer-card-api/entities/card_token"

	errors "github.com/adhfoundation/payment-layer-error-package/pkg/errors"

	mock "github.com/stretchr/testify/mock"
)

// CardTokenRepositoryInterface is an autogenerated mock type for the CardTokenRepositoryInterface type
type CardTokenRepositoryInterface struct {
	mock.Mock
}

// DeleteByCardToken provides a mock function with given fields: ctx, cardID
func (_m *CardTokenRepositoryInterface) DeleteByCardToken(ctx context.Context, cardID string) *errors.ErrorOutput {
	ret := _m.Called(ctx, cardID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteByCardToken")
	}

	var r0 *errors.ErrorOutput
	if rf, ok := ret.Get(0).(func(context.Context, string) *errors.ErrorOutput); ok {
		r0 = rf(ctx, cardID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*errors.ErrorOutput)
		}
	}

	return r0
}

// GetByCardID provides a mock function with given fields: ctx, cardID
func (_m *CardTokenRepositoryInterface) GetByCardID(ctx context.Context, cardID string) ([]*card_token.CardToken, *errors.ErrorOutput) {
	ret := _m.Called(ctx, cardID)

	if len(ret) == 0 {
		panic("no return value specified for GetByCardID")
	}

	var r0 []*card_token.CardToken
	var r1 *errors.ErrorOutput
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]*card_token.CardToken, *errors.ErrorOutput)); ok {
		return rf(ctx, cardID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []*card_token.CardToken); ok {
		r0 = rf(ctx, cardID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*card_token.CardToken)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) *errors.ErrorOutput); ok {
		r1 = rf(ctx, cardID)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*errors.ErrorOutput)
		}
	}

	return r0, r1
}

// GetByCardIDAndGateway provides a mock function with given fields: ctx, cardID, gateway
func (_m *CardTokenRepositoryInterface) GetByCardIDAndGateway(ctx context.Context, cardID string, gateway string) (*card_token.CardToken, *errors.ErrorOutput) {
	ret := _m.Called(ctx, cardID, gateway)

	if len(ret) == 0 {
		panic("no return value specified for GetByCardIDAndGateway")
	}

	var r0 *card_token.CardToken
	var r1 *errors.ErrorOutput
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (*card_token.CardToken, *errors.ErrorOutput)); ok {
		return rf(ctx, cardID, gateway)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *card_token.CardToken); ok {
		r0 = rf(ctx, cardID, gateway)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*card_token.CardToken)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) *errors.ErrorOutput); ok {
		r1 = rf(ctx, cardID, gateway)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*errors.ErrorOutput)
		}
	}

	return r0, r1
}

// Insert provides a mock function with given fields: ctx, cardToken
func (_m *CardTokenRepositoryInterface) Insert(ctx context.Context, cardToken *card_token.CardToken) *errors.ErrorOutput {
	ret := _m.Called(ctx, cardToken)

	if len(ret) == 0 {
		panic("no return value specified for Insert")
	}

	var r0 *errors.ErrorOutput
	if rf, ok := ret.Get(0).(func(context.Context, *card_token.CardToken) *errors.ErrorOutput); ok {
		r0 = rf(ctx, cardToken)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*errors.ErrorOutput)
		}
	}

	return r0
}

// NewCardTokenRepositoryInterface creates a new instance of CardTokenRepositoryInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCardTokenRepositoryInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *CardTokenRepositoryInterface {
	mock := &CardTokenRepositoryInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}