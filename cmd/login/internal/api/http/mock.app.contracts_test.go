// Code generated by MockGen. DO NOT EDIT.
// Source: http.go

// Package api_test is a generated GoMock package.
package http_test

import (
	context "context"
	app "medods/cmd/login/internal/app"
	reflect "reflect"

	uuid "github.com/gofrs/uuid"
	gomock "github.com/golang/mock/gomock"
)

// Mockapplication is a mock of application interface.
type Mockapplication struct {
	ctrl     *gomock.Controller
	recorder *MockapplicationMockRecorder
}

// MockapplicationMockRecorder is the mock recorder for Mockapplication.
type MockapplicationMockRecorder struct {
	mock *Mockapplication
}

// NewMockapplication creates a new mock instance.
func NewMockapplication(ctrl *gomock.Controller) *Mockapplication {
	mock := &Mockapplication{ctrl: ctrl}
	mock.recorder = &MockapplicationMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockapplication) EXPECT() *MockapplicationMockRecorder {
	return m.recorder
}

// Login mocks base method.
func (m *Mockapplication) Login(ctx context.Context, userID uuid.UUID, ip string) (*app.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", ctx, userID, ip)
	ret0, _ := ret[0].(*app.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockapplicationMockRecorder) Login(ctx, userID, ip interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*Mockapplication)(nil).Login), ctx, userID, ip)
}

// Refresh mocks base method.
func (m *Mockapplication) Refresh(ctx context.Context, accessToken, refreshToken, ip string) (*app.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Refresh", ctx, accessToken, refreshToken, ip)
	ret0, _ := ret[0].(*app.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Refresh indicates an expected call of Refresh.
func (mr *MockapplicationMockRecorder) Refresh(ctx, accessToken, refreshToken, ip interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Refresh", reflect.TypeOf((*Mockapplication)(nil).Refresh), ctx, accessToken, refreshToken, ip)
}
