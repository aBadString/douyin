package main

import (
	"douyin/initialize"
	"github.com/gin-gonic/gin"
)

func main() {
	var ginEngine *gin.Engine = gin.Default()
	initialize.InitRouter(ginEngine)
	_ = ginEngine.Run()
}
