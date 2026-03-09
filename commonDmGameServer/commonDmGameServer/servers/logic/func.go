package logic

import (
	"dmGameServer/common"
	"dmGameServer/generateConfig"
	"dmGameServer/model"
	"dmGameServer/servers/logic/control"
	"dmGameServer/untils"
	"dmGameServer/zlog"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	_ "net/http/pprof"
	"runtime/debug"
)

// Start 逻辑服启动
func Start() {
	// 初始化配置
	config.InitConfig()
	// 初始化逻辑管理器
	control.InitLogicMgr()
	// 初始化缓存
	model.InitCache()
	// 初始化弹幕管理
	control.InitBarrageMgr()
	r := gin.Default()
	// 解决跨域
	r.Use(common.Cors())
	// 捕获panic
	r.Use(RecoveryWithLog())
	if common.GetConfConfig().LogLv < 1 {
		// 默认日志
		r.Use(gin.Logger())
	}
	// 初始化游戏路由
	control.InitControl(r)
	// 初始化存储的结构体
	model.InitSaveMgr()
	// 性能分析
	go func() {
		zlog.Logger.Info().Msgf("%v", http.ListenAndServe("0.0.0.0:6060", nil))
	}()
	zlog.Logger.Info().Msgf("逻辑服启动,监听端口%v 服务器版本:%v", common.GetConfConfig().Port, "1.0.1")
	r.Run(fmt.Sprintf(":%v", common.GetConfConfig().Port))
}
func RecoveryWithLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				stack := debug.Stack()
				untils.PanicPoss(err, stack)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			}
		}()

		c.Next()
	}
}
