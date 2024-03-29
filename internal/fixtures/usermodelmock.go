// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/asankov/gira/cmd/api/server (interfaces: UserModel)

// Package fixtures is a generated GoMock package.
package fixtures

import (
	reflect "reflect"

	models "github.com/asankov/gira/pkg/models"
	gomock "github.com/golang/mock/gomock"
)

// UserModelMock is a mock of UserModel interface.
type UserModelMock struct {
	ctrl     *gomock.Controller
	recorder *UserModelMockMockRecorder
}

// UserModelMockMockRecorder is the mock recorder for UserModelMock.
type UserModelMockMockRecorder struct {
	mock *UserModelMock
}

// NewUserModelMock creates a new mock instance.
func NewUserModelMock(ctrl *gomock.Controller) *UserModelMock {
	mock := &UserModelMock{ctrl: ctrl}
	mock.recorder = &UserModelMockMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *UserModelMock) EXPECT() *UserModelMockMockRecorder {
	return m.recorder
}

// AssociateTokenWithUser mocks base method.
func (m *UserModelMock) AssociateTokenWithUser(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AssociateTokenWithUser", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AssociateTokenWithUser indicates an expected call of AssociateTokenWithUser.
func (mr *UserModelMockMockRecorder) AssociateTokenWithUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AssociateTokenWithUser", reflect.TypeOf((*UserModelMock)(nil).AssociateTokenWithUser), arg0, arg1)
}

// Authenticate mocks base method.
func (m *UserModelMock) Authenticate(arg0, arg1 string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Authenticate", arg0, arg1)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Authenticate indicates an expected call of Authenticate.
func (mr *UserModelMockMockRecorder) Authenticate(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Authenticate", reflect.TypeOf((*UserModelMock)(nil).Authenticate), arg0, arg1)
}

// GetUserByToken mocks base method.
func (m *UserModelMock) GetUserByToken(arg0 string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByToken", arg0)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByToken indicates an expected call of GetUserByToken.
func (mr *UserModelMockMockRecorder) GetUserByToken(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByToken", reflect.TypeOf((*UserModelMock)(nil).GetUserByToken), arg0)
}

// Insert mocks base method.
func (m *UserModelMock) Insert(arg0 *models.User) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", arg0)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Insert indicates an expected call of Insert.
func (mr *UserModelMockMockRecorder) Insert(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*UserModelMock)(nil).Insert), arg0)
}

// InvalidateToken mocks base method.
func (m *UserModelMock) InvalidateToken(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InvalidateToken", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// InvalidateToken indicates an expected call of InvalidateToken.
func (mr *UserModelMockMockRecorder) InvalidateToken(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InvalidateToken", reflect.TypeOf((*UserModelMock)(nil).InvalidateToken), arg0, arg1)
}
