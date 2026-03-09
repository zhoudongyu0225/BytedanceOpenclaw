package zlog

import (
	"dmGameServer/common"
	"fmt"
	"github.com/rs/zerolog"
	"os"
)

var Logger zerolog.Logger

// InitLog 初始化日志
func InitLog() {
	//  0 debug  1 info
	if common.GetConfConfig().Log == "" {
		Logger = zerolog.New(os.Stdout).With().Timestamp().Caller().Logger().Level(zerolog.Level(common.GetConfConfig().LogLv))
	} else {
		logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println("open log file failed, err:", err)
		}
		Logger = zerolog.New(logFile).With().Timestamp().Caller().Logger().Level(zerolog.Level(common.GetConfConfig().LogLv))
	}

}
