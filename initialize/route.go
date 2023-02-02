package initialize

import (
	"douyin/base"
	"douyin/service"
	"github.com/gin-gonic/gin"
)

// InitRouter 初始化 Gin 的路由
func InitRouter(router gin.IRouter) {
	router.GET("/", base.HandlerFuncConverter(service.UserInfo))

}
