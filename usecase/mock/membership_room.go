// Code generated by MockGen. DO NOT EDIT.
// Source: membership_room.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockMembershipRoomUseCase is a mock of MembershipRoomUseCase interface.
type MockMembershipRoomUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockMembershipRoomUseCaseMockRecorder
}

// MockMembershipRoomUseCaseMockRecorder is the mock recorder for MockMembershipRoomUseCase.
type MockMembershipRoomUseCaseMockRecorder struct {
	mock *MockMembershipRoomUseCase
}

// NewMockMembershipRoomUseCase creates a new mock instance.
func NewMockMembershipRoomUseCase(ctrl *gomock.Controller) *MockMembershipRoomUseCase {
	mock := &MockMembershipRoomUseCase{ctrl: ctrl}
	mock.recorder = &MockMembershipRoomUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMembershipRoomUseCase) EXPECT() *MockMembershipRoomUseCaseMockRecorder {
	return m.recorder
}

// CreateMembershipRoom mocks base method.
func (m *MockMembershipRoomUseCase) CreateMembershipRoom(ctx context.Context, membershipID, roomID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateMembershipRoom", ctx, membershipID, roomID)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateMembershipRoom indicates an expected call of CreateMembershipRoom.
func (mr *MockMembershipRoomUseCaseMockRecorder) CreateMembershipRoom(ctx, membershipID, roomID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMembershipRoom", reflect.TypeOf((*MockMembershipRoomUseCase)(nil).CreateMembershipRoom), ctx, membershipID, roomID)
}

// DeleteMembershipRoom mocks base method.
func (m *MockMembershipRoomUseCase) DeleteMembershipRoom(ctx context.Context, membershipID, roomID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteMembershipRoom", ctx, membershipID, roomID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteMembershipRoom indicates an expected call of DeleteMembershipRoom.
func (mr *MockMembershipRoomUseCaseMockRecorder) DeleteMembershipRoom(ctx, membershipID, roomID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteMembershipRoom", reflect.TypeOf((*MockMembershipRoomUseCase)(nil).DeleteMembershipRoom), ctx, membershipID, roomID)
}
