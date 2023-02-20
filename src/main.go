package main

import (
	"douyin/conf"
	"douyin/initialize"
	"douyin/redis"
	"douyin/repository"
	"os"
	"strconv"
)

func main() {
	configFile := "app.json"
	if len(os.Args) >= 2 {
		configFile = os.Args[1]
	}

	conf.Properties = initialize.InitApplicationProperties(configFile)
	repository.ORM = initialize.InitORM(conf.Properties.DatabaseUrl)
	ginEngine := initialize.InitGin()

	if conf.Properties.Redis.Enable {
		redisConf := conf.Properties.Redis
		redisClient, err := initialize.InitRedisClient(redisConf.Addr, redisConf.Password, redisConf.Db)
		if err == nil {
			redis.RedisClient = redisClient
		}
	}

	port := 8080
	if conf.Properties.Port != 0 {
		port = conf.Properties.Port
	}
	https := conf.Properties.Https
	if https.Enable {
		_ = ginEngine.RunTLS(":"+strconv.Itoa(port), https.CertFile, https.KeyFile)
	} else {
		_ = ginEngine.Run(":" + strconv.Itoa(port))
	}
}
