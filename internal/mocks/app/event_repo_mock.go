// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ozonmp/ssn-service-api/internal/retranslator/repo (interfaces: EventRepo)

// Package appmocks is a generated GoMock package.
package appmocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	subscription "github.com/ozonmp/ssn-service-api/internal/model/subscription"
	repo "github.com/ozonmp/ssn-service-api/internal/retranslator/repo"
)

// MockEventRepo is a mock of EventRepo interface.
type MockEventRepo struct {
	ctrl     *gomock.Controller
	recorder *MockEventRepoMockRecorder
}

// MockEventRepoMockRecorder is the mock recorder for MockEventRepo.
type MockEventRepoMockRecorder struct {
	mock *MockEventRepo
}

// NewMockEventRepo creates a new mock instance.
func NewMockEventRepo(ctrl *gomock.Controller) *MockEventRepo {
	mock := &MockEventRepo{ctrl: ctrl}
	mock.recorder = &MockEventRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEventRepo) EXPECT() *MockEventRepoMockRecorder {
	return m.recorder
}

// Lock mocks base method.
func (m *MockEventRepo) Lock(arg0 context.Context, arg1 uint64, arg2 repo.QueryerExecer) ([]subscription.ServiceEvent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Lock", arg0, arg1, arg2)
	ret0, _ := ret[0].([]subscription.ServiceEvent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Lock indicates an expected call of Lock.
func (mr *MockEventRepoMockRecorder) Lock(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Lock", reflect.TypeOf((*MockEventRepo)(nil).Lock), arg0, arg1, arg2)
}

// Remove mocks base method.
func (m *MockEventRepo) Remove(arg0 context.Context, arg1 []uint64, arg2 repo.QueryerExecer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remove", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Remove indicates an expected call of Remove.
func (mr *MockEventRepoMockRecorder) Remove(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockEventRepo)(nil).Remove), arg0, arg1, arg2)
}

// Unlock mocks base method.
func (m *MockEventRepo) Unlock(arg0 context.Context, arg1 []uint64, arg2 repo.QueryerExecer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unlock", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Unlock indicates an expected call of Unlock.
func (mr *MockEventRepoMockRecorder) Unlock(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unlock", reflect.TypeOf((*MockEventRepo)(nil).Unlock), arg0, arg1, arg2)
}
