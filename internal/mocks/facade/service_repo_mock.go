// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ozonmp/ssn-service-api/internal/facade/repo/subscription/service (interfaces: ServiceRepo)

// Package facademocks is a generated GoMock package.
package facademocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	subscription "github.com/ozonmp/ssn-service-api/internal/facade/model/subscription"
	repo "github.com/ozonmp/ssn-service-api/internal/facade/repo"
)

// MockServiceRepo is a mock of ServiceRepo interface.
type MockServiceRepo struct {
	ctrl     *gomock.Controller
	recorder *MockServiceRepoMockRecorder
}

// MockServiceRepoMockRecorder is the mock recorder for MockServiceRepo.
type MockServiceRepoMockRecorder struct {
	mock *MockServiceRepo
}

// NewMockServiceRepo creates a new mock instance.
func NewMockServiceRepo(ctrl *gomock.Controller) *MockServiceRepo {
	mock := &MockServiceRepo{ctrl: ctrl}
	mock.recorder = &MockServiceRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServiceRepo) EXPECT() *MockServiceRepoMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockServiceRepo) Add(arg0 context.Context, arg1 *subscription.Service, arg2 repo.QueryerExecer) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", arg0, arg1, arg2)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Add indicates an expected call of Add.
func (mr *MockServiceRepoMockRecorder) Add(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockServiceRepo)(nil).Add), arg0, arg1, arg2)
}

// Remove mocks base method.
func (m *MockServiceRepo) Remove(arg0 context.Context, arg1, arg2 uint64, arg3 repo.QueryerExecer) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remove", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Remove indicates an expected call of Remove.
func (mr *MockServiceRepoMockRecorder) Remove(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockServiceRepo)(nil).Remove), arg0, arg1, arg2, arg3)
}

// Update mocks base method.
func (m *MockServiceRepo) Update(arg0 context.Context, arg1 *subscription.Service, arg2 repo.QueryerExecer) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1, arg2)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockServiceRepoMockRecorder) Update(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockServiceRepo)(nil).Update), arg0, arg1, arg2)
}