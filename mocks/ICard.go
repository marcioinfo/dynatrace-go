// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// ICard is an autogenerated mock type for the ICard type
type ICard struct {
	mock.Mock
}

// GenerateFingerprint provides a mock function with given fields: salt, cardNumber
func (_m *ICard) GenerateFingerprint(salt string, cardNumber string) {
	_m.Called(salt, cardNumber)
}

// InitID provides a mock function with given fields:
func (_m *ICard) InitID() {
	_m.Called()
}

// IsValid provides a mock function with given fields:
func (_m *ICard) IsValid() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for IsValid")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewICard creates a new instance of ICard. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewICard(t interface {
	mock.TestingT
	Cleanup(func())
}) *ICard {
	mock := &ICard{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}