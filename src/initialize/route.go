package initialize

import (
	"douyin/base"
	"douyin/conf"
	"douyin/service"
	"github.com/gin-gonic/gin"
)

func InitGin() *gin.Engine {
	var ginEngine = gin.Default()
	initRouter(ginEngine)
	return ginEngine
}

// initRouter 初始化 Gin 的路由
func initRouter(router gin.IRouter) {

	router.Static(conf.Properties.DataUrl, conf.Properties.DataPath)

	apiRouter := router.Group("/douyin")
	{
		apiRouter.POST("/user/register/", service.Register)
		apiRouter.POST("/user/login/", service.Login)

		authRouter := apiRouter.Group("/", Auth)
		{
			// 基础接口
			authRouter.GET("/feed/", base.HandlerFuncConverter(service.Feed))
			authRouter.GET("/user/", base.HandlerFuncConverter(service.UserInfo))
			authRouter.POST("/publish/action/", base.HandlerFuncConverter(service.PublishVideo))
			authRouter.GET("/publish/list/", base.HandlerFuncConverter(service.PublishList))

			// 互动接口
			authRouter.POST("/favorite/action/", base.HandlerFuncConverter(service.FavoriteAction))
			authRouter.GET("/favorite/list/", base.HandlerFuncConverter(service.FavoriteList))
			authRouter.POST("/comment/action/", base.HandlerFuncConverter(service.CommentAction))
			authRouter.GET("/comment/list/", base.HandlerFuncConverter(service.VideoCommentList))

			// 社交接口
			authRouter.POST("/relation/action/", base.HandlerFuncConverter(service.RelationAction))
			authRouter.GET("/relation/follow/list/", base.HandlerFuncConverter(service.FollowList))
			authRouter.GET("/relation/follower/list/", base.HandlerFuncConverter(service.FollowerList))
			authRouter.GET("/relation/friend/list/", base.HandlerFuncConverter(service.FriendList))
			authRouter.GET("/message/chat/", base.HandlerFuncConverter(service.MessageChat))
			authRouter.POST("/message/action/", base.HandlerFuncConverter(service.MessageAction))
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
