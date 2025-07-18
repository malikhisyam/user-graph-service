package wizards

import (
	"github.com/gin-gonic/gin"
)

func RegisterServer(router *gin.Engine) {
	relation := router.Group("/relations")
	{
		// Follow User 
		relation.POST("/followings", RelationHttp.Follow)
		// Unfollow User
		relation.DELETE("/followings", RelationHttp.Unfollow)
		// See If A User Followed By A User
		relation.GET("/:userId/followings/:targetUserId", RelationHttp.IsFollowing)
		// Get Specific User His/Her Followers
		relation.GET("/:userId/followers", RelationHttp.GetFollowers)
		// Get Specific User His/Her Followings
		relation.GET("/:userId/followings", RelationHttp.GetFollowings)
	}
}