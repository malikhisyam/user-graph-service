package responses

import "time"

type FollowerWithUserInfo struct {
	ID         string `json:"id"`
	FollowerID string `json:"follower_id"`
	Name       string `json:"name"`
	Username   string `json:"username"`
}

type FollowingWithUserInfo struct {
	ID          string
	FollowerID  string
	FollowingID string
	Name string
	Username    string
	AvatarURL   string
	CreatedAt   time.Time
}
