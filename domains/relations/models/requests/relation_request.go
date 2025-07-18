package requests

import "github.com/google/uuid"

type FollowRequest struct {
	FollowerID  uuid.UUID `json:"follower_id" binding:"required"`
	FollowingID uuid.UUID `json:"following_id" binding:"required"`
}

type UnfollowRequest struct {
	FollowerID  uuid.UUID `json:"follower_id" binding:"required"`
	FollowingID uuid.UUID `json:"following_id" binding:"required"`
}

type IsFollowingRequest struct {
	FollowerID  uuid.UUID `json:"follower_id" binding:"required"`
	FollowingID uuid.UUID `json:"following_id" binding:"required"`
}

type GetFollowersRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
}

type GetFollowingsRequest struct {
	UserID uuid.UUID `json:"user_id" binding:"required"`
}

