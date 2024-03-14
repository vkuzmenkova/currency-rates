// Code generated by mockery v2.42.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	models "github.com/vkuzmenkova/currency-rates/models"

	uuid "github.com/google/uuid"
)

// Repo is an autogenerated mock type for the Repo type
type Repo struct {
	mock.Mock
}

// GetLastRate provides a mock function with given fields: ctx, base, currencyCode
func (_m *Repo) GetLastRate(ctx context.Context, base uint8, currencyCode uint8) {
	_m.Called(ctx, base, currencyCode)
}

// GetRateByUUID provides a mock function with given fields: ctx, _a1
func (_m *Repo) GetRateByUUID(ctx context.Context, _a1 uuid.UUID) (models.CurrencyRate, error) {
	ret := _m.Called(ctx, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetRateByUUID")
	}

	var r0 models.CurrencyRate
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (models.CurrencyRate, error)); ok {
		return rf(ctx, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) models.CurrencyRate); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Get(0).(models.CurrencyRate)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(ctx, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InsertRate provides a mock function with given fields: ctx, uuidUpdate, base, currencyCode, value
func (_m *Repo) InsertRate(ctx context.Context, uuidUpdate string, base uint8, currencyCode uint8, value float64) error {
	ret := _m.Called(ctx, uuidUpdate, base, currencyCode, value)

	if len(ret) == 0 {
		panic("no return value specified for InsertRate")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, uint8, uint8, float64) error); ok {
		r0 = rf(ctx, uuidUpdate, base, currencyCode, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewRepo creates a new instance of Repo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRepo(t interface {
	mock.TestingT
	Cleanup(func())
}) *Repo {
	mock := &Repo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}