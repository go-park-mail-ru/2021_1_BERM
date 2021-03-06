// Code generated by MockGen. DO NOT EDIT.
// Source: user/internal/app/review (interfaces: UseCase)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"
	models "user/internal/app/models"

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

// Create mocks base method.
func (m *MockUseCase) Create(arg0 models.Review, arg1 context.Context) (*models.Review, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(*models.Review)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockUseCaseMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUseCase)(nil).Create), arg0, arg1)
}

// GetAllReviewByUserId mocks base method.
func (m *MockUseCase) GetAllReviewByUserId(arg0 uint64, arg1 context.Context) (*models.UserReviews, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllReviewByUserId", arg0, arg1)
	ret0, _ := ret[0].(*models.UserReviews)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllReviewByUserId indicates an expected call of GetAllReviewByUserId.
func (mr *MockUseCaseMockRecorder) GetAllReviewByUserId(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllReviewByUserId", reflect.TypeOf((*MockUseCase)(nil).GetAllReviewByUserId), arg0, arg1)
}
