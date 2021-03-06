// Code generated by MockGen. DO NOT EDIT.
// Source: upsert.go

package upsert

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockUpsert is a mock of Upsert interface
type MockUpsert struct {
	ctrl     *gomock.Controller
	recorder *MockUpsertMockRecorder
}

// MockUpsertMockRecorder is the mock recorder for MockUpsert
type MockUpsertMockRecorder struct {
	mock *MockUpsert
}

// NewMockUpsert creates a new mock instance
func NewMockUpsert(ctrl *gomock.Controller) *MockUpsert {
	mock := &MockUpsert{ctrl: ctrl}
	mock.recorder = &MockUpsertMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUpsert) EXPECT() *MockUpsertMockRecorder {
	return m.recorder
}

// SendViaAPI mocks base method
func (m *MockUpsert) SendViaAPI(data string, apiKey string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendViaAPI", data)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendViaAPI indicates an expected call of SendViaAPI
func (mr *MockUpsertMockRecorder) SendViaAPI(data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendViaAPI", reflect.TypeOf((*MockUpsert)(nil).SendViaAPI), data)
}
