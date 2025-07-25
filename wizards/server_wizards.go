package wizards

import (
	"github.com/gin-gonic/gin"
	// "github.com/malikhisyam/user-graph-service/shared/middlewares"
)

func RegisterServer(router *gin.Engine) {
	api := router.Group("/api")
	v1 := api.Group("/v1")
	relation := v1.Group("/relations")
	{
		// relation.Use(middlewares.AuthMiddleware())
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