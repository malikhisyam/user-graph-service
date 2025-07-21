package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/malikhisyam/user-graph-service/domains/relations/models/requests"
	"github.com/malikhisyam/user-graph-service/domains/relations/models/responses"
	"github.com/malikhisyam/user-graph-service/domains/relations/usecases"
)

type RelationHttp struct {
	relationUc usecases.RelationUseCase
}

func NewRelationHttp(relationUc usecases.RelationUseCase) *RelationHttp {
	return &RelationHttp{
		relationUc: relationUc,
	}
}

func (h *RelationHttp) Follow(c *gin.Context) {
	var req requests.FollowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.relationUc.Follow(c.Request.Context(), req.FollowerID, req.FollowingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "follow success"})
}

func (h *RelationHttp) Unfollow(c *gin.Context) {
	var req requests.UnfollowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.relationUc.Unfollow(c.Request.Context(), req.FollowerID, req.FollowingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "unfollow success"})
}

func (h *RelationHttp) IsFollowing(c *gin.Context) {
	var req requests.IsFollowingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isFollowing, err := h.relationUc.IsFollowing(c.Request.Context(), req.FollowerID, req.FollowingID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := responses.IsFollowingResponse{
		IsFollowing: isFollowing,
	}
	c.JSON(http.StatusOK, resp)
}

func (h *RelationHttp) GetFollowers(c *gin.Context) {
	var req requests.GetFollowersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"gg": err.Error()})
		return
	}

	follows, err := h.relationUc.GetFollowers(c.Request.Context(), req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var followerResponses []responses.FollowerResponse
	for _, f := range follows {
		followerResponses = append(followerResponses, responses.FollowerResponse{
			ID:          f.ID,
			FollowerID:  f.FollowerID,
			FollowingID: f.FollowingID,
			CreatedAt:   f.CreatedAt,
		})
	}

	resp := responses.GetFollowersResponse{
		Followers: followerResponses,
	}

	c.JSON(http.StatusOK, resp)
}

func (h *RelationHttp) GetFollowings(c *gin.Context) {
	userId := c.Param("userId")

	followings, err := h.relationUc.GetFollowings(c.Request.Context(), userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var followingResponses []responses.FolllowingResponse
	for _, f := range followings {
		followingResponses = append(followingResponses, responses.FolllowingResponse{
			ID:          f.ID,
			FollowerID:  f.FollowerID,
			FollowingID: f.FollowingID,
			CreatedAt:   f.CreatedAt,
		})
	}

	resp := responses.GetFollowingsResponse{
		Followings: followingResponses,
	}

	c.JSON(http.StatusOK, resp)
}

