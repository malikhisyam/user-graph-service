package http

import (
	"net/http"
	"strconv"

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
	userId := c.Param("userId")

	// Pagination
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	offset := (page - 1) * limit

	nameFilter := c.DefaultQuery("name", "")

	followers, err := h.relationUc.GetFollowers(c.Request.Context(), userId, limit, offset, nameFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var followerResponses []responses.FollowerResponse
	for _, f := range followers {
		followerResponses = append(followerResponses, responses.FollowerResponse{
			ID:          f.ID,
			FollowerID:  f.FollowerID,
			DisplayName: f.Name,
			Username:    f.Username,
		})
	}

	resp := responses.GetFollowersResponse{
		Followers: followerResponses,
	}

	c.JSON(http.StatusOK, resp)
}



func (h *RelationHttp) GetFollowings(c *gin.Context) {
	userId := c.Param("userId")

	// Pagination
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	offset := (page - 1) * limit

	// Name filter
	nameFilter := c.DefaultQuery("name", "")

	// Usecase
	followings, err := h.relationUc.GetFollowings(c.Request.Context(), userId, limit, offset, nameFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var followingResponses []responses.FollowingResponse
	for _, f := range followings {
		followingResponses = append(followingResponses, responses.FollowingResponse{
			ID:           f.ID,
			FollowerID:   f.FollowerID,
			FollowingID:  f.FollowingID,
			Name:  f.Name,
			Username:     f.Username,
			CreatedAt:    f.CreatedAt,
		})
	}

	resp := responses.GetFollowingsResponse{
		Followings: followingResponses,
	}

	c.JSON(http.StatusOK, resp)
}


