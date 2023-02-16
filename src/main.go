package main

import (
	"douyin/conf"
	"douyin/initialize"
	"douyin/repository"
)

func main() {
	conf.Properties = initialize.InitApplicationProperties("app.json")
	repository.ORM = initialize.InitORM(conf.Properties.DatabaseUrl)
	ginEngine := initialize.InitGin()
	_ = ginEngine.Run()
}
