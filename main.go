package main

import (
	"douyin/conf"
	"douyin/initialize"
	"douyin/repository"
)

func main() {
	conf.Properties = initialize.InitApplicationProperties("app.json")
	repository.ORM = initialize.InitORM("visitor:visitor@tcp(localhost:3306)/douyin?charset=utf8&parseTime=True&loc=Local")
	ginEngine := initialize.InitGin()
	_ = ginEngine.Run()
}
