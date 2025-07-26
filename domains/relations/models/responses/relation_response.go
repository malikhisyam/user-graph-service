package responses

import (
	"time"
)

type IsFollowingResponse struct {
	IsFollowing bool `json:"is_following"`
}

type FollowerResponse struct {
    ID          string    `json:"id"`
    FollowerID  string    `json:"follower_id"`
    DisplayName string    `json:"display_name"`
    Username    string    `json:"username"`
    CreatedAt   time.Time `json:"created_at"`
}

type FollowingResponse struct {
	ID           string    `json:"id"`
	FollowerID   string    `json:"follower_id"`
	FollowingID  string    `json:"following_id"`
	Name string    `json:"name"`
	Username     string    `json:"username"`
	CreatedAt    time.Time `json:"created_at"`
}
type GetFollowersResponse struct {
    Followers []FollowerResponse `json:"followers"`
}

type GetFollowingsResponse struct {
	Followings []FollowingResponse `json:"followings"`
}
