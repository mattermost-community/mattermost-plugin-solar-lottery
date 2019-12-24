// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mattermost/mattermost-plugin-solar-lottery/server/store (interfaces: RotationStore)

// Package mock_store is a generated GoMock package.
package mock_store

import (
	gomock "github.com/golang/mock/gomock"
	store "github.com/mattermost/mattermost-plugin-solar-lottery/server/store"
	reflect "reflect"
)

// MockRotationStore is a mock of RotationStore interface
type MockRotationStore struct {
	ctrl     *gomock.Controller
	recorder *MockRotationStoreMockRecorder
}

// MockRotationStoreMockRecorder is the mock recorder for MockRotationStore
type MockRotationStoreMockRecorder struct {
	mock *MockRotationStore
}

// NewMockRotationStore creates a new mock instance
func NewMockRotationStore(ctrl *gomock.Controller) *MockRotationStore {
	mock := &MockRotationStore{ctrl: ctrl}
	mock.recorder = &MockRotationStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRotationStore) EXPECT() *MockRotationStoreMockRecorder {
	return m.recorder
}

// DeleteRotation mocks base method
func (m *MockRotationStore) DeleteRotation(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteRotation", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteRotation indicates an expected call of DeleteRotation
func (mr *MockRotationStoreMockRecorder) DeleteRotation(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteRotation", reflect.TypeOf((*MockRotationStore)(nil).DeleteRotation), arg0)
}

// LoadKnownRotations mocks base method
func (m *MockRotationStore) LoadKnownRotations() (store.IDMap, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadKnownRotations")
	ret0, _ := ret[0].(store.IDMap)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadKnownRotations indicates an expected call of LoadKnownRotations
func (mr *MockRotationStoreMockRecorder) LoadKnownRotations() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadKnownRotations", reflect.TypeOf((*MockRotationStore)(nil).LoadKnownRotations))
}

// LoadRotation mocks base method
func (m *MockRotationStore) LoadRotation(arg0 string) (*store.Rotation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadRotation", arg0)
	ret0, _ := ret[0].(*store.Rotation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadRotation indicates an expected call of LoadRotation
func (mr *MockRotationStoreMockRecorder) LoadRotation(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadRotation", reflect.TypeOf((*MockRotationStore)(nil).LoadRotation), arg0)
}

// StoreKnownRotations mocks base method
func (m *MockRotationStore) StoreKnownRotations(arg0 store.IDMap) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreKnownRotations", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// StoreKnownRotations indicates an expected call of StoreKnownRotations
func (mr *MockRotationStoreMockRecorder) StoreKnownRotations(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreKnownRotations", reflect.TypeOf((*MockRotationStore)(nil).StoreKnownRotations), arg0)
}

// StoreRotation mocks base method
func (m *MockRotationStore) StoreRotation(arg0 *store.Rotation) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreRotation", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// StoreRotation indicates an expected call of StoreRotation
func (mr *MockRotationStoreMockRecorder) StoreRotation(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreRotation", reflect.TypeOf((*MockRotationStore)(nil).StoreRotation), arg0)
}
