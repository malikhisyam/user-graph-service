package repositories

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/malikhisyam/user-graph-service/domains/relations/entities"
	"github.com/malikhisyam/user-graph-service/domains/relations/models/responses"
	"github.com/malikhisyam/user-graph-service/infrastructures"
	"github.com/malikhisyam/user-graph-service/shared/util"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RelationRepository interface {
	Follow(ctx context.Context, followerID, followingID uuid.UUID) error
	Unfollow(ctx context.Context, followerID, followingID uuid.UUID) error
	IsFollowing(ctx context.Context, followerID, followingID uuid.UUID) (bool, error)
	GetFollowers(ctx context.Context, userID string, limit, offset int, nameFilter string) ([]responses.FollowerWithUserInfo, error)
	GetFollowings(ctx context.Context, userID string, limit, offset int, nameFilter string) ([]responses.FollowingWithUserInfo, error)
}

type relationRepository struct {
	db infrastructures.Database
	redisCache *redis.Client
	logger util.Logger
}

func NewRelationRepository(db infrastructures.Database, redisClient *redis.Client, logger util.Logger) RelationRepository {
	return &relationRepository{
		db: db,
		redisCache: redisClient,
		logger: logger,
	}
}

func followKey(followerID, followingID uuid.UUID) string {
	return fmt.Sprintf("follow:%s:%s", followerID.String(), followingID.String())
}


func (r *relationRepository) Follow(ctx context.Context, followerID, followingID uuid.UUID) error {
	var existing entities.Follows

	err := r.db.GetInstance().
		WithContext(ctx).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		First(&existing).Error

	if err == nil {
		r.logger.Warn("Attempted to create a follow relationship that already exists",
			zap.String("follower_id", followerID.String()),
			zap.String("following_id", followingID.String()),
		)
		return fmt.Errorf("user already following")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		r.logger.Error("Failed to check for existing follow relationship",
			zap.Error(err),
			zap.String("follower_id", followerID.String()),
			zap.String("following_id", followingID.String()),
		)
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
		r.logger.Error("Failed to create follow relationship in database",
			zap.Error(result.Error),
			zap.String("follower_id", followerID.String()),
			zap.String("following_id", followingID.String()),
		)
		return result.Error
	}

	cacheKey := fmt.Sprintf("follow:%s:%s", followerID.String(), followingID.String()) // Assuming a key generation function
	if err := r.redisCache.Set(ctx, cacheKey, "1", 10*time.Minute).Err(); err != nil {
		r.logger.Error("Failed to set follow relationship in Redis cache",
			zap.Error(err),
			zap.String("cache_key", cacheKey),
		)
	}

	r.logger.Info("User followed successfully",
		zap.String("follower_id", followerID.String()),
		zap.String("following_id", followingID.String()),
	)
	return nil
}


func (r *relationRepository) Unfollow(ctx context.Context, followerID, followingID uuid.UUID) error {
	r.logger.Info("Attempting to unfollow user",
		zap.String("follower_id", followerID.String()),
		zap.String("following_id", followingID.String()),
	)

	result := r.db.GetInstance().
		Unscoped().
		WithContext(ctx).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Delete(&entities.Follows{})

	if result.Error != nil {
		r.logger.Error("Failed to delete follow relationship from database",
			zap.Error(result.Error),
			zap.String("follower_id", followerID.String()),
			zap.String("following_id", followingID.String()),
		)
		return result.Error
	}

	if result.RowsAffected == 0 {
		r.logger.Warn("Unfollow attempt on a non-existent relationship",
			zap.String("follower_id", followerID.String()),
			zap.String("following_id", followingID.String()),
		)
		return fmt.Errorf("follow relationship not found")
	}

	cacheKey := followKey(followerID, followingID)
	if err := r.redisCache.Del(ctx, cacheKey).Err(); err != nil {
		r.logger.Error("Failed to delete follow relationship from Redis cache",
			zap.Error(err),
			zap.String("cache_key", cacheKey),
		)
	}

	r.logger.Info("User unfollowed successfully",
		zap.String("follower_id", followerID.String()),
		zap.String("following_id", followingID.String()),
	)

	return nil
}


func (r *relationRepository) IsFollowing(ctx context.Context, followerID, followingID uuid.UUID) (bool, error) {
	cacheKey := followKey(followerID, followingID)
	cached, err := r.redisCache.Get(ctx, cacheKey).Result()

	if err == nil {
		r.logger.Debug("Cache hit for IsFollowing check",
			zap.String("cache_key", cacheKey),
			zap.String("result", cached),
		)
		return cached == "1", nil
	}

	if err != redis.Nil {
		r.logger.Error("Redis error during IsFollowing check",
			zap.Error(err),
			zap.String("cache_key", cacheKey),
		)
	} else {
		r.logger.Debug("Cache miss for IsFollowing check", zap.String("cache_key", cacheKey))
	}

	var follow entities.Follows
	dbErr := r.db.GetInstance().
		WithContext(ctx).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		First(&follow).Error

	if dbErr != nil {
		if errors.Is(dbErr, gorm.ErrRecordNotFound) {
			if cacheErr := r.redisCache.Set(ctx, cacheKey, "0", 10*time.Minute).Err(); cacheErr != nil {
				r.logger.Error("Failed to set 'not following' status in cache",
					zap.Error(cacheErr),
					zap.String("cache_key", cacheKey),
				)
			}
			return false, nil
		}

		r.logger.Error("Database error during IsFollowing check",
			zap.Error(dbErr),
			zap.String("follower_id", followerID.String()),
			zap.String("following_id", followingID.String()),
		)
		return false, dbErr
	}

	if cacheErr := r.redisCache.Set(ctx, cacheKey, "1", 10*time.Minute).Err(); cacheErr != nil {
		r.logger.Error("Failed to set 'following' status in cache",
			zap.Error(cacheErr),
			zap.String("cache_key", cacheKey),
		)
	}
	return true, nil
}

func (r *relationRepository) GetFollowers(ctx context.Context, userID string, limit, offset int, nameFilter string) ([]responses.FollowerWithUserInfo, error) {
	var followers []responses.FollowerWithUserInfo

	db := r.db.GetInstance().WithContext(ctx).
		Table("follows").
		Select("follows.id, follows.follower_id, users.name, users.username").
		Joins("JOIN users ON follows.follower_id = users.id").
		Where("follows.following_id = ?", userID).
		Order("follows.created_at DESC").
		Limit(limit).
		Offset(offset)

	if nameFilter != "" {
		db = db.Where("LOWER(users.name) LIKE ?", "%"+strings.ToLower(nameFilter)+"%")
	}

	err := db.Find(&followers).Error
	return followers, err
}

func (r *relationRepository) GetFollowings(ctx context.Context, userID string, limit, offset int, nameFilter string) ([]responses.FollowingWithUserInfo, error) {
	var results []responses.FollowingWithUserInfo

	db := r.db.GetInstance().WithContext(ctx)

	query := db.
		Table("follows AS f").
		Select(`
			f.id,
			f.follower_id,
			f.following_id,
			u.name,
			u.username
		`).
		Joins("JOIN users u ON f.following_id = u.id").
		Where("f.follower_id = ?", userID)

	if nameFilter != "" {
		query = query.Where("(u.name ILIKE ? OR u.username ILIKE ?)", "%"+nameFilter+"%", "%"+nameFilter+"%")
	}

	err := query.
		Order("f.created_at DESC").
		Limit(limit).
		Offset(offset).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}

