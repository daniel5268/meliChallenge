// Code generated by mockery v2.13.1. DO NOT EDIT.

package mocks

import (
	meteorology "github.com/daniel5268/meliChallenge/src/domain/meteorology"
	mock "github.com/stretchr/testify/mock"
)

// ClimateRecordsService is an autogenerated mock type for the ClimateRecordsService type
type ClimateRecordsService struct {
	mock.Mock
}

// GetClimateRecord provides a mock function with given fields: day
func (_m *ClimateRecordsService) GetClimateRecord(day int64) (meteorology.ClimateRecord, error) {
	ret := _m.Called(day)

	var r0 meteorology.ClimateRecord
	if rf, ok := ret.Get(0).(func(int64) meteorology.ClimateRecord); ok {
		r0 = rf(day)
	} else {
		r0 = ret.Get(0).(meteorology.ClimateRecord)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(day)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewClimateRecordsService interface {
	mock.TestingT
	Cleanup(func())
}

// NewClimateRecordsService creates a new instance of ClimateRecordsService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewClimateRecordsService(t mockConstructorTestingTNewClimateRecordsService) *ClimateRecordsService {
	mock := &ClimateRecordsService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
