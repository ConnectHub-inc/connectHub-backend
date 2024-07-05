//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package usecase

import (
	"context"
	"fmt"

	"github.com/tusmasoma/connectHub-backend/entity"
	"github.com/tusmasoma/connectHub-backend/internal/log"
	"github.com/tusmasoma/connectHub-backend/repository"
)

type MembershipUseCase interface {
	ListMemberships(ctx context.Context, workspaceID string) ([]entity.Membership, error)
	ListRoomMemberships(ctx context.Context, channelID string) ([]entity.Membership, error)
	GetMembership(ctx context.Context, membershipID string) (*entity.Membership, error)
	UpdateMembership(ctx context.Context, params *UpdateMembershipParams, membership entity.Membership) error
}

type membershipUseCase struct {
	mr repository.MembershipRepository
}

func NewMembershipUseCase(mr repository.MembershipRepository) MembershipUseCase {
	return &membershipUseCase{
		mr: mr,
	}
}

func (muc *membershipUseCase) ListMemberships(ctx context.Context, workspaceID string) ([]entity.Membership, error) {
	memberships, err := muc.mr.List(ctx, []repository.QueryCondition{{Field: "WorkspaceID", Value: workspaceID}})
	if err != nil {
		log.Error("Failed to list memberships", log.Fstring("workspaceID", workspaceID))
		return nil, err
	}
	return memberships, nil
}

func (muc *membershipUseCase) ListRoomMemberships(ctx context.Context, channelID string) ([]entity.Membership, error) {
	memberships, err := muc.mr.ListRoomMemberships(ctx, channelID)
	if err != nil {
		log.Error("Failed to list room memberships", log.Fstring("channelID", channelID))
		return nil, err
	}
	return memberships, nil
}

func (muc *membershipUseCase) GetMembership(ctx context.Context, membershipID string) (*entity.Membership, error) {
	membership, err := muc.mr.Get(ctx, membershipID)
	if err != nil {
		log.Error("Failed to get membership", log.Fstring("membershipID", membershipID))
		return nil, err
	}
	return membership, nil
}

type UpdateMembershipParams struct {
	UserID          string `json:"userID"`
	Name            string `json:"name"`
	ProfileImageURL string `json:"profile_image_url"`
}

func (muc *membershipUseCase) UpdateMembership(ctx context.Context, params *UpdateMembershipParams, membership entity.Membership) error {
	if membership.UserID != params.UserID {
		log.Warn("User don't have permission to update membership", log.Fstring("userID", membership.UserID))
		return fmt.Errorf("don't have permission to update membership")
	}

	membership.Name = params.Name
	membership.ProfileImageURL = params.ProfileImageURL

	if err := muc.mr.Update(ctx, membership); err != nil {
		log.Error(
			"Failed to update membership")
		return err
	}
	return nil
}
