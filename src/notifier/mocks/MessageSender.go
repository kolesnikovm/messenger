// Code generated by mockery v2.35.4. DO NOT EDIT.

package mocks

import (
	context "context"

	entity "github.com/kolesnikovm/messenger/entity"
	mock "github.com/stretchr/testify/mock"

	ulid "github.com/oklog/ulid/v2"
)

// MockMessageSender is an autogenerated mock type for the MessageSender type
type MockMessageSender struct {
	mock.Mock
}

type MockMessageSender_Expecter struct {
	mock *mock.Mock
}

func (_m *MockMessageSender) EXPECT() *MockMessageSender_Expecter {
	return &MockMessageSender_Expecter{mock: &_m.Mock}
}

// Get provides a mock function with given fields: ctx, recipientID, sessionID
func (_m *MockMessageSender) Get(ctx context.Context, recipientID string, sessionID ulid.ULID) (<-chan *entity.Message, func()) {
	ret := _m.Called(ctx, recipientID, sessionID)

	var r0 <-chan *entity.Message
	var r1 func()
	if rf, ok := ret.Get(0).(func(context.Context, string, ulid.ULID) (<-chan *entity.Message, func())); ok {
		return rf(ctx, recipientID, sessionID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, ulid.ULID) <-chan *entity.Message); ok {
		r0 = rf(ctx, recipientID, sessionID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan *entity.Message)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, ulid.ULID) func()); ok {
		r1 = rf(ctx, recipientID, sessionID)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(func())
		}
	}

	return r0, r1
}

// MockMessageSender_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type MockMessageSender_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - recipientID string
//   - sessionID ulid.ULID
func (_e *MockMessageSender_Expecter) Get(ctx interface{}, recipientID interface{}, sessionID interface{}) *MockMessageSender_Get_Call {
	return &MockMessageSender_Get_Call{Call: _e.mock.On("Get", ctx, recipientID, sessionID)}
}

func (_c *MockMessageSender_Get_Call) Run(run func(ctx context.Context, recipientID string, sessionID ulid.ULID)) *MockMessageSender_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(ulid.ULID))
	})
	return _c
}

func (_c *MockMessageSender_Get_Call) Return(stream <-chan *entity.Message, cleanup func()) *MockMessageSender_Get_Call {
	_c.Call.Return(stream, cleanup)
	return _c
}

func (_c *MockMessageSender_Get_Call) RunAndReturn(run func(context.Context, string, ulid.ULID) (<-chan *entity.Message, func())) *MockMessageSender_Get_Call {
	_c.Call.Return(run)
	return _c
}

// Send provides a mock function with given fields: _a0, _a1
func (_m *MockMessageSender) Send(_a0 context.Context, _a1 entity.Message) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, entity.Message) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockMessageSender_Send_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Send'
type MockMessageSender_Send_Call struct {
	*mock.Call
}

// Send is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 entity.Message
func (_e *MockMessageSender_Expecter) Send(_a0 interface{}, _a1 interface{}) *MockMessageSender_Send_Call {
	return &MockMessageSender_Send_Call{Call: _e.mock.On("Send", _a0, _a1)}
}

func (_c *MockMessageSender_Send_Call) Run(run func(_a0 context.Context, _a1 entity.Message)) *MockMessageSender_Send_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(entity.Message))
	})
	return _c
}

func (_c *MockMessageSender_Send_Call) Return(_a0 error) *MockMessageSender_Send_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockMessageSender_Send_Call) RunAndReturn(run func(context.Context, entity.Message) error) *MockMessageSender_Send_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockMessageSender creates a new instance of MockMessageSender. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockMessageSender(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockMessageSender {
	mock := &MockMessageSender{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
