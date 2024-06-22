// Code generated by MockGen. DO NOT EDIT.
// Source: room.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	entity "github.com/tusmasoma/connectHub-backend/entity"
	repository "github.com/tusmasoma/connectHub-backend/repository"
)

// MockRoomRepository is a mock of RoomRepository interface.
type MockRoomRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRoomRepositoryMockRecorder
}

// MockRoomRepositoryMockRecorder is the mock recorder for MockRoomRepository.
type MockRoomRepositoryMockRecorder struct {
	mock *MockRoomRepository
}

// NewMockRoomRepository creates a new mock instance.
func NewMockRoomRepository(ctrl *gomock.Controller) *MockRoomRepository {
	mock := &MockRoomRepository{ctrl: ctrl}
	mock.recorder = &MockRoomRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRoomRepository) EXPECT() *MockRoomRepositoryMockRecorder {
	return m.recorder
}

// BatchCreate mocks base method.
func (m *MockRoomRepository) BatchCreate(ctx context.Context, rooms []entity.Room) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BatchCreate", ctx, rooms)
	ret0, _ := ret[0].(error)
	return ret0
}

// BatchCreate indicates an expected call of BatchCreate.
func (mr *MockRoomRepositoryMockRecorder) BatchCreate(ctx, rooms interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BatchCreate", reflect.TypeOf((*MockRoomRepository)(nil).BatchCreate), ctx, rooms)
}

// Create mocks base method.
func (m *MockRoomRepository) Create(ctx context.Context, room entity.Room) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, room)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockRoomRepositoryMockRecorder) Create(ctx, room interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRoomRepository)(nil).Create), ctx, room)
}

// CreateOrUpdate mocks base method.
func (m *MockRoomRepository) CreateOrUpdate(ctx context.Context, id string, qcs []repository.QueryCondition, room entity.Room) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrUpdate", ctx, id, qcs, room)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateOrUpdate indicates an expected call of CreateOrUpdate.
func (mr *MockRoomRepositoryMockRecorder) CreateOrUpdate(ctx, id, qcs, room interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrUpdate", reflect.TypeOf((*MockRoomRepository)(nil).CreateOrUpdate), ctx, id, qcs, room)
}

// Delete mocks base method.
func (m *MockRoomRepository) Delete(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockRoomRepositoryMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRoomRepository)(nil).Delete), ctx, id)
}

// Get mocks base method.
func (m *MockRoomRepository) Get(ctx context.Context, id string) (*entity.Room, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id)
	ret0, _ := ret[0].(*entity.Room)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockRoomRepositoryMockRecorder) Get(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockRoomRepository)(nil).Get), ctx, id)
}

// List mocks base method.
func (m *MockRoomRepository) List(ctx context.Context, qcs []repository.QueryCondition) ([]entity.Room, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, qcs)
	ret0, _ := ret[0].([]entity.Room)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockRoomRepositoryMockRecorder) List(ctx, qcs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockRoomRepository)(nil).List), ctx, qcs)
}

// ListUserWorkspaceRooms mocks base method.
func (m *MockRoomRepository) ListUserWorkspaceRooms(ctx context.Context, userID, workspaceID string) ([]entity.Room, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListUserWorkspaceRooms", ctx, userID, workspaceID)
	ret0, _ := ret[0].([]entity.Room)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListUserWorkspaceRooms indicates an expected call of ListUserWorkspaceRooms.
func (mr *MockRoomRepositoryMockRecorder) ListUserWorkspaceRooms(ctx, userID, workspaceID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListUserWorkspaceRooms", reflect.TypeOf((*MockRoomRepository)(nil).ListUserWorkspaceRooms), ctx, userID, workspaceID)
}

// Update mocks base method.
func (m *MockRoomRepository) Update(ctx context.Context, id string, room entity.Room) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, id, room)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockRoomRepositoryMockRecorder) Update(ctx, id, room interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockRoomRepository)(nil).Update), ctx, id, room)
}
