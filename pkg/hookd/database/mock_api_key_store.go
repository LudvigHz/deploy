// Code generated by mockery v1.0.0. DO NOT EDIT.

package database

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockApiKeyStore is an autogenerated mock type for the ApiKeyStore type
type MockApiKeyStore struct {
	mock.Mock
}

// ApiKeys provides a mock function with given fields: ctx, id
func (_m *MockApiKeyStore) ApiKeys(ctx context.Context, id string) (ApiKeys, error) {
	ret := _m.Called(ctx, id)

	var r0 ApiKeys
	if rf, ok := ret.Get(0).(func(context.Context, string) ApiKeys); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(ApiKeys)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RotateApiKey provides a mock function with given fields: ctx, team, groupId, key
func (_m *MockApiKeyStore) RotateApiKey(ctx context.Context, team string, groupId string, key []byte) error {
	ret := _m.Called(ctx, team, groupId, key)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, []byte) error); ok {
		r0 = rf(ctx, team, groupId, key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}