// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/item/item_handler.go

// Package mock_item is a generated GoMock package.
package mock_item

import (
	gomock "github.com/golang/mock/gomock"
	echo "github.com/labstack/echo/v4"
	reflect "reflect"
)

// MockItemHandler is a mock of ItemHandler interface
type MockItemHandler struct {
	ctrl     *gomock.Controller
	recorder *MockItemHandlerMockRecorder
}

// MockItemHandlerMockRecorder is the mock recorder for MockItemHandler
type MockItemHandlerMockRecorder struct {
	mock *MockItemHandler
}

// NewMockItemHandler creates a new mock instance
func NewMockItemHandler(ctrl *gomock.Controller) *MockItemHandler {
	mock := &MockItemHandler{ctrl: ctrl}
	mock.recorder = &MockItemHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockItemHandler) EXPECT() *MockItemHandlerMockRecorder {
	return m.recorder
}

// Get mocks base method
func (m *MockItemHandler) Get(c echo.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", c)
	ret0, _ := ret[0].(error)
	return ret0
}

// Get indicates an expected call of Get
func (mr *MockItemHandlerMockRecorder) Get(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockItemHandler)(nil).Get), c)
}