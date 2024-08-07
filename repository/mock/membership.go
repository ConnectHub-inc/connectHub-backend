// Code generated by MockGen. DO NOT EDIT.
// Source: membership.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	entity "github.com/tusmasoma/connectHub-backend/entity"
	repository "github.com/tusmasoma/connectHub-backend/repository"
)

// MockMembershipRepository is a mock of MembershipRepository interface.
type MockMembershipRepository struct {
	ctrl     *gomock.Controller
	recorder *MockMembershipRepositoryMockRecorder
}

// MockMembershipRepositoryMockRecorder is the mock recorder for MockMembershipRepository.
type MockMembershipRepositoryMockRecorder struct {
	mock *MockMembershipRepository
}

// NewMockMembershipRepository creates a new mock instance.
func NewMockMembershipRepository(ctrl *gomock.Controller) *MockMembershipRepository {
	mock := &MockMembershipRepository{ctrl: ctrl}
	mock.recorder = &MockMembershipRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMembershipRepository) EXPECT() *MockMembershipRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockMembershipRepository) Create(ctx context.Context, membership entity.Membership) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, membership)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockMembershipRepositoryMockRecorder) Create(ctx, membership interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockMembershipRepository)(nil).Create), ctx, membership)
}

// Delete mocks base method.
func (m *MockMembershipRepository) Delete(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockMembershipRepositoryMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockMembershipRepository)(nil).Delete), ctx, id)
}

// Get mocks base method.
func (m *MockMembershipRepository) Get(ctx context.Context, id string) (*entity.Membership, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id)
	ret0, _ := ret[0].(*entity.Membership)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockMembershipRepositoryMockRecorder) Get(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockMembershipRepository)(nil).Get), ctx, id)
}

// List mocks base method.
func (m *MockMembershipRepository) List(ctx context.Context, qcs []repository.QueryCondition) ([]entity.Membership, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, qcs)
	ret0, _ := ret[0].([]entity.Membership)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockMembershipRepositoryMockRecorder) List(ctx, qcs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockMembershipRepository)(nil).List), ctx, qcs)
}

// ListChannelMemberships mocks base method.
func (m *MockMembershipRepository) ListChannelMemberships(ctx context.Context, channelID string) ([]entity.Membership, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListChannelMemberships", ctx, channelID)
	ret0, _ := ret[0].([]entity.Membership)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListChannelMemberships indicates an expected call of ListChannelMemberships.
func (mr *MockMembershipRepositoryMockRecorder) ListChannelMemberships(ctx, channelID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListChannelMemberships", reflect.TypeOf((*MockMembershipRepository)(nil).ListChannelMemberships), ctx, channelID)
}

// SoftDelete mocks base method.
func (m *MockMembershipRepository) SoftDelete(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SoftDelete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// SoftDelete indicates an expected call of SoftDelete.
func (mr *MockMembershipRepositoryMockRecorder) SoftDelete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SoftDelete", reflect.TypeOf((*MockMembershipRepository)(nil).SoftDelete), ctx, id)
}

// Update mocks base method.
func (m *MockMembershipRepository) Update(ctx context.Context, membership entity.Membership) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, membership)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockMembershipRepositoryMockRecorder) Update(ctx, membership interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockMembershipRepository)(nil).Update), ctx, membership)
}
