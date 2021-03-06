// Code generated by MockGen. DO NOT EDIT.
// Source: authorizationservice/internal/app/session/tools (interfaces: SessionTools)

// Package mock is a generated GoMock package.
package mock

import (
	models "authorizationservice/internal/app/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockSessionTools is a mock of SessionTools interface.
type MockSessionTools struct {
	ctrl     *gomock.Controller
	recorder *MockSessionToolsMockRecorder
}

// MockSessionToolsMockRecorder is the mock recorder for MockSessionTools.
type MockSessionToolsMockRecorder struct {
	mock *MockSessionTools
}

// NewMockSessionTools creates a new mock instance.
func NewMockSessionTools(ctrl *gomock.Controller) *MockSessionTools {
	mock := &MockSessionTools{ctrl: ctrl}
	mock.recorder = &MockSessionToolsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSessionTools) EXPECT() *MockSessionToolsMockRecorder {
	return m.recorder
}

// BeforeCreate mocks base method.
func (m *MockSessionTools) BeforeCreate(arg0 models.Session) (models.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BeforeCreate", arg0)
	ret0, _ := ret[0].(models.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BeforeCreate indicates an expected call of BeforeCreate.
func (mr *MockSessionToolsMockRecorder) BeforeCreate(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BeforeCreate", reflect.TypeOf((*MockSessionTools)(nil).BeforeCreate), arg0)
}

// DecodingTarantoolToSession mocks base method.
func (m *MockSessionTools) DecodingTarantoolToSession(arg0 []interface{}) *models.Session {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DecodingTarantoolToSession", arg0)
	ret0, _ := ret[0].(*models.Session)
	return ret0
}

// DecodingTarantoolToSession indicates an expected call of DecodingTarantoolToSession.
func (mr *MockSessionToolsMockRecorder) DecodingTarantoolToSession(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DecodingTarantoolToSession", reflect.TypeOf((*MockSessionTools)(nil).DecodingTarantoolToSession), arg0)
}

// EncodingSessionToTarantool mocks base method.
func (m *MockSessionTools) EncodingSessionToTarantool(arg0 *models.Session) []interface{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EncodingSessionToTarantool", arg0)
	ret0, _ := ret[0].([]interface{})
	return ret0
}

// EncodingSessionToTarantool indicates an expected call of EncodingSessionToTarantool.
func (mr *MockSessionToolsMockRecorder) EncodingSessionToTarantool(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EncodingSessionToTarantool", reflect.TypeOf((*MockSessionTools)(nil).EncodingSessionToTarantool), arg0)
}
