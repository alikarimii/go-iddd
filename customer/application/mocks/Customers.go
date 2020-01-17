// Code generated by mockery v1.0.0. DO NOT EDIT.

// +build test

package mocks

import mock "github.com/stretchr/testify/mock"
import shared "go-iddd/shared"
import values "go-iddd/customer/domain/values"

// Customers is an autogenerated mock type for the Customers type
type Customers struct {
	mock.Mock
}

// EventStream provides a mock function with given fields: id
func (_m *Customers) EventStream(id values.CustomerID) (shared.DomainEvents, error) {
	ret := _m.Called(id)

	var r0 shared.DomainEvents
	if rf, ok := ret.Get(0).(func(values.CustomerID) shared.DomainEvents); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(shared.DomainEvents)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(values.CustomerID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Persist provides a mock function with given fields: id, recordedEvents
func (_m *Customers) Persist(id values.CustomerID, recordedEvents shared.DomainEvents) error {
	ret := _m.Called(id, recordedEvents)

	var r0 error
	if rf, ok := ret.Get(0).(func(values.CustomerID, shared.DomainEvents) error); ok {
		r0 = rf(id, recordedEvents)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Register provides a mock function with given fields: id, recordedEvents
func (_m *Customers) Register(id values.CustomerID, recordedEvents shared.DomainEvents) error {
	ret := _m.Called(id, recordedEvents)

	var r0 error
	if rf, ok := ret.Get(0).(func(values.CustomerID, shared.DomainEvents) error); ok {
		r0 = rf(id, recordedEvents)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
