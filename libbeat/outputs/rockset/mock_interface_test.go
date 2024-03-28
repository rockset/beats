// Code generated by MockGen. DO NOT EDIT.
// Source: rockset.go

// Package rockset is a generated GoMock package.
package rockset

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	openapi "github.com/rockset/rockset-go-client/openapi"
)

// MockRocker is a mock of Rocker interface.
type MockRocker struct {
	ctrl     *gomock.Controller
	recorder *MockRockerMockRecorder
}

// MockRockerMockRecorder is the mock recorder for MockRocker.
type MockRockerMockRecorder struct {
	mock *MockRocker
}

// NewMockRocker creates a new mock instance.
func NewMockRocker(ctrl *gomock.Controller) *MockRocker {
	mock := &MockRocker{ctrl: ctrl}
	mock.recorder = &MockRockerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRocker) EXPECT() *MockRockerMockRecorder {
	return m.recorder
}

// AddDocuments mocks base method.
func (m *MockRocker) AddDocuments(ctx context.Context, workspace, collection string, documents []interface{}) ([]openapi.DocumentStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddDocuments", ctx, workspace, collection, documents)
	ret0, _ := ret[0].([]openapi.DocumentStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddDocuments indicates an expected call of AddDocuments.
func (mr *MockRockerMockRecorder) AddDocuments(ctx, workspace, collection, documents interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddDocuments", reflect.TypeOf((*MockRocker)(nil).AddDocuments), ctx, workspace, collection, documents)
}

// GetOrganization mocks base method.
func (m *MockRocker) GetOrganization(ctx context.Context) (openapi.Organization, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrganization", ctx)
	ret0, _ := ret[0].(openapi.Organization)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrganization indicates an expected call of GetOrganization.
func (mr *MockRockerMockRecorder) GetOrganization(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrganization", reflect.TypeOf((*MockRocker)(nil).GetOrganization), ctx)
}
