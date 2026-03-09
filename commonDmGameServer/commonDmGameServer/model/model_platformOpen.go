package model

import (
	"context"
	"dmGameServer/common"
	"dmGameServer/zlog"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

const (
	// slg跨平台开启
	PlatformOpen = "PlatformOpen"
)

var platformOpenInfo = &PlatformOpenInfo{}

// PlatformOpenInfo 平台开关
type PlatformOpenInfo struct {
	gameId       int32 `bson:"gameId"`
	PlatformOpen bool  `bson:"PlatformOpen"`
}

// UpdateSlgCrossPlatformOpen 获取跨平台的开关信息
func UpdateSlgCrossPlatformOpen() {
	zlog.Logger.Info().Msgf("获取跨平台的开关信息")
	// 获取mogodb中的PlatformOpenInfo的 gameId是5的、
	// 1、如果没有这个数据，就默认是关闭的
	collection := AnchorDb.Collection(PlatformOpen)
	filter := bson.M{"gameid": common.GameId}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := collection.FindOne(ctx, filter).Decode(platformOpenInfo)
	if err != nil {
		zlog.Logger.Error().Msgf("GetSlgCrossPlatformOpen err=%v", err)
		return
	}
	zlog.Logger.Info().Msgf("GetSlgCrossPlatformOpen PlatformOpen=%v", platformOpenInfo.PlatformOpen)
}

// GetPlatformOpenInfo 获取跨平台的开关信息
func GetPlatformOpenInfo() bool {
	// 默认打开
	if platformOpenInfo == nil {
		return true
	}
	return platformOpenInfo.PlatformOpen
}
