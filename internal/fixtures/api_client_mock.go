// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/asankov/gira/cmd/front-end/server (interfaces: APIClient)

// Package fixtures is a generated GoMock package.
package fixtures

import (
	client "github.com/asankov/gira/pkg/client"
	models "github.com/asankov/gira/pkg/models"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// APIClientMock is a mock of APIClient interface
type APIClientMock struct {
	ctrl     *gomock.Controller
	recorder *APIClientMockMockRecorder
}

// APIClientMockMockRecorder is the mock recorder for APIClientMock
type APIClientMockMockRecorder struct {
	mock *APIClientMock
}

// NewAPIClientMock creates a new mock instance
func NewAPIClientMock(ctrl *gomock.Controller) *APIClientMock {
	mock := &APIClientMock{ctrl: ctrl}
	mock.recorder = &APIClientMockMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *APIClientMock) EXPECT() *APIClientMockMockRecorder {
	return m.recorder
}

// ChangeGameProgress mocks base method
func (m *APIClientMock) ChangeGameProgress(arg0, arg1 string, arg2 *models.UserGameProgress) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeGameProgress", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangeGameProgress indicates an expected call of ChangeGameProgress
func (mr *APIClientMockMockRecorder) ChangeGameProgress(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeGameProgress", reflect.TypeOf((*APIClientMock)(nil).ChangeGameProgress), arg0, arg1, arg2)
}

// ChangeGameStatus mocks base method
func (m *APIClientMock) ChangeGameStatus(arg0, arg1 string, arg2 models.Status) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeGameStatus", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangeGameStatus indicates an expected call of ChangeGameStatus
func (mr *APIClientMockMockRecorder) ChangeGameStatus(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeGameStatus", reflect.TypeOf((*APIClientMock)(nil).ChangeGameStatus), arg0, arg1, arg2)
}

// CreateGame mocks base method
func (m *APIClientMock) CreateGame(arg0 *models.Game, arg1 string) (*models.Game, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateGame", arg0, arg1)
	ret0, _ := ret[0].(*models.Game)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateGame indicates an expected call of CreateGame
func (mr *APIClientMockMockRecorder) CreateGame(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateGame", reflect.TypeOf((*APIClientMock)(nil).CreateGame), arg0, arg1)
}

// CreateUser mocks base method
func (m *APIClientMock) CreateUser(arg0 *models.User) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser
func (mr *APIClientMockMockRecorder) CreateUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*APIClientMock)(nil).CreateUser), arg0)
}

// DeleteUserGame mocks base method
func (m *APIClientMock) DeleteUserGame(arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUserGame", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUserGame indicates an expected call of DeleteUserGame
func (mr *APIClientMockMockRecorder) DeleteUserGame(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUserGame", reflect.TypeOf((*APIClientMock)(nil).DeleteUserGame), arg0, arg1)
}

// GetGames mocks base method
func (m *APIClientMock) GetGames(arg0 string, arg1 *client.GetGamesOptions) ([]*models.Game, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGames", arg0, arg1)
	ret0, _ := ret[0].([]*models.Game)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGames indicates an expected call of GetGames
func (mr *APIClientMockMockRecorder) GetGames(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGames", reflect.TypeOf((*APIClientMock)(nil).GetGames), arg0, arg1)
}

// GetUser mocks base method
func (m *APIClientMock) GetUser(arg0 string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", arg0)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser
func (mr *APIClientMockMockRecorder) GetUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*APIClientMock)(nil).GetUser), arg0)
}

// GetUserGames mocks base method
func (m *APIClientMock) GetUserGames(arg0 string) (map[models.Status][]*models.UserGame, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserGames", arg0)
	ret0, _ := ret[0].(map[models.Status][]*models.UserGame)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserGames indicates an expected call of GetUserGames
func (mr *APIClientMockMockRecorder) GetUserGames(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserGames", reflect.TypeOf((*APIClientMock)(nil).GetUserGames), arg0)
}

// LinkGameToUser mocks base method
func (m *APIClientMock) LinkGameToUser(arg0, arg1 string) (*models.UserGame, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LinkGameToUser", arg0, arg1)
	ret0, _ := ret[0].(*models.UserGame)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LinkGameToUser indicates an expected call of LinkGameToUser
func (mr *APIClientMockMockRecorder) LinkGameToUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LinkGameToUser", reflect.TypeOf((*APIClientMock)(nil).LinkGameToUser), arg0, arg1)
}

// LoginUser mocks base method
func (m *APIClientMock) LoginUser(arg0 *models.User) (*models.UserLoginResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoginUser", arg0)
	ret0, _ := ret[0].(*models.UserLoginResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoginUser indicates an expected call of LoginUser
func (mr *APIClientMockMockRecorder) LoginUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoginUser", reflect.TypeOf((*APIClientMock)(nil).LoginUser), arg0)
}

// LogoutUser mocks base method
func (m *APIClientMock) LogoutUser(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LogoutUser", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// LogoutUser indicates an expected call of LogoutUser
func (mr *APIClientMockMockRecorder) LogoutUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogoutUser", reflect.TypeOf((*APIClientMock)(nil).LogoutUser), arg0)
}
