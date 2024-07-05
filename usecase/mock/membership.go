// Code generated by MockGen. DO NOT EDIT.
// Source: membership.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	entity "github.com/tusmasoma/connectHub-backend/entity"
	usecase "github.com/tusmasoma/connectHub-backend/usecase"
)

// MockMembershipUseCase is a mock of MembershipUseCase interface.
type MockMembershipUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockMembershipUseCaseMockRecorder
}

// MockMembershipUseCaseMockRecorder is the mock recorder for MockMembershipUseCase.
type MockMembershipUseCaseMockRecorder struct {
	mock *MockMembershipUseCase
}

// NewMockMembershipUseCase creates a new mock instance.
func NewMockMembershipUseCase(ctrl *gomock.Controller) *MockMembershipUseCase {
	mock := &MockMembershipUseCase{ctrl: ctrl}
	mock.recorder = &MockMembershipUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMembershipUseCase) EXPECT() *MockMembershipUseCaseMockRecorder {
	return m.recorder
}

// GetMembership mocks base method.
func (m *MockMembershipUseCase) GetMembership(ctx context.Context, membershipID string) (*entity.Membership, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMembership", ctx, membershipID)
	ret0, _ := ret[0].(*entity.Membership)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMembership indicates an expected call of GetMembership.
func (mr *MockMembershipUseCaseMockRecorder) GetMembership(ctx, membershipID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMembership", reflect.TypeOf((*MockMembershipUseCase)(nil).GetMembership), ctx, membershipID)
}

// ListMemberships mocks base method.
func (m *MockMembershipUseCase) ListMemberships(ctx context.Context, workspaceID string) ([]entity.Membership, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListMemberships", ctx, workspaceID)
	ret0, _ := ret[0].([]entity.Membership)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListMemberships indicates an expected call of ListMemberships.
func (mr *MockMembershipUseCaseMockRecorder) ListMemberships(ctx, workspaceID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListMemberships", reflect.TypeOf((*MockMembershipUseCase)(nil).ListMemberships), ctx, workspaceID)
}

// ListRoomMemberships mocks base method.
func (m *MockMembershipUseCase) ListRoomMemberships(ctx context.Context, channelID string) ([]entity.Membership, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListRoomMemberships", ctx, channelID)
	ret0, _ := ret[0].([]entity.Membership)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListRoomMemberships indicates an expected call of ListRoomMemberships.
func (mr *MockMembershipUseCaseMockRecorder) ListRoomMemberships(ctx, channelID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListRoomMemberships", reflect.TypeOf((*MockMembershipUseCase)(nil).ListRoomMemberships), ctx, channelID)
}

// UpdateMembership mocks base method.
func (m *MockMembershipUseCase) UpdateMembership(ctx context.Context, params *usecase.UpdateMembershipParams, membership entity.Membership) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMembership", ctx, params, membership)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateMembership indicates an expected call of UpdateMembership.
func (mr *MockMembershipUseCaseMockRecorder) UpdateMembership(ctx, params, membership interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMembership", reflect.TypeOf((*MockMembershipUseCase)(nil).UpdateMembership), ctx, params, membership)
}
