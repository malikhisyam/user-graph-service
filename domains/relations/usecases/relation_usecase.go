package usecases

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/malikhisyam/user-graph-service/domains/relations/entities"
	"github.com/malikhisyam/user-graph-service/domains/relations/repositories"
)

var (
	ErrCannotFollowSelf   = errors.New("cannot follow yourself")
	ErrCannotUnfollowSelf = errors.New("cannot unfollow yourself")
)

type RelationUseCase interface {
	Follow(ctx context.Context, followerID, followingID uuid.UUID) error
	Unfollow(ctx context.Context, followerID, followingID uuid.UUID) error
	IsFollowing(ctx context.Context, followerID, followingID uuid.UUID) (bool, error)
	GetFollowers(ctx context.Context, userID string) ([]entities.Follows, error)
	GetFollowings(ctx context.Context, userID string) ([]entities.Follows, error)
}

type relationUsecase struct {
	relationRepo repositories.RelationRepository
}

func NewRelationUseCase(relationRepo repositories.RelationRepository) RelationUseCase {
	return &relationUsecase{
		relationRepo: relationRepo,
	}
}

func (u *relationUsecase) Follow(ctx context.Context, followerID, followingID uuid.UUID) error {
	if followerID == followingID {
		return ErrCannotFollowSelf
	}

	return u.relationRepo.Follow(ctx, followerID, followingID)
}

func (u *relationUsecase) Unfollow(ctx context.Context, followerID, followingID uuid.UUID) error {
	if followerID == followingID {
		return ErrCannotUnfollowSelf
	}

	return u.relationRepo.Unfollow(ctx, followerID, followingID)
}

func (u *relationUsecase) IsFollowing(ctx context.Context, followerID, followingID uuid.UUID) (bool, error) {
	return u.relationRepo.IsFollowing(ctx, followerID, followingID)
}

func (u *relationUsecase) GetFollowers(ctx context.Context, userID string) ([]entities.Follows, error) {
	return u.relationRepo.GetFollowers(ctx, userID)
}

func (u *relationUsecase) GetFollowings(ctx context.Context, userID string) ([]entities.Follows, error) {
	return u.relationRepo.GetFollowings(ctx, userID)
}


