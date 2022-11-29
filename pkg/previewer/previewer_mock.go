// Code generated by MockGen. DO NOT EDIT.
// Source: previewer.go

// Package mock_previewer is a generated GoMock package.
package previewer

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockPreviewInterface is a mock of PreviewInterface interface.
type MockPreviewInterface struct {
	ctrl     *gomock.Controller
	recorder *MockPreviewInterfaceMockRecorder
}

// MockPreviewInterfaceMockRecorder is the mock recorder for MockPreviewInterface.
type MockPreviewInterfaceMockRecorder struct {
	mock *MockPreviewInterface
}

// NewMockPreviewInterface creates a new mock instance.
func NewMockPreviewInterface(ctrl *gomock.Controller) *MockPreviewInterface {
	mock := &MockPreviewInterface{ctrl: ctrl}
	mock.recorder = &MockPreviewInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPreviewInterface) EXPECT() *MockPreviewInterfaceMockRecorder {
	return m.recorder
}

// Fill mocks base method.
func (m *MockPreviewInterface) Fill(params *FillParams) (*FillResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Fill", params)
	ret0, _ := ret[0].(*FillResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Fill indicates an expected call of Fill.
func (mr *MockPreviewInterfaceMockRecorder) Fill(params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fill", reflect.TypeOf((*MockPreviewInterface)(nil).Fill), params)
}
