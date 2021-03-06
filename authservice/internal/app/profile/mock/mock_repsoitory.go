// Code generated by MockGen. DO NOT EDIT.
// Source: authorizationservice/internal/app/profile (interfaces: Repository)

// Package mock is a generated GoMock package.
package mock

import (
	models "authorizationservice/internal/app/models"
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
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

// Authentication mocks base method.
func (m *MockRepository) Authentication(arg0, arg1 string, arg2 context.Context) (*models.UserBasicInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Authentication", arg0, arg1, arg2)
	ret0, _ := ret[0].(*models.UserBasicInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Authentication indicates an expected call of Authentication.
func (mr *MockRepositoryMockRecorder) Authentication(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Authentication", reflect.TypeOf((*MockRepository)(nil).Authentication), arg0, arg1, arg2)
}

// Create mocks base method.
func (m *MockRepository) Create(arg0 models.NewUser, arg1 context.Context) (*models.UserBasicInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(*models.UserBasicInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockRepositoryMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRepository)(nil).Create), arg0, arg1)
}
