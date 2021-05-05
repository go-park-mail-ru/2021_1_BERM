// Code generated by MockGen. DO NOT EDIT.
// Source: post/internal/app/vacancy (interfaces: Repository)

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
func (m *MockRepository) Change(arg0 models.Vacancy, arg1 context.Context) error {
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
func (m *MockRepository) Create(arg0 models.Vacancy, arg1 context.Context) (uint64, error) {
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
func (m *MockRepository) CreateArchive(arg0 models.Vacancy, arg1 context.Context) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateArchive", arg0, arg1)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateArchive indicates an expected call of CreateArchive.
func (mr *MockRepositoryMockRecorder) CreateArchive(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateArchive", reflect.TypeOf((*MockRepository)(nil).CreateArchive), arg0, arg1)
}

// DeleteVacancy mocks base method.
func (m *MockRepository) DeleteVacancy(arg0 uint64, arg1 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteVacancy", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteVacancy indicates an expected call of DeleteVacancy.
func (mr *MockRepositoryMockRecorder) DeleteVacancy(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteVacancy", reflect.TypeOf((*MockRepository)(nil).DeleteVacancy), arg0, arg1)
}

// FindByCustomerID mocks base method.
func (m *MockRepository) FindByCustomerID(arg0 uint64, arg1 context.Context) ([]models.Vacancy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByCustomerID", arg0, arg1)
	ret0, _ := ret[0].([]models.Vacancy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByCustomerID indicates an expected call of FindByCustomerID.
func (mr *MockRepositoryMockRecorder) FindByCustomerID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByCustomerID", reflect.TypeOf((*MockRepository)(nil).FindByCustomerID), arg0, arg1)
}

// FindByExecutorID mocks base method.
func (m *MockRepository) FindByExecutorID(arg0 uint64, arg1 context.Context) ([]models.Vacancy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByExecutorID", arg0, arg1)
	ret0, _ := ret[0].([]models.Vacancy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByExecutorID indicates an expected call of FindByExecutorID.
func (mr *MockRepositoryMockRecorder) FindByExecutorID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByExecutorID", reflect.TypeOf((*MockRepository)(nil).FindByExecutorID), arg0, arg1)
}

// FindByID mocks base method.
func (m *MockRepository) FindByID(arg0 uint64, arg1 context.Context) (*models.Vacancy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", arg0, arg1)
	ret0, _ := ret[0].(*models.Vacancy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockRepositoryMockRecorder) FindByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockRepository)(nil).FindByID), arg0, arg1)
}

// GetActualVacancies mocks base method.
func (m *MockRepository) GetActualVacancies(arg0 context.Context) ([]models.Vacancy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetActualVacancies", arg0)
	ret0, _ := ret[0].([]models.Vacancy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetActualVacancies indicates an expected call of GetActualVacancies.
func (mr *MockRepositoryMockRecorder) GetActualVacancies(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetActualVacancies", reflect.TypeOf((*MockRepository)(nil).GetActualVacancies), arg0)
}

// GetArchiveVacancies mocks base method.
func (m *MockRepository) GetArchiveVacancies(arg0 context.Context) ([]models.Vacancy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetArchiveVacancies", arg0)
	ret0, _ := ret[0].([]models.Vacancy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetArchiveVacancies indicates an expected call of GetArchiveVacancies.
func (mr *MockRepositoryMockRecorder) GetArchiveVacancies(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetArchiveVacancies", reflect.TypeOf((*MockRepository)(nil).GetArchiveVacancies), arg0)
}

// SearchVacancy mocks base method.
func (m *MockRepository) SearchVacancy(arg0 string, arg1 context.Context) ([]models.Vacancy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchVacancy", arg0, arg1)
	ret0, _ := ret[0].([]models.Vacancy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchVacancy indicates an expected call of SearchVacancy.
func (mr *MockRepositoryMockRecorder) SearchVacancy(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchVacancy", reflect.TypeOf((*MockRepository)(nil).SearchVacancy), arg0, arg1)
}

// UpdateExecutor mocks base method.
func (m *MockRepository) UpdateExecutor(arg0 models.Vacancy, arg1 context.Context) error {
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
