// Code generated by MockGen. DO NOT EDIT.
// Source: FL_2/store (interfaces: UserRepository)

// Package mock is a generated GoMock package.
package mock

import (
	model "FL_2/model"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	pq "github.com/lib/pq"
)

// MockUserRepository is a mock of UserRepository interface.
type MockUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryMockRecorder
}

// MockUserRepositoryMockRecorder is the mock recorder for MockUserRepository.
type MockUserRepositoryMockRecorder struct {
	mock *MockUserRepository
}

// NewMockUserRepository creates a new mock instance.
func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
	mock := &MockUserRepository{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepository) EXPECT() *MockUserRepositoryMockRecorder {
	return m.recorder
}

// AddSpec mocks base method.
func (m *MockUserRepository) AddSpec(arg0 string) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddSpec", arg0)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddSpec indicates an expected call of AddSpec.
func (mr *MockUserRepositoryMockRecorder) AddSpec(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddSpec", reflect.TypeOf((*MockUserRepository)(nil).AddSpec), arg0)
}

// AddUser mocks base method.
func (m *MockUserRepository) AddUser(arg0 *model.User) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUser", arg0)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddUser indicates an expected call of AddUser.
func (mr *MockUserRepositoryMockRecorder) AddUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUser", reflect.TypeOf((*MockUserRepository)(nil).AddUser), arg0)
}

// AddUserSpec mocks base method.
func (m *MockUserRepository) AddUserSpec(arg0, arg1 uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUserSpec", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddUserSpec indicates an expected call of AddUserSpec.
func (mr *MockUserRepositoryMockRecorder) AddUserSpec(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUserSpec", reflect.TypeOf((*MockUserRepository)(nil).AddUserSpec), arg0, arg1)
}

// ChangeUser mocks base method.
func (m *MockUserRepository) ChangeUser(arg0 *model.User) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeUser", arg0)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ChangeUser indicates an expected call of ChangeUser.
func (mr *MockUserRepositoryMockRecorder) ChangeUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeUser", reflect.TypeOf((*MockUserRepository)(nil).ChangeUser), arg0)
}

// DelSpecialize mocks base method.
func (m *MockUserRepository) DelSpecialize(arg0, arg1 uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DelSpecialize", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DelSpecialize indicates an expected call of DelSpecialize.
func (mr *MockUserRepositoryMockRecorder) DelSpecialize(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DelSpecialize", reflect.TypeOf((*MockUserRepository)(nil).DelSpecialize), arg0, arg1)
}

// FindSpecializeByName mocks base method.
func (m *MockUserRepository) FindSpecializeByName(arg0 string) (model.Specialize, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindSpecializeByName", arg0)
	ret0, _ := ret[0].(model.Specialize)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindSpecializeByName indicates an expected call of FindSpecializeByName.
func (mr *MockUserRepositoryMockRecorder) FindSpecializeByName(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindSpecializeByName", reflect.TypeOf((*MockUserRepository)(nil).FindSpecializeByName), arg0)
}

// FindSpecializesByUserEmail mocks base method.
func (m *MockUserRepository) FindSpecializesByUserEmail(arg0 string) (pq.StringArray, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindSpecializesByUserEmail", arg0)
	ret0, _ := ret[0].(pq.StringArray)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindSpecializesByUserEmail indicates an expected call of FindSpecializesByUserEmail.
func (mr *MockUserRepositoryMockRecorder) FindSpecializesByUserEmail(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindSpecializesByUserEmail", reflect.TypeOf((*MockUserRepository)(nil).FindSpecializesByUserEmail), arg0)
}

// FindSpecializesByUserID mocks base method.
func (m *MockUserRepository) FindSpecializesByUserID(arg0 uint64) (pq.StringArray, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindSpecializesByUserID", arg0)
	ret0, _ := ret[0].(pq.StringArray)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindSpecializesByUserID indicates an expected call of FindSpecializesByUserID.
func (mr *MockUserRepositoryMockRecorder) FindSpecializesByUserID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindSpecializesByUserID", reflect.TypeOf((*MockUserRepository)(nil).FindSpecializesByUserID), arg0)
}

// FindUserByEmail mocks base method.
func (m *MockUserRepository) FindUserByEmail(arg0 string) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUserByEmail", arg0)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindUserByEmail indicates an expected call of FindUserByEmail.
func (mr *MockUserRepositoryMockRecorder) FindUserByEmail(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUserByEmail", reflect.TypeOf((*MockUserRepository)(nil).FindUserByEmail), arg0)
}

// FindUserByID mocks base method.
func (m *MockUserRepository) FindUserByID(arg0 uint64) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUserByID", arg0)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindUserByID indicates an expected call of FindUserByID.
func (mr *MockUserRepositoryMockRecorder) FindUserByID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUserByID", reflect.TypeOf((*MockUserRepository)(nil).FindUserByID), arg0)
}

// IsUserHaveSpec mocks base method.
func (m *MockUserRepository) IsUserHaveSpec(arg0, arg1 uint64) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsUserHaveSpec", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsUserHaveSpec indicates an expected call of IsUserHaveSpec.
func (mr *MockUserRepositoryMockRecorder) IsUserHaveSpec(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsUserHaveSpec", reflect.TypeOf((*MockUserRepository)(nil).IsUserHaveSpec), arg0, arg1)
}