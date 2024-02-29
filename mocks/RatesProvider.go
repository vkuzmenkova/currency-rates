// Code generated by mockery v2.42.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	models "github.com/vkuzmenkova/currency-rates/models"
)

// RatesProvider is an autogenerated mock type for the RatesProvider type
type RatesProvider struct {
	mock.Mock
}

// GetRate provides a mock function with given fields: base, currencyCode
func (_m *RatesProvider) GetRate(base string, currencyCode string) (*models.CurrencyRate, error) {
	ret := _m.Called(base, currencyCode)

	if len(ret) == 0 {
		panic("no return value specified for GetRate")
	}

	var r0 *models.CurrencyRate
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (*models.CurrencyRate, error)); ok {
		return rf(base, currencyCode)
	}
	if rf, ok := ret.Get(0).(func(string, string) *models.CurrencyRate); ok {
		r0 = rf(base, currencyCode)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.CurrencyRate)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(base, currencyCode)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewRatesProvider creates a new instance of RatesProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRatesProvider(t interface {
	mock.TestingT
	Cleanup(func())
}) *RatesProvider {
	mock := &RatesProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
