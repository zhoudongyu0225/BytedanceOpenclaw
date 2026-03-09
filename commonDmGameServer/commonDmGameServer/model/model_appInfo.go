package model

import (
	"context"
	"dmGameServer/zlog"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
	"time"
)

type AppConfigInfo struct {
	AppId            string   `json:"appId"`
	GameId           string   `json:"gameId"`           // 游戏ID
	PlatformId       string   `json:"platformId"`       // 平台
	RegionId         string   `json:"regionId"`         // 地区
	PartitionId      string   `json:"partitionId"`      // 区服
	ModeId           string   `json:"modeId"`           // 游戏模式/官游/野游
	Secret           string   `json:"secret"`           // 密钥
	EnableSecondAuth bool     `json:"enableSecondAuth"` // 是否开启二次验证 开始就严重白名单
	TopGiftList      []string `json:"topGiftList"`      // 礼物置顶
	InteractiveJoin  string   `json:"interactiveJoin"`  // 快捷加入
	InteractiveGift  string   `json:"interactiveGift"`  // 快捷礼物
	Version          string   `json:"version"`          // 版本号
	ExtId            string   `json:"extId"`            //小程序id
	// 语音Appid
	VoiceAppId string `json:"voiceAppId"`
	// 语音密钥
	VoiceSecret string `json:"voiceSecret"`
}

// app信息
func GetAppConfigInfo(xid string) *AppConfigInfo {
	appConfigJson := &AppConfigInfo{}
	list := strings.Split(xid, ".")
	if len(list) != 3 {
		zlog.Logger.Error().Msgf("客户端文件的获取 err:%v %v", "xid格式不对", xid)
		return appConfigJson
	}
	gameId := list[0]     // 1
	platformId := list[1] //3
	modeId := list[2]     // 1
	// 1.3.6.1.0
	filter := bson.M{"gameid": gameId, "platformid": platformId, "modeid": modeId}
	collection := AnchorDb.Collection(CollectionAppConfigJson)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, filter).Decode(appConfigJson)
	if err != nil {
		// zlog.Logger.Error().Msgf("GetAppConfigInfo文件的获取 err:%v %v", err, xid)
		return appConfigJson
	}
	return appConfigJson
}
