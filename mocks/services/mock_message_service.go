// Code generated by MockGen. DO NOT EDIT.
// Source: message-app/services (interfaces: MessageService)

// Package services is a generated GoMock package.
package services

import (
	models "message-app/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockMessageService is a mock of MessageService interface.
type MockMessageService struct {
	ctrl     *gomock.Controller
	recorder *MockMessageServiceMockRecorder
}

// MockMessageServiceMockRecorder is the mock recorder for MockMessageService.
type MockMessageServiceMockRecorder struct {
	mock *MockMessageService
}

// NewMockMessageService creates a new mock instance.
func NewMockMessageService(ctrl *gomock.Controller) *MockMessageService {
	mock := &MockMessageService{ctrl: ctrl}
	mock.recorder = &MockMessageServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMessageService) EXPECT() *MockMessageServiceMockRecorder {
	return m.recorder
}

// CreateMessage mocks base method.
func (m *MockMessageService) CreateMessage(arg0 *models.CreateMessageReq) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateMessage", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateMessage indicates an expected call of CreateMessage.
func (mr *MockMessageServiceMockRecorder) CreateMessage(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMessage", reflect.TypeOf((*MockMessageService)(nil).CreateMessage), arg0)
}
