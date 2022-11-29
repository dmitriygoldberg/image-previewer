// Code generated by MockGen. DO NOT EDIT.
// Source: downloader.go

// Package mock_previewer is a generated GoMock package.
package previewer

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockImageDownloader is a mock of ImageDownloader interface.
type MockImageDownloader struct {
	ctrl     *gomock.Controller
	recorder *MockImageDownloaderMockRecorder
}

// MockImageDownloaderMockRecorder is the mock recorder for MockImageDownloader.
type MockImageDownloaderMockRecorder struct {
	mock *MockImageDownloader
}

// NewMockImageDownloader creates a new mock instance.
func NewMockImageDownloader(ctrl *gomock.Controller) *MockImageDownloader {
	mock := &MockImageDownloader{ctrl: ctrl}
	mock.recorder = &MockImageDownloaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockImageDownloader) EXPECT() *MockImageDownloaderMockRecorder {
	return m.recorder
}

// DownloadByURL mocks base method.
func (m *MockImageDownloader) DownloadByURL(ctx context.Context, url string, headers map[string][]string) (*Image, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DownloadByURL", ctx, url, headers)
	ret0, _ := ret[0].(*Image)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DownloadByURL indicates an expected call of DownloadByURL.
func (mr *MockImageDownloaderMockRecorder) DownloadByURL(ctx, url, headers interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DownloadByURL", reflect.TypeOf((*MockImageDownloader)(nil).DownloadByURL), ctx, url, headers)
}
