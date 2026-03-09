package common

import (
	"fmt"
	"github.com/goccy/go-json"
	"os"
)

var confConfig *ConfConfig

// ConfConfig 配置
type ConfConfig struct {
	Log           string `json:"log"`
	Port          string `json:"port"`
	RedisUrl      string `json:"redisUrl"`
	RedisPassword string `json:"redisPassword"`
	MongoUrl      string `json:"mongoUrl"`
	LogLv         int    `json:"logLv"`  // 日志等级 0 debug 1info
	GameId        int    `json:"gameId"` // 游戏id
}

func InitConf(LogPath string) {
	// 打开配置文件
	configFile, err := os.Open(LogPath)
	if err != nil {
		fmt.Println("无法打开配置文件:", err, LogPath)
		return
	}
	defer configFile.Close()
	confConfig = &ConfConfig{}
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&confConfig)
	if err != nil {
		fmt.Println("无法解码配置文件::", err, LogPath)
		return
	}
	fmt.Println("confConfig", confConfig)
}

func GetConfConfig() *ConfConfig {
	return confConfig
}
