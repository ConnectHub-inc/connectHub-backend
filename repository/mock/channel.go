// Code generated by MockGen. DO NOT EDIT.
// Source: channel.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	entity "github.com/tusmasoma/connectHub-backend/entity"
	repository "github.com/tusmasoma/connectHub-backend/repository"
)

// MockChannelRepository is a mock of ChannelRepository interface.
type MockChannelRepository struct {
	ctrl     *gomock.Controller
	recorder *MockChannelRepositoryMockRecorder
}

// MockChannelRepositoryMockRecorder is the mock recorder for MockChannelRepository.
type MockChannelRepositoryMockRecorder struct {
	mock *MockChannelRepository
}

// NewMockChannelRepository creates a new mock instance.
func NewMockChannelRepository(ctrl *gomock.Controller) *MockChannelRepository {
	mock := &MockChannelRepository{ctrl: ctrl}
	mock.recorder = &MockChannelRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockChannelRepository) EXPECT() *MockChannelRepositoryMockRecorder {
	return m.recorder
}

// BatchCreate mocks base method.
func (m *MockChannelRepository) BatchCreate(ctx context.Context, channels []entity.Channel) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BatchCreate", ctx, channels)
	ret0, _ := ret[0].(error)
	return ret0
}

// BatchCreate indicates an expected call of BatchCreate.
func (mr *MockChannelRepositoryMockRecorder) BatchCreate(ctx, channels interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BatchCreate", reflect.TypeOf((*MockChannelRepository)(nil).BatchCreate), ctx, channels)
}

// Create mocks base method.
func (m *MockChannelRepository) Create(ctx context.Context, channel entity.Channel) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, channel)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockChannelRepositoryMockRecorder) Create(ctx, channel interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockChannelRepository)(nil).Create), ctx, channel)
}

// CreateOrUpdate mocks base method.
func (m *MockChannelRepository) CreateOrUpdate(ctx context.Context, id string, qcs []repository.QueryCondition, channel entity.Channel) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrUpdate", ctx, id, qcs, channel)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateOrUpdate indicates an expected call of CreateOrUpdate.
func (mr *MockChannelRepositoryMockRecorder) CreateOrUpdate(ctx, id, qcs, channel interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrUpdate", reflect.TypeOf((*MockChannelRepository)(nil).CreateOrUpdate), ctx, id, qcs, channel)
}

// Delete mocks base method.
func (m *MockChannelRepository) Delete(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockChannelRepositoryMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockChannelRepository)(nil).Delete), ctx, id)
}

// Get mocks base method.
func (m *MockChannelRepository) Get(ctx context.Context, id string) (*entity.Channel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id)
	ret0, _ := ret[0].(*entity.Channel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockChannelRepositoryMockRecorder) Get(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockChannelRepository)(nil).Get), ctx, id)
}

// List mocks base method.
func (m *MockChannelRepository) List(ctx context.Context, qcs []repository.QueryCondition) ([]entity.Channel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, qcs)
	ret0, _ := ret[0].([]entity.Channel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockChannelRepositoryMockRecorder) List(ctx, qcs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockChannelRepository)(nil).List), ctx, qcs)
}

// ListMembershipChannels mocks base method.
func (m *MockChannelRepository) ListMembershipChannels(ctx context.Context, membershipID string) ([]entity.Channel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListMembershipChannels", ctx, membershipID)
	ret0, _ := ret[0].([]entity.Channel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListMembershipChannels indicates an expected call of ListMembershipChannels.
func (mr *MockChannelRepositoryMockRecorder) ListMembershipChannels(ctx, membershipID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListMembershipChannels", reflect.TypeOf((*MockChannelRepository)(nil).ListMembershipChannels), ctx, membershipID)
}

// Update mocks base method.
func (m *MockChannelRepository) Update(ctx context.Context, id string, channel entity.Channel) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, id, channel)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockChannelRepositoryMockRecorder) Update(ctx, id, channel interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockChannelRepository)(nil).Update), ctx, id, channel)
}
