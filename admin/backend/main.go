package main

import (
	"admin-backend/config"
	"admin-backend/logger"
	"admin-backend/router"
	"admin-backend/store"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 1. 加载配置
	if err := config.Init(); err != nil {
		panic("配置加载失败: " + err.Error())
	}

	// 2. 初始化日志
	if err := logger.Init(config.Conf.LogConfig); err != nil {
		panic("日志初始化失败: " + err.Error())
	}

	// 3. 初始化MongoDB连接
	if err := store.InitMongo(config.Conf.MongoConfig); err != nil {
		zap.L().Fatal("MongoDB连接失败", zap.Error(err))
	}
	defer store.CloseMongo()

	// 4. 初始化Redis连接
	if err := store.InitRedis(config.Conf.RedisConfig); err != nil {
		zap.L().Fatal("Redis连接失败", zap.Error(err))
	}
	defer store.CloseRedis()

	// 5. 设置 Gin 模式
	gin.SetMode(config.Conf.ServerConfig.Mode)

	// 6. 初始化路由
	r := router.InitRouter()

	// 7. 启动服务
	srv := &http.Server{
		Addr:    ":" + config.Conf.ServerConfig.Port,
		Handler: r,
	}

	// 优雅关机
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("服务启动失败", zap.Error(err))
		}
	}()

	zap.L().Info("服务启动成功", zap.String("port", config.Conf.ServerConfig.Port))

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Info("开始关闭服务...")

	// 5秒超时关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("服务关闭异常", zap.Error(err))
	}

	zap.L().Info("服务已退出")
}
