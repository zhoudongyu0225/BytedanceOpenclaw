package main

import (
	"dmGameServer/common"
	"dmGameServer/model"
	"dmGameServer/servers/logic"
	"dmGameServer/servers/logic/control"
	"dmGameServer/untils"
	"dmGameServer/zlog"
	"flag"
	"math/rand"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"syscall"
	"time"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			stack := debug.Stack()
			untils.PanicPoss(err, stack)
		}
	}()
	// 服务器启动配置
	path := flag.String("f", "", "run config")
	flag.Parse()
	// 日志文件
	LogPath := *path
	common.ServiceSign = strings.Replace(LogPath, "conf/", "", -1)
	common.ServiceSign = strings.Replace(common.ServiceSign, ".json", "", -1)
	common.ServiceType = strings.Split(common.ServiceSign, "_")[0]
	// 初始化程序配置
	common.InitConf(LogPath)
	// 初始化日志
	zlog.InitLog()
	// 初始化游戏id
	common.GameId = common.GetConfConfig().GameId
	// 设置随机时间戳
	rand.Seed(time.Now().UnixNano())
	// 初始化redis
	model.InitRedis()
	// 初始化mongodb
	model.InitMongoDB()
	zlog.Logger.Info().Msgf("服务器启动 ServiceSign%v ServiceType%v gameServer:%v", common.ServiceSign, common.ServiceType, common.GameId)
	// 创建一个通道来接收操作系统的信号
	signalCh := make(chan os.Signal, 1)
	// 监听 Interrupt 信号（Ctrl+C）
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	// 启动一个 goroutine 来等待信号
	go func() {
		// 等待信号
		sig := <-signalCh
		zlog.Logger.Info().Msgf("接收到信号 %v", sig)
		// 在关闭前执行一些清理逻辑
		// 在关闭前执行一些清理逻辑
		model.IsStop = true
		// 等待数据存储5秒
		time.Sleep(5 * time.Second)
		cleanup()
		// 退出程序
		os.Exit(0)
	}()
	// 根据服务器类型启动服务器
	switch common.ServiceType {
	case common.LOGIC_SERVER_TYPE:
		// 逻辑服
		logic.Start()
	default:
		zlog.Logger.Panic().Msgf("不是期望类型 %v", common.ServiceType)
	}

}

func cleanup() {
	// 在这里执行关闭前的清理逻辑
	zlog.Logger.Info().Msgf("服务器关闭 %v", common.ServiceSign)
	// 例如，关闭数据库连接、释放资源等
	control.SaveOut()
	zlog.Logger.Info().Msgf("服务器关闭完成 %v", common.ServiceSign)

}
