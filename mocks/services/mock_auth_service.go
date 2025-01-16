// Code generated by MockGen. DO NOT EDIT.
// Source: message-app/services (interfaces: AuthenticationService)

// Package services is a generated GoMock package.
package services

import (
	models "message-app/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockAuthenticationService is a mock of AuthenticationService interface.
type MockAuthenticationService struct {
	ctrl     *gomock.Controller
	recorder *MockAuthenticationServiceMockRecorder
}

// MockAuthenticationServiceMockRecorder is the mock recorder for MockAuthenticationService.
type MockAuthenticationServiceMockRecorder struct {
	mock *MockAuthenticationService
}

// NewMockAuthenticationService creates a new mock instance.
func NewMockAuthenticationService(ctrl *gomock.Controller) *MockAuthenticationService {
	mock := &MockAuthenticationService{ctrl: ctrl}
	mock.recorder = &MockAuthenticationServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthenticationService) EXPECT() *MockAuthenticationServiceMockRecorder {
	return m.recorder
}

// CreateAuthToken mocks base method.
func (m *MockAuthenticationService) CreateAuthToken(arg0 *models.AuthenticationReq) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAuthToken", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAuthToken indicates an expected call of CreateAuthToken.
func (mr *MockAuthenticationServiceMockRecorder) CreateAuthToken(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAuthToken", reflect.TypeOf((*MockAuthenticationService)(nil).CreateAuthToken), arg0)
}
