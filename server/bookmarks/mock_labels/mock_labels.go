// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/jfrerich/mattermost-plugin-bookmarks/server/bookmarks (interfaces: ILabels)

// Package mock_bookmarks is a generated GoMock package.
package mock_bookmarks

import (
	gomock "github.com/golang/mock/gomock"
	bookmarks "github.com/jfrerich/mattermost-plugin-bookmarks/server/bookmarks"
	reflect "reflect"
)

// MockILabels is a mock of ILabels interface
type MockILabels struct {
	ctrl     *gomock.Controller
	recorder *MockILabelsMockRecorder
}

// MockILabelsMockRecorder is the mock recorder for MockILabels
type MockILabelsMockRecorder struct {
	mock *MockILabels
}

// NewMockILabels creates a new mock instance
func NewMockILabels(ctrl *gomock.Controller) *MockILabels {
	mock := &MockILabels{ctrl: ctrl}
	mock.recorder = &MockILabelsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockILabels) EXPECT() *MockILabelsMockRecorder {
	return m.recorder
}

// AddLabel mocks base method
func (m *MockILabels) AddLabel(arg0 string) (*bookmarks.Label, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddLabel", arg0)
	ret0, _ := ret[0].(*bookmarks.Label)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddLabel indicates an expected call of AddLabel
func (mr *MockILabelsMockRecorder) AddLabel(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddLabel", reflect.TypeOf((*MockILabels)(nil).AddLabel), arg0)
}

// DeleteByID mocks base method
func (m *MockILabels) DeleteByID(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteByID", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteByID indicates an expected call of DeleteByID
func (mr *MockILabelsMockRecorder) DeleteByID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteByID", reflect.TypeOf((*MockILabels)(nil).DeleteByID), arg0)
}

// GetByIDs mocks base method
func (m *MockILabels) GetByIDs() map[string]*bookmarks.Label {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByIDs")
	ret0, _ := ret[0].(map[string]*bookmarks.Label)
	return ret0
}

// GetByIDs indicates an expected call of GetByIDs
func (mr *MockILabelsMockRecorder) GetByIDs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByIDs", reflect.TypeOf((*MockILabels)(nil).GetByIDs))
}

// GetIDFromName mocks base method
func (m *MockILabels) GetIDFromName(arg0 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIDFromName", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetIDFromName indicates an expected call of GetIDFromName
func (mr *MockILabelsMockRecorder) GetIDFromName(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIDFromName", reflect.TypeOf((*MockILabels)(nil).GetIDFromName), arg0)
}

// GetLabelByName mocks base method
func (m *MockILabels) GetLabelByName(arg0 string) *bookmarks.Label {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLabelByName", arg0)
	ret0, _ := ret[0].(*bookmarks.Label)
	return ret0
}

// GetLabelByName indicates an expected call of GetLabelByName
func (mr *MockILabelsMockRecorder) GetLabelByName(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLabelByName", reflect.TypeOf((*MockILabels)(nil).GetLabelByName), arg0)
}

// GetNameFromID mocks base method
func (m *MockILabels) GetNameFromID(arg0 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNameFromID", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNameFromID indicates an expected call of GetNameFromID
func (mr *MockILabelsMockRecorder) GetNameFromID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNameFromID", reflect.TypeOf((*MockILabels)(nil).GetNameFromID), arg0)
}

// StoreLabels mocks base method
func (m *MockILabels) StoreLabels() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreLabels")
	ret0, _ := ret[0].(error)
	return ret0
}

// StoreLabels indicates an expected call of StoreLabels
func (mr *MockILabelsMockRecorder) StoreLabels() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreLabels", reflect.TypeOf((*MockILabels)(nil).StoreLabels))
}

// UpdateLabel mocks base method
func (m *MockILabels) UpdateLabel(arg0 string, arg1 *bookmarks.Label) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateLabel", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateLabel indicates an expected call of UpdateLabel
func (mr *MockILabelsMockRecorder) UpdateLabel(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateLabel", reflect.TypeOf((*MockILabels)(nil).UpdateLabel), arg0, arg1)
}