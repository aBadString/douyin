package main

import (
	"douyin/conf"
	"douyin/initialize"
	"douyin/repository"
	"os"
)

func main() {
	configFile := "app.json"
	if len(os.Args) >= 2 {
		configFile = os.Args[1]
	}

	conf.Properties = initialize.InitApplicationProperties(configFile)
	repository.ORM = initialize.InitORM(conf.Properties.DatabaseUrl)
	ginEngine := initialize.InitGin()
	_ = ginEngine.Run()
}
