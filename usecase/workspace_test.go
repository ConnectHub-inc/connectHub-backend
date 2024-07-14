package usecase

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/tusmasoma/connectHub-backend/repository/mock"
)

func TestWorkspaceUseCase_CreateWorkspace(t *testing.T) {
	t.Parallel()

	patterns := []struct {
		name  string
		setup func(
			m *mock.MockWorkspaceRepository,
		)
		arg struct {
			ctx  context.Context
			name string
		}
		wantErr error
	}{
		{
			name: "success",
			setup: func(m *mock.MockWorkspaceRepository) {
				m.EXPECT().Create(
					gomock.Any(),
					gomock.Any(),
				).Return(nil)
			},
			arg: struct {
				ctx  context.Context
				name string
			}{
				ctx:  context.Background(),
				name: "test",
			},
			wantErr: nil,
		},
		{
			name: "Fail: name is required",
			arg: struct {
				ctx  context.Context
				name string
			}{
				ctx:  context.Background(),
				name: "",
			},
			wantErr: fmt.Errorf("name is required"),
		},
	}
	for _, tt := range patterns {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			wr := mock.NewMockWorkspaceRepository(ctrl)

			if tt.setup != nil {
				tt.setup(wr)
			}

			usecase := NewWorkspaceUseCase(wr)
			err := usecase.CreateWorkspace(tt.arg.ctx, tt.arg.name)

			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("CreateWorkspace() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("CreateWorkspace() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
