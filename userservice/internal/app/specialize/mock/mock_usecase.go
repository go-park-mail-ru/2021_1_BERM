// Code generated by MockGen. DO NOT EDIT.
// Source: user/internal/app/specialize (interfaces: UseCase)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockUseCase is a mock of UseCase interface.
type MockUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockUseCaseMockRecorder
}

// MockUseCaseMockRecorder is the mock recorder for MockUseCase.
type MockUseCaseMockRecorder struct {
	mock *MockUseCase
}

// NewMockUseCase creates a new mock instance.
func NewMockUseCase(ctrl *gomock.Controller) *MockUseCase {
	mock := &MockUseCase{ctrl: ctrl}
	mock.recorder = &MockUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUseCase) EXPECT() *MockUseCaseMockRecorder {
	return m.recorder
}

// AssociateWithUser mocks base method.
func (m *MockUseCase) AssociateWithUser(arg0 uint64, arg1 string, arg2 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AssociateWithUser", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AssociateWithUser indicates an expected call of AssociateWithUser.
func (mr *MockUseCaseMockRecorder) AssociateWithUser(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AssociateWithUser", reflect.TypeOf((*MockUseCase)(nil).AssociateWithUser), arg0, arg1, arg2)
}

// Create mocks base method.
func (m *MockUseCase) Create(arg0 string, arg1 context.Context) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockUseCaseMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUseCase)(nil).Create), arg0, arg1)
}

// Remove mocks base method.
func (m *MockUseCase) Remove(arg0 uint64, arg1 string, arg2 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remove", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Remove indicates an expected call of Remove.
func (mr *MockUseCaseMockRecorder) Remove(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockUseCase)(nil).Remove), arg0, arg1, arg2)
}
