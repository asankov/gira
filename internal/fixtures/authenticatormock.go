// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/gira-games/api/cmd/api/server (interfaces: Authenticator)

// Package fixtures is a generated GoMock package.
package fixtures

import (
	models "github.com/gira-games/api/pkg/models"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// AuthenticatorMock is a mock of Authenticator interface
type AuthenticatorMock struct {
	ctrl     *gomock.Controller
	recorder *AuthenticatorMockMockRecorder
}

// AuthenticatorMockMockRecorder is the mock recorder for AuthenticatorMock
type AuthenticatorMockMockRecorder struct {
	mock *AuthenticatorMock
}

// NewAuthenticatorMock creates a new mock instance
func NewAuthenticatorMock(ctrl *gomock.Controller) *AuthenticatorMock {
	mock := &AuthenticatorMock{ctrl: ctrl}
	mock.recorder = &AuthenticatorMockMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *AuthenticatorMock) EXPECT() *AuthenticatorMockMockRecorder {
	return m.recorder
}

// DecodeToken mocks base method
func (m *AuthenticatorMock) DecodeToken(arg0 string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DecodeToken", arg0)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DecodeToken indicates an expected call of DecodeToken
func (mr *AuthenticatorMockMockRecorder) DecodeToken(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DecodeToken", reflect.TypeOf((*AuthenticatorMock)(nil).DecodeToken), arg0)
}

// NewTokenForUser mocks base method
func (m *AuthenticatorMock) NewTokenForUser(arg0 *models.User) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewTokenForUser", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewTokenForUser indicates an expected call of NewTokenForUser
func (mr *AuthenticatorMockMockRecorder) NewTokenForUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewTokenForUser", reflect.TypeOf((*AuthenticatorMock)(nil).NewTokenForUser), arg0)
}
