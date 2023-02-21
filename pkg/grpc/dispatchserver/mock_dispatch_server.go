// Code generated by mockery v2.20.0. DO NOT EDIT.

package dispatchserver

import (
	context "context"

	pb "github.com/nais/deploy/pkg/pb"
	mock "github.com/stretchr/testify/mock"
)

// MockDispatchServer is an autogenerated mock type for the DispatchServer type
type MockDispatchServer struct {
	pb.UnimplementedDispatchServer
	mock.Mock
}

// Deployments provides a mock function with given fields: _a0, _a1
func (_m *MockDispatchServer) Deployments(_a0 *pb.GetDeploymentOpts, _a1 pb.Dispatch_DeploymentsServer) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(*pb.GetDeploymentOpts, pb.Dispatch_DeploymentsServer) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// HandleDeploymentStatus provides a mock function with given fields: ctx, status
func (_m *MockDispatchServer) HandleDeploymentStatus(ctx context.Context, status *pb.DeploymentStatus) error {
	ret := _m.Called(ctx, status)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *pb.DeploymentStatus) error); ok {
		r0 = rf(ctx, status)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ReportStatus provides a mock function with given fields: _a0, _a1
func (_m *MockDispatchServer) ReportStatus(_a0 context.Context, _a1 *pb.DeploymentStatus) (*pb.ReportStatusOpts, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *pb.ReportStatusOpts
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *pb.DeploymentStatus) (*pb.ReportStatusOpts, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *pb.DeploymentStatus) *pb.ReportStatusOpts); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*pb.ReportStatusOpts)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *pb.DeploymentStatus) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SendDeploymentRequest provides a mock function with given fields: ctx, deployment
func (_m *MockDispatchServer) SendDeploymentRequest(ctx context.Context, deployment *pb.DeploymentRequest) error {
	ret := _m.Called(ctx, deployment)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *pb.DeploymentRequest) error); ok {
		r0 = rf(ctx, deployment)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// StreamStatus provides a mock function with given fields: _a0, _a1
func (_m *MockDispatchServer) StreamStatus(_a0 context.Context, _a1 chan<- *pb.DeploymentStatus) {
	_m.Called(_a0, _a1)
}

// mustEmbedUnimplementedDispatchServer provides a mock function with given fields:
func (_m *MockDispatchServer) mustEmbedUnimplementedDispatchServer() {
	_m.Called()
}

type mockConstructorTestingTNewMockDispatchServer interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockDispatchServer creates a new instance of MockDispatchServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockDispatchServer(t mockConstructorTestingTNewMockDispatchServer) *MockDispatchServer {
	mock := &MockDispatchServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
