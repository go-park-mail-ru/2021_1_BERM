// Code generated by MockGen. DO NOT EDIT.
// Source: user/internal/app/specialize (interfaces: Repository)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	pq "github.com/lib/pq"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// AssociateSpecializationWithUser mocks base method.
func (m *MockRepository) AssociateSpecializationWithUser(arg0, arg1 uint64, arg2 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AssociateSpecializationWithUser", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AssociateSpecializationWithUser indicates an expected call of AssociateSpecializationWithUser.
func (mr *MockRepositoryMockRecorder) AssociateSpecializationWithUser(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AssociateSpecializationWithUser", reflect.TypeOf((*MockRepository)(nil).AssociateSpecializationWithUser), arg0, arg1, arg2)
}

// Create mocks base method.
func (m *MockRepository) Create(arg0 string, arg1 context.Context) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockRepositoryMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRepository)(nil).Create), arg0, arg1)
}

// FindByName mocks base method.
func (m *MockRepository) FindByName(arg0 string, arg1 context.Context) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByName", arg0, arg1)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByName indicates an expected call of FindByName.
func (mr *MockRepositoryMockRecorder) FindByName(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByName", reflect.TypeOf((*MockRepository)(nil).FindByName), arg0, arg1)
}

// FindByUserID mocks base method.
func (m *MockRepository) FindByUserID(arg0 uint64, arg1 context.Context) (pq.StringArray, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByUserID", arg0, arg1)
	ret0, _ := ret[0].(pq.StringArray)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByUserID indicates an expected call of FindByUserID.
func (mr *MockRepositoryMockRecorder) FindByUserID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByUserID", reflect.TypeOf((*MockRepository)(nil).FindByUserID), arg0, arg1)
}

// RemoveAssociateSpecializationWithUser mocks base method.
func (m *MockRepository) RemoveAssociateSpecializationWithUser(arg0, arg1 uint64, arg2 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveAssociateSpecializationWithUser", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveAssociateSpecializationWithUser indicates an expected call of RemoveAssociateSpecializationWithUser.
func (mr *MockRepositoryMockRecorder) RemoveAssociateSpecializationWithUser(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveAssociateSpecializationWithUser", reflect.TypeOf((*MockRepository)(nil).RemoveAssociateSpecializationWithUser), arg0, arg1, arg2)
}