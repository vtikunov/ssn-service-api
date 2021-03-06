// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ozonmp/ssn-service-api/internal/repo/subscription/service (interfaces: ServiceRepo)

// Package apimocks is a generated GoMock package.
package apimocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	subscription "github.com/ozonmp/ssn-service-api/internal/model/subscription"
	repo "github.com/ozonmp/ssn-service-api/internal/repo"
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
func (m *MockServiceRepo) Add(arg0 context.Context, arg1 *subscription.Service, arg2 repo.QueryerExecer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MockServiceRepoMockRecorder) Add(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockServiceRepo)(nil).Add), arg0, arg1, arg2)
}

// Describe mocks base method.
func (m *MockServiceRepo) Describe(arg0 context.Context, arg1 uint64, arg2 repo.QueryerExecer) (*subscription.Service, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Describe", arg0, arg1, arg2)
	ret0, _ := ret[0].(*subscription.Service)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Describe indicates an expected call of Describe.
func (mr *MockServiceRepoMockRecorder) Describe(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Describe", reflect.TypeOf((*MockServiceRepo)(nil).Describe), arg0, arg1, arg2)
}

// List mocks base method.
func (m *MockServiceRepo) List(arg0 context.Context, arg1, arg2 uint64, arg3 repo.QueryerExecer) ([]*subscription.Service, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].([]*subscription.Service)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockServiceRepoMockRecorder) List(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockServiceRepo)(nil).List), arg0, arg1, arg2, arg3)
}

// Remove mocks base method.
func (m *MockServiceRepo) Remove(arg0 context.Context, arg1 uint64, arg2 repo.QueryerExecer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remove", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Remove indicates an expected call of Remove.
func (mr *MockServiceRepoMockRecorder) Remove(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockServiceRepo)(nil).Remove), arg0, arg1, arg2)
}

// Update mocks base method.
func (m *MockServiceRepo) Update(arg0 context.Context, arg1 *subscription.Service, arg2 repo.QueryerExecer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockServiceRepoMockRecorder) Update(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockServiceRepo)(nil).Update), arg0, arg1, arg2)
}
