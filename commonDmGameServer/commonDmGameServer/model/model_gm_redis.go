package model

import (
	"context"
	"dmGameServer/untils"
	"dmGameServer/zlog"
	"fmt"
)

// IsHaveDayAnchorInfo 今日主播是否存在 不存在加入
func IsHaveDayAnchorInfo(id string) bool {
	// 判断是否存在
	// 1.获取redis
	client := GetRedisClient()
	// 2.获取key
	key := fmt.Sprintf(DayAnchorInfo, untils.GetMidnightTimestamp())
	// 3.判断是否存在
	isHave, err := client.HExists(Ctx, key, id).Result()
	if err != nil {
		zlog.Logger.Error().Msgf("IsHaveDayPlayerInfo err:%v", err)
		return true
	}
	if !isHave {
		// 不存在加入
		client.HSet(Ctx, key, id, 0)
		return false
	}
	return true
}

// IsHaveDayPlayerInfo 今日玩家是否存在 不存在加入
func IsHaveDayPlayerInfo(id string) bool {
	// 判断是否存在
	// 1.获取redis
	client := GetRedisClient()
	// 2.获取key
	key := fmt.Sprintf(DayPlayerInfo, untils.GetMidnightTimestamp())
	// 3.判断是否存在
	isHave, err := client.HExists(Ctx, key, id).Result()
	if err != nil {
		zlog.Logger.Error().Msgf("IsHaveDayPlayerInfo err:%v", err)
		return true
	}
	if !isHave {
		// 不存在加入
		client.HSet(Ctx, key, id, 0)
		return false
	}

	return true
}

// IsAnchorActivePlayer v1是主播id v2是时间戳凌晨时间戳 hash //当日主播下活跃的玩家
func IsAnchorActivePlayer(id, AnchorId string) bool {
	// 判断是否存在
	// 1.获取redis
	client := GetRedisClient()
	// 2.获取key
	key := fmt.Sprintf(AnchorActivePlayer, AnchorId, untils.GetMidnightTimestamp())
	// 3.判断是否存在
	isHave, err := client.HExists(Ctx, key, id).Result()
	if err != nil {
		zlog.Logger.Error().Msgf("IsAnchorActivePlayer err:%v", err)
		return true
	}
	if !isHave {
		// 不存在加入
		client.HSet(Ctx, key, id, 0)
		return false
	}
	return true
}

// IsAnchorPlayer AnchorPlayer 玩家在这个主播玩过的游戏
func IsAnchorPlayer(id, AnchorId string) bool {
	// 判断是否存在
	// 1.获取redis
	client := GetRedisClient()
	// 2.获取key
	key := fmt.Sprintf(AnchorPlayer, AnchorId)
	// 3.判断是否存在
	isHave, err := client.HExists(Ctx, key, id).Result()
	if err != nil {
		zlog.Logger.Error().Msgf("IsAnchorActivePlayer err:%v", err)
		return true
	}
	if !isHave {
		// 不存在加入
		client.HSet(Ctx, key, id, 0)
		return false
	}
	return true
}

type AnchorPlayerRankInfo struct {
	NickName  string `json:"nickName"`
	Rank      int64  `json:"rank"`
	GiftValue int64  `json:"giftValue"`
}

// 当前时间
var curTime int64 = 0

// GetAnchorPlayerRank 前3名
func GetAnchorPlayerRank(anchorId string, GameId, PlatformId, RegionId, ModeId, PartitionId int32) string {
	//// 30秒清理一次
	//if curTime+30 < int64(time.Now().Unix()) {
	//	// 清理
	//	GetPlayerMgr().Rw.Lock()
	//	GetPlayerMgr().Players = make(map[string]*pb.OpenVo)
	//	GetPlayerMgr().Rw.Unlock()
	//	curTime = int64(time.Now().Unix())
	//}
	//
	//key := fmt.Sprintf(AnchorPlayerRank, anchorId)
	//ctx := context.Background()
	//// 获取排行榜数据
	//result, err := GetRedisClient().ZRevRangeWithScores(ctx, key, int64(0), int64(-1)).Result()
	//if err != nil {
	//	zlog.Logger.Error().Msgf("获取排行榜错误 err:%v [%v]", err, anchorId)
	//	return ""
	//}
	//vStr := ""
	//rank := int64(1) // 排名从 x+1 开始
	//// RankVo
	//for _, z := range result {
	//	pId := z.Member.(string)
	//	// 不需要冲缓存获取 // 理解的数据
	//	openVo, _ := GetPlayerMgr().GetOpenVo(fmt.Sprintf("%v.%v.%v.%v.%v.%v", pId, GameId, PlatformId, RegionId, ModeId, PartitionId),
	//		pId, fmt.Sprintf("CollectionPlayerInfo.%v.%v.%v.%v.%v", GameId, PlatformId, RegionId, ModeId, PartitionId))
	//
	//	a := &AnchorPlayerRankInfo{
	//		NickName:  openVo.NickName,
	//		Rank:      rank,
	//		GiftValue: openVo.GiftValue,
	//	}
	//	rank++
	//	data, _ := json.Marshal(a)
	//	vStr += string(data)
	//}
	//return vStr
	return ""
}

// CurAnchorOfTransactions 当前直播的情况
type CurAnchorOfTransactions struct {
	AccountId string `json:"accountId"` // 唯一id
	LiveUrl   string `json:"liveUrl"`   // 直播地址
}

// AddCurrCurrAnchor 正在直播的主播账号集合
func UpdateCurrCurrAnchor(curAnchorOfTransactions *CurAnchorOfTransactions) bool {
	key := CurrAnchorInfo
	ctx := context.Background()
	_, err := GetRedisClient().HSet(ctx, key, curAnchorOfTransactions.AccountId, curAnchorOfTransactions.LiveUrl).Result()
	if err != nil {
		zlog.Logger.Error().Msgf("AddCurrCurrAnchor err:%v", err)
		return false
	}
	return true
}

// DelCurrCurrAnchor 删掉
func DelCurrCurrAnchor(accountId string) bool {
	key := CurrAnchorInfo
	ctx := context.Background()
	// hash添加
	_, err := GetRedisClient().HDel(ctx, key, accountId).Result()
	if err != nil {
		zlog.Logger.Error().Msgf("DelCurrCurrAnchor err:%v", err)
		return false
	}
	return true
}

// GetCurrCurrAnchor key是主播id value是直播地址
func GetCurrCurrAnchor() map[string]string {
	key := CurrAnchorInfo
	ctx := context.Background()
	// hash获取
	result, err := GetRedisClient().HGetAll(ctx, key).Result()
	if err != nil {
		zlog.Logger.Error().Msgf("GetCurrCurrAnchor err:%v", err)
		return nil
	}
	m := make(map[string]string)
	for k, v := range result {
		m[k] = v
	}
	return m
}
