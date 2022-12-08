// Code generated by MockGen. DO NOT EDIT.
// Source: ./service/cake.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	cake "gitlab.com/cake-store-RESTFul/service/cake"
	common "gitlab.com/cake-store-RESTFul/service/common"
)

// MockCake is a mock of Cake interface.
type MockCake struct {
	ctrl     *gomock.Controller
	recorder *MockCakeMockRecorder
}

// MockCakeMockRecorder is the mock recorder for MockCake.
type MockCakeMockRecorder struct {
	mock *MockCake
}

// NewMockCake creates a new mock instance.
func NewMockCake(ctrl *gomock.Controller) *MockCake {
	mock := &MockCake{ctrl: ctrl}
	mock.recorder = &MockCakeMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCake) EXPECT() *MockCakeMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockCake) Create(ctx context.Context, req cake.CreateRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockCakeMockRecorder) Create(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockCake)(nil).Create), ctx, req)
}

// Delete mocks base method.
func (m *MockCake) Delete(ctx context.Context, id int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockCakeMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockCake)(nil).Delete), ctx, id)
}

// GetDetail mocks base method.
func (m *MockCake) GetDetail(ctx context.Context, id int) (cake.CakeResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDetail", ctx, id)
	ret0, _ := ret[0].(cake.CakeResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDetail indicates an expected call of GetDetail.
func (mr *MockCakeMockRecorder) GetDetail(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDetail", reflect.TypeOf((*MockCake)(nil).GetDetail), ctx, id)
}

// GetList mocks base method.
func (m *MockCake) GetList(ctx context.Context, req cake.GetListRequest, paginateReq common.PaginationRequest) (cake.CakesResponse, common.PaginationResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetList", ctx, req, paginateReq)
	ret0, _ := ret[0].(cake.CakesResponse)
	ret1, _ := ret[1].(common.PaginationResponse)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetList indicates an expected call of GetList.
func (mr *MockCakeMockRecorder) GetList(ctx, req, paginateReq interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetList", reflect.TypeOf((*MockCake)(nil).GetList), ctx, req, paginateReq)
}

// Update mocks base method.
func (m *MockCake) Update(ctx context.Context, req cake.UpdateRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockCakeMockRecorder) Update(ctx, req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockCake)(nil).Update), ctx, req)
}