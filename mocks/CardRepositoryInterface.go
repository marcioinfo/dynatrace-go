// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"
	cards "payment-layer-card-api/entities/cards"

	errors "github.com/adhfoundation/payment-layer-error-package/pkg/errors"

	mock "github.com/stretchr/testify/mock"
)

// CardRepositoryInterface is an autogenerated mock type for the CardRepositoryInterface type
type CardRepositoryInterface struct {
	mock.Mock
}

// DeleteByID provides a mock function with given fields: ctx, id
func (_m *CardRepositoryInterface) DeleteByID(ctx context.Context, id string) (string, *errors.ErrorOutput) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for DeleteByID")
	}

	var r0 string
	var r1 *errors.ErrorOutput
	if rf, ok := ret.Get(0).(func(context.Context, string) (string, *errors.ErrorOutput)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) *errors.ErrorOutput); ok {
		r1 = rf(ctx, id)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*errors.ErrorOutput)
		}
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: ctx, id
func (_m *CardRepositoryInterface) GetByID(ctx context.Context, id string) (*cards.Card, *errors.ErrorOutput) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetByID")
	}

	var r0 *cards.Card
	var r1 *errors.ErrorOutput
	if rf, ok := ret.Get(0).(func(context.Context, string) (*cards.Card, *errors.ErrorOutput)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *cards.Card); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*cards.Card)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) *errors.ErrorOutput); ok {
		r1 = rf(ctx, id)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*errors.ErrorOutput)
		}
	}

	return r0, r1
}

// GetByIDAndCustomerID provides a mock function with given fields: ctx, id, customerId
func (_m *CardRepositoryInterface) GetByIDAndCustomerID(ctx context.Context, id string, customerId string) (*cards.Card, *errors.ErrorOutput) {
	ret := _m.Called(ctx, id, customerId)

	if len(ret) == 0 {
		panic("no return value specified for GetByIDAndCustomerID")
	}

	var r0 *cards.Card
	var r1 *errors.ErrorOutput
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (*cards.Card, *errors.ErrorOutput)); ok {
		return rf(ctx, id, customerId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *cards.Card); ok {
		r0 = rf(ctx, id, customerId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*cards.Card)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) *errors.ErrorOutput); ok {
		r1 = rf(ctx, id, customerId)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*errors.ErrorOutput)
		}
	}

	return r0, r1
}

// GetCardByFingerprint provides a mock function with given fields: ctx, fingerprint
func (_m *CardRepositoryInterface) GetCardByFingerprint(ctx context.Context, fingerprint string) (*cards.Card, *errors.ErrorOutput) {
	ret := _m.Called(ctx, fingerprint)

	if len(ret) == 0 {
		panic("no return value specified for GetCardByFingerprint")
	}

	var r0 *cards.Card
	var r1 *errors.ErrorOutput
	if rf, ok := ret.Get(0).(func(context.Context, string) (*cards.Card, *errors.ErrorOutput)); ok {
		return rf(ctx, fingerprint)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *cards.Card); ok {
		r0 = rf(ctx, fingerprint)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*cards.Card)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) *errors.ErrorOutput); ok {
		r1 = rf(ctx, fingerprint)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*errors.ErrorOutput)
		}
	}

	return r0, r1
}

// GetCardsByCustomerID provides a mock function with given fields: ctx, customerID
func (_m *CardRepositoryInterface) GetCardsByCustomerID(ctx context.Context, customerID string) ([]*cards.Card, *errors.ErrorOutput) {
	ret := _m.Called(ctx, customerID)

	if len(ret) == 0 {
		panic("no return value specified for GetCardsByCustomerID")
	}

	var r0 []*cards.Card
	var r1 *errors.ErrorOutput
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]*cards.Card, *errors.ErrorOutput)); ok {
		return rf(ctx, customerID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []*cards.Card); ok {
		r0 = rf(ctx, customerID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*cards.Card)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) *errors.ErrorOutput); ok {
		r1 = rf(ctx, customerID)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*errors.ErrorOutput)
		}
	}

	return r0, r1
}

// Insert provides a mock function with given fields: ctx, card
func (_m *CardRepositoryInterface) Insert(ctx context.Context, card *cards.Card) *errors.ErrorOutput {
	ret := _m.Called(ctx, card)

	if len(ret) == 0 {
		panic("no return value specified for Insert")
	}

	var r0 *errors.ErrorOutput
	if rf, ok := ret.Get(0).(func(context.Context, *cards.Card) *errors.ErrorOutput); ok {
		r0 = rf(ctx, card)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*errors.ErrorOutput)
		}
	}

	return r0
}

// NewCardRepositoryInterface creates a new instance of CardRepositoryInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCardRepositoryInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *CardRepositoryInterface {
	mock := &CardRepositoryInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}