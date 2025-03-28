// Code generated by mockery v2.52.3. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	task_executor "github.com/uala-challenge/simple-toolkit/pkg/utilities/task_executor"

	types "github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// ExecuteTasks provides a mock function with given fields: ctx, messages
func (_m *Service) ExecuteTasks(ctx context.Context, messages []types.Message) map[string]task_executor.Result {
	ret := _m.Called(ctx, messages)

	if len(ret) == 0 {
		panic("no return value specified for ExecuteTasks")
	}

	var r0 map[string]task_executor.Result
	if rf, ok := ret.Get(0).(func(context.Context, []types.Message) map[string]task_executor.Result); ok {
		r0 = rf(ctx, messages)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]task_executor.Result)
		}
	}

	return r0
}

// NewService creates a new instance of Service. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewService(t interface {
	mock.TestingT
	Cleanup(func())
}) *Service {
	mock := &Service{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
