package responses

import (
	"time"

	"github.com/google/uuid"
)

type IsFollowingResponse struct {
	IsFollowing bool `json:"is_following"`
}

type FollowerResponse struct {
	ID          uuid.UUID `json:"id"`
	FollowerID  uuid.UUID `json:"follower_id"`
	FollowingID uuid.UUID `json:"following_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type FolllowingResponse struct {
	ID          uuid.UUID `json:"id"`
	FollowerID  uuid.UUID `json:"follower_id"`
	FollowingID uuid.UUID `json:"following_id"`
	CreatedAt   time.Time `json:"created_at"`
}
type GetFollowersResponse struct {
	Followers []FollowerResponse `json:"followers"`
}

type GetFollowingsResponse struct {
	Followings []FolllowingResponse`json:"followings"`
}
