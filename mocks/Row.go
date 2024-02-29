// Code generated by mockery v2.42.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// Row is an autogenerated mock type for the Row type
type Row struct {
	mock.Mock
}

// Scan provides a mock function with given fields: dest
func (_m *Row) Scan(dest ...interface{}) error {
	var _ca []interface{}
	_ca = append(_ca, dest...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Scan")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(...interface{}) error); ok {
		r0 = rf(dest...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewRow creates a new instance of Row. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRow(t interface {
	mock.TestingT
	Cleanup(func())
}) *Row {
	mock := &Row{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
