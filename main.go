package main

import (
	"douyin/initialize"
	"github.com/gin-gonic/gin"
)

func main() {
	var ginEngine *gin.Engine = gin.Default()
	initialize.InitRouter(ginEngine)
	initialize.InitORM("visitor:visitor@tcp(localhost:3306)/douyin?charset=utf8&parseTime=True&loc=Local")
	_ = ginEngine.Run()
}
