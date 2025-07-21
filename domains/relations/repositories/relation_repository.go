package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/malikhisyam/user-graph-service/domains/relations/entities"
	"github.com/malikhisyam/user-graph-service/infrastructures"
	"gorm.io/gorm"
)

type RelationRepository interface {
	Follow(ctx context.Context, followerID, followingID uuid.UUID) error
	Unfollow(ctx context.Context, followerID, followingID uuid.UUID) error
	IsFollowing(ctx context.Context, followerID, followingID uuid.UUID) (bool, error)
	GetFollowers(ctx context.Context, userID uuid.UUID) ([]entities.Follows, error)
	GetFollowings(ctx context.Context, userID string) ([]entities.Follows, error)
}

type relationRepository struct {
	db infrastructures.Database
}

func NewRelationRepository(db infrastructures.Database) RelationRepository {
	return &relationRepository{
		db: db,
	}
}

func (r *relationRepository) Follow(ctx context.Context, followerID, followingID uuid.UUID) error {
	var existing entities.Follows

	err := r.db.GetInstance().
	WithContext(ctx).
	Where("follower_id = ? AND following_id = ?", followerID, followingID).
	First(&existing).Error

	if err == nil {
		return fmt.Errorf("user already following")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	follow := entities.Follows{
		ID:          uuid.New(),
		FollowerID:  followerID,
		FollowingID: followingID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	result := r.db.GetInstance().WithContext(ctx).Create(&follow)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *relationRepository) Unfollow(ctx context.Context, followerID, followingID uuid.UUID) error {
	result := r.db.GetInstance().
		WithContext(ctx).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Delete(&entities.Follows{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *relationRepository) IsFollowing(ctx context.Context, followerID, followingID uuid.UUID) (bool, error) {
	var follow entities.Follows

	err := r.db.GetInstance().
		WithContext(ctx).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		First(&follow).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (r *relationRepository) GetFollowers(ctx context.Context, userID uuid.UUID) ([]entities.Follows, error) {
	var followers []entities.Follows

	err := r.db.GetInstance().
		WithContext(ctx).
		Where("following_id = ?", userID).
		Find(&followers).Error

	if err != nil {
		return nil, err
	}

	return followers, nil
}

func (r *relationRepository) GetFollowings(ctx context.Context, userID string) ([]entities.Follows, error) {
	var followings []entities.Follows

	err := r.db.GetInstance().
		WithContext(ctx).
		Where("follower_id = ?", userID).
		Find(&followings).Error

	if err != nil {
		return nil, err
	}

	return followings, nil
}
