// Code generated by mockery v2.35.4. DO NOT EDIT.

package mocks

import (
	context "context"

	entity "github.com/kolesnikovm/messenger/entity"
	mock "github.com/stretchr/testify/mock"

	ulid "github.com/oklog/ulid/v2"
)

// MockMessages is an autogenerated mock type for the Messages type
type MockMessages struct {
	mock.Mock
}

type MockMessages_Expecter struct {
	mock *mock.Mock
}

func (_m *MockMessages) EXPECT() *MockMessages_Expecter {
	return &MockMessages_Expecter{mock: &_m.Mock}
}

// BatchInsert provides a mock function with given fields: ctx, messages
func (_m *MockMessages) BatchInsert(ctx context.Context, messages []*entity.Message) error {
	ret := _m.Called(ctx, messages)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []*entity.Message) error); ok {
		r0 = rf(ctx, messages)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockMessages_BatchInsert_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'BatchInsert'
type MockMessages_BatchInsert_Call struct {
	*mock.Call
}

// BatchInsert is a helper method to define mock.On call
//   - ctx context.Context
//   - messages []*entity.Message
func (_e *MockMessages_Expecter) BatchInsert(ctx interface{}, messages interface{}) *MockMessages_BatchInsert_Call {
	return &MockMessages_BatchInsert_Call{Call: _e.mock.On("BatchInsert", ctx, messages)}
}

func (_c *MockMessages_BatchInsert_Call) Run(run func(ctx context.Context, messages []*entity.Message)) *MockMessages_BatchInsert_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]*entity.Message))
	})
	return _c
}

func (_c *MockMessages_BatchInsert_Call) Return(_a0 error) *MockMessages_BatchInsert_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockMessages_BatchInsert_Call) RunAndReturn(run func(context.Context, []*entity.Message) error) *MockMessages_BatchInsert_Call {
	_c.Call.Return(run)
	return _c
}

// GetMessageHistory provides a mock function with given fields: ctx, fromMessageID, chatID, messageCount, direction
func (_m *MockMessages) GetMessageHistory(ctx context.Context, fromMessageID ulid.ULID, chatID string, messageCount uint32, direction string) ([]*entity.Message, error) {
	ret := _m.Called(ctx, fromMessageID, chatID, messageCount, direction)

	var r0 []*entity.Message
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ulid.ULID, string, uint32, string) ([]*entity.Message, error)); ok {
		return rf(ctx, fromMessageID, chatID, messageCount, direction)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ulid.ULID, string, uint32, string) []*entity.Message); ok {
		r0 = rf(ctx, fromMessageID, chatID, messageCount, direction)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*entity.Message)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, ulid.ULID, string, uint32, string) error); ok {
		r1 = rf(ctx, fromMessageID, chatID, messageCount, direction)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockMessages_GetMessageHistory_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetMessageHistory'
type MockMessages_GetMessageHistory_Call struct {
	*mock.Call
}

// GetMessageHistory is a helper method to define mock.On call
//   - ctx context.Context
//   - fromMessageID ulid.ULID
//   - chatID string
//   - messageCount uint32
//   - direction string
func (_e *MockMessages_Expecter) GetMessageHistory(ctx interface{}, fromMessageID interface{}, chatID interface{}, messageCount interface{}, direction interface{}) *MockMessages_GetMessageHistory_Call {
	return &MockMessages_GetMessageHistory_Call{Call: _e.mock.On("GetMessageHistory", ctx, fromMessageID, chatID, messageCount, direction)}
}

func (_c *MockMessages_GetMessageHistory_Call) Run(run func(ctx context.Context, fromMessageID ulid.ULID, chatID string, messageCount uint32, direction string)) *MockMessages_GetMessageHistory_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(ulid.ULID), args[2].(string), args[3].(uint32), args[4].(string))
	})
	return _c
}

func (_c *MockMessages_GetMessageHistory_Call) Return(_a0 []*entity.Message, _a1 error) *MockMessages_GetMessageHistory_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockMessages_GetMessageHistory_Call) RunAndReturn(run func(context.Context, ulid.ULID, string, uint32, string) ([]*entity.Message, error)) *MockMessages_GetMessageHistory_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockMessages creates a new instance of MockMessages. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockMessages(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockMessages {
	mock := &MockMessages{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
