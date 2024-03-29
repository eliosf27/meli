// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/item/item_service.go

// Package mock_item is a generated GoMock package.
package item

import (
	gomock "github.com/golang/mock/gomock"
	entities "meli/internal/entities"
	reflect "reflect"
)

// MockItemService is a mock of ItemService interface
type MockItemService struct {
	ctrl     *gomock.Controller
	recorder *MockItemServiceMockRecorder
}

// MockItemServiceMockRecorder is the mock recorder for MockItemService
type MockItemServiceMockRecorder struct {
	mock *MockItemService
}

// NewMockItemService creates a new mock instance
func NewMockItemService(ctrl *gomock.Controller) *MockItemService {
	mock := &MockItemService{ctrl: ctrl}
	mock.recorder = &MockItemServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockItemService) EXPECT() *MockItemServiceMockRecorder {
	return m.recorder
}

// FetchItemById mocks base method
func (m *MockItemService) FetchItemById(id string) entities.Item {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchItemById", id)
	ret0, _ := ret[0].(entities.Item)
	return ret0
}

// FetchItemById indicates an expected call of FetchItemById
func (mr *MockItemServiceMockRecorder) FetchItemById(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchItemById", reflect.TypeOf((*MockItemService)(nil).FetchItemById), id)
}
