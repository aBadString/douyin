package initialize

import (
	"douyin/base"
	"douyin/service"
	"github.com/gin-gonic/gin"
)

// InitRouter 初始化 Gin 的路由
func InitRouter(router gin.IRouter) {

	router.Static("/static", "./public")

	apiRouter := router.Group("/douyin")
	{
		//apiRouter.POST("/user/register/", service.Register)
		apiRouter.POST("/user/login/", service.Login)

		authRouter := apiRouter.Group("/", Auth)
		{
			// 基础接口
			apiRouter.GET("/feed/", base.HandlerFuncConverter(service.Feed))
			authRouter.GET("/user/", base.HandlerFuncConverter(service.UserInfo))
			authRouter.POST("/publish/action/", base.HandlerFuncConverter(service.Publish))
			authRouter.GET("/publish/list/", base.HandlerFuncConverter(service.PublishList))

			// 互动接口
			//authRouter.POST("/favorite/action/",  base.HandlerFuncConverter(service.FavoriteAction))
			authRouter.GET("/favorite/list/", base.HandlerFuncConverter(service.FavoriteList))
			//authRouter.POST("/comment/action/",  base.HandlerFuncConverter(service.CommentAction))
			//authRouter.GET("/comment/list/",  base.HandlerFuncConverter(service.CommentList))

			// 社交接口
			authRouter.POST("/relation/action/", base.HandlerFuncConverter(service.RelationAction))
			authRouter.GET("/relation/follow/list/", base.HandlerFuncConverter(service.FollowList))
			authRouter.GET("/relation/follower/list/", base.HandlerFuncConverter(service.FollowerList))
			authRouter.GET("/relation/friend/list/", base.HandlerFuncConverter(service.FriendList))
			//authRouter.GET("/message/chat/",  base.HandlerFuncConverter(service.MessageChat))
			//authRouter.POST("/message/action/",  base.HandlerFuncConverter(service.MessageAction))
		}
	}
}

func Auth(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		token = c.PostForm("token")
	}
	if token == "" {
		return
	}

	currentUserId, err := service.GetUserIdByToken(token)
	if err != nil {
		return
	}

	// 将 current_user_id 放到上下文中
	c.Set("current_user_id", currentUserId)
}
