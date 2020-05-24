// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/asankov/gira/cmd/front-end/server (interfaces: Renderer)

// Package fixtures is a generated GoMock package.
package fixtures

import (
	gomock "github.com/golang/mock/gomock"
	http "net/http"
	reflect "reflect"
)

// RendererMock is a mock of Renderer interface
type RendererMock struct {
	ctrl     *gomock.Controller
	recorder *RendererMockMockRecorder
}

// RendererMockMockRecorder is the mock recorder for RendererMock
type RendererMockMockRecorder struct {
	mock *RendererMock
}

// NewRendererMock creates a new mock instance
func NewRendererMock(ctrl *gomock.Controller) *RendererMock {
	mock := &RendererMock{ctrl: ctrl}
	mock.recorder = &RendererMockMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *RendererMock) EXPECT() *RendererMockMockRecorder {
	return m.recorder
}

// Render mocks base method
func (m *RendererMock) Render(arg0 http.ResponseWriter, arg1 *http.Request, arg2 interface{}, arg3 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Render", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// Render indicates an expected call of Render
func (mr *RendererMockMockRecorder) Render(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Render", reflect.TypeOf((*RendererMock)(nil).Render), arg0, arg1, arg2, arg3)
}
