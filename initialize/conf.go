package initialize

import (
	"douyin/conf"
	"encoding/json"
	"os"
)

func InitApplicationProperties(configFile string) *conf.ApplicationProperties {
	var properties = &conf.ApplicationProperties{
		DataPath:  "./public/",
		DataUrl:   "/static/",
		Hostname:  "http://localhost:8080",
		SecretKey: "default",
	}

	bytes, err := os.ReadFile(configFile)
	if err != nil {
		return properties
	}

	err = json.Unmarshal(bytes, properties)
	if err != nil {
		return properties
	}

	return properties
}
