// Code generated by mockery v2.13.1. DO NOT EDIT.

package mocks

import (
	meteorology "github.com/daniel5268/meliChallenge/src/domain/meteorology"
	mock "github.com/stretchr/testify/mock"
)

// ClimateRecordsRepository is an autogenerated mock type for the ClimateRecordsRepository type
type ClimateRecordsRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: _a0
func (_m *ClimateRecordsRepository) Create(_a0 ...*meteorology.ClimateRecord) error {
	_va := make([]interface{}, len(_a0))
	for _i := range _a0 {
		_va[_i] = _a0[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(...*meteorology.ClimateRecord) error); ok {
		r0 = rf(_a0...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindBetweenDays provides a mock function with given fields: _a0, _a1
func (_m *ClimateRecordsRepository) FindBetweenDays(_a0 int64, _a1 int64) ([]meteorology.ClimateRecord, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []meteorology.ClimateRecord
	if rf, ok := ret.Get(0).(func(int64, int64) []meteorology.ClimateRecord); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]meteorology.ClimateRecord)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64, int64) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindByDay provides a mock function with given fields: _a0
func (_m *ClimateRecordsRepository) FindByDay(_a0 int64) (meteorology.ClimateRecord, error) {
	ret := _m.Called(_a0)

	var r0 meteorology.ClimateRecord
	if rf, ok := ret.Get(0).(func(int64) meteorology.ClimateRecord); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(meteorology.ClimateRecord)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewClimateRecordsRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewClimateRecordsRepository creates a new instance of ClimateRecordsRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewClimateRecordsRepository(t mockConstructorTestingTNewClimateRecordsRepository) *ClimateRecordsRepository {
	mock := &ClimateRecordsRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
