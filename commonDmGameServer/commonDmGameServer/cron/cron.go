package cron

import (
	"context"
	"dmGameServer/model"
	"dmGameServer/zlog"
	"github.com/robfig/cron"
)

// InitCron 初始化定时器
func InitCron() {
	// 创建一个新的Cron调度器
	c := cron.New()
	// 添加每周五下午17点执行清理任务的定时器
	c.AddFunc("0 17 * * 5", func() {
		// 清理排行榜
		// ClearRank()
	})
	// 启动定时任务
	c.Start()
	zlog.Logger.Info().Msgf("初始化定时器")
}

func ClearRank() {
	zlog.Logger.Info().Msgf("清理排行榜")
	var cursor uint64
	var keys []string
	pattern := "*RankPlayerInfo*"
	ctx := context.Background()
	for {
		var err error
		keys, cursor, err = model.GetRedisClient().Scan(ctx, cursor, pattern, 10).Result()
		if err != nil {
			zlog.Logger.Error().Msgf("Error scanning keys:%v", err)
			return
		}
		// 删除匹配的键
		if len(keys) > 0 {
			err := model.GetRedisClient().Del(ctx, keys...).Err()
			if err != nil {
				zlog.Logger.Error().Msgf("Error deleting keys:%v", err)
				return
			}
		}
		// 如果cursor为0，表示遍历完成
		if cursor == 0 {
			break
		}
	}
}
