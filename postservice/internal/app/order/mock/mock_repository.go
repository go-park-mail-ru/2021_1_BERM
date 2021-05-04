// Code generated by MockGen. DO NOT EDIT.
// Source: post/internal/app/order (interfaces: Repository)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	models "post/internal/app/models"
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

// Change mocks base method.
func (m *MockRepository) Change(arg0 models.Order, arg1 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Change", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Change indicates an expected call of Change.
func (mr *MockRepositoryMockRecorder) Change(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Change", reflect.TypeOf((*MockRepository)(nil).Change), arg0, arg1)
}

// Create mocks base method.
func (m *MockRepository) Create(arg0 models.Order, arg1 context.Context) (uint64, error) {
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

// CreateArchive mocks base method.
func (m *MockRepository) CreateArchive(arg0 models.Order, arg1 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateArchive", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateArchive indicates an expected call of CreateArchive.
func (mr *MockRepositoryMockRecorder) CreateArchive(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateArchive", reflect.TypeOf((*MockRepository)(nil).CreateArchive), arg0, arg1)
}

// DeleteOrder mocks base method.
func (m *MockRepository) DeleteOrder(arg0 uint64, arg1 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteOrder", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteOrder indicates an expected call of DeleteOrder.
func (mr *MockRepositoryMockRecorder) DeleteOrder(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteOrder", reflect.TypeOf((*MockRepository)(nil).DeleteOrder), arg0, arg1)
}

// FindArchiveByID mocks base method.
func (m *MockRepository) FindArchiveByID(arg0 uint64, arg1 context.Context) (*models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindArchiveByID", arg0, arg1)
	ret0, _ := ret[0].(*models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindArchiveByID indicates an expected call of FindArchiveByID.
func (mr *MockRepositoryMockRecorder) FindArchiveByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindArchiveByID", reflect.TypeOf((*MockRepository)(nil).FindArchiveByID), arg0, arg1)
}

// FindByCustomerID mocks base method.
func (m *MockRepository) FindByCustomerID(arg0 uint64, arg1 context.Context) ([]models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByCustomerID", arg0, arg1)
	ret0, _ := ret[0].([]models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByCustomerID indicates an expected call of FindByCustomerID.
func (mr *MockRepositoryMockRecorder) FindByCustomerID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByCustomerID", reflect.TypeOf((*MockRepository)(nil).FindByCustomerID), arg0, arg1)
}

// FindByExecutorID mocks base method.
func (m *MockRepository) FindByExecutorID(arg0 uint64, arg1 context.Context) ([]models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByExecutorID", arg0, arg1)
	ret0, _ := ret[0].([]models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByExecutorID indicates an expected call of FindByExecutorID.
func (mr *MockRepositoryMockRecorder) FindByExecutorID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByExecutorID", reflect.TypeOf((*MockRepository)(nil).FindByExecutorID), arg0, arg1)
}

// FindByID mocks base method.
func (m *MockRepository) FindByID(arg0 uint64, arg1 context.Context) (*models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", arg0, arg1)
	ret0, _ := ret[0].(*models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockRepositoryMockRecorder) FindByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockRepository)(nil).FindByID), arg0, arg1)
}

// GetActualOrders mocks base method.
func (m *MockRepository) GetActualOrders(arg0 context.Context) ([]models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetActualOrders", arg0)
	ret0, _ := ret[0].([]models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetActualOrders indicates an expected call of GetActualOrders.
func (mr *MockRepositoryMockRecorder) GetActualOrders(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetActualOrders", reflect.TypeOf((*MockRepository)(nil).GetActualOrders), arg0)
}

// GetArchiveOrdersByCustomerID mocks base method.
func (m *MockRepository) GetArchiveOrdersByCustomerID(arg0 uint64, arg1 context.Context) ([]models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetArchiveOrdersByCustomerID", arg0, arg1)
	ret0, _ := ret[0].([]models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetArchiveOrdersByCustomerID indicates an expected call of GetArchiveOrdersByCustomerID.
func (mr *MockRepositoryMockRecorder) GetArchiveOrdersByCustomerID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetArchiveOrdersByCustomerID", reflect.TypeOf((*MockRepository)(nil).GetArchiveOrdersByCustomerID), arg0, arg1)
}

// GetArchiveOrdersByExecutorID mocks base method.
func (m *MockRepository) GetArchiveOrdersByExecutorID(arg0 uint64, arg1 context.Context) ([]models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetArchiveOrdersByExecutorID", arg0, arg1)
	ret0, _ := ret[0].([]models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetArchiveOrdersByExecutorID indicates an expected call of GetArchiveOrdersByExecutorID.
func (mr *MockRepositoryMockRecorder) GetArchiveOrdersByExecutorID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetArchiveOrdersByExecutorID", reflect.TypeOf((*MockRepository)(nil).GetArchiveOrdersByExecutorID), arg0, arg1)
}

// SearchOrders mocks base method.
func (m *MockRepository) SearchOrders(arg0 string, arg1 context.Context) ([]models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchOrders", arg0, arg1)
	ret0, _ := ret[0].([]models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchOrders indicates an expected call of SearchOrders.
func (mr *MockRepositoryMockRecorder) SearchOrders(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchOrders", reflect.TypeOf((*MockRepository)(nil).SearchOrders), arg0, arg1)
}

// UpdateExecutor mocks base method.
func (m *MockRepository) UpdateExecutor(arg0 models.Order, arg1 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateExecutor", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateExecutor indicates an expected call of UpdateExecutor.
func (mr *MockRepositoryMockRecorder) UpdateExecutor(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateExecutor", reflect.TypeOf((*MockRepository)(nil).UpdateExecutor), arg0, arg1)
}