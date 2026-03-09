package model

import (
	"context"
	"dmGameServer/common"
	pb "dmGameServer/pb"
	"dmGameServer/untils"
	"dmGameServer/zlog"
	"fmt"
)

const (
	GmmeRank = "GameRank.%v.%v.%v.%v" //  排行榜类型 游戏id  服务器类型 //月排名
)

//var  RankDayList

type RankVo struct {
	OpenId string
	// 值 用于外部显示的
	Value int64
}

// 更新 【月】rank value是增加的
func UpdateGameMonthRank(openId string, value int64, pt int32) {
	zlog.Logger.Info().Msgf("更新排行榜数据%v %v %v", openId, value, pt)
	key := fmt.Sprintf(GmmeRank, pb.RankType_PlayerRankMonth, common.GameId, pt, untils.GetFirstOfMonthTimestamp())
	//加入排行榜
	_, err := GetRedisClient().ZIncrBy(context.Background(), key, float64(value), openId).Result()
	if err != nil {
		zlog.Logger.Error().Msgf("UpdateTbgCache err:%v [%v]", err, key)
		return
	}
	//排序
	_, err = GetRedisClient().ZRevRangeWithScores(context.Background(), key, 0, 1000).Result()
	if err != nil {
		zlog.Logger.Error().Msgf("UpdateTbgCache err:%v [%v]", err, key)
		return
	}
}

// 获取 【月】 rank
// 获取主播通用排行榜 前100
func GetGameMonthRank(pt int32) []*RankVo {
	key := fmt.Sprintf(GmmeRank, pb.RankType_PlayerRankMonth, common.GameId, pt, untils.GetFirstOfMonthTimestamp())
	zs, err := GetRedisClient().ZRevRangeWithScores(context.Background(), key, 0, 99).Result()
	if err != nil {
		zlog.Logger.Error().Msgf("GetTbgAnchorCommRank err:%v [%v]", err, key)
		return nil
	}
	var res []*RankVo
	for _, z := range zs {
		rankVo := &RankVo{}
		rankVo.OpenId = z.Member.(string)
		// 值 用于外部显示的
		rankVo.Value = int64(z.Score)
		res = append(res, rankVo)
	}
	return res
}

//--------------------------------周------------------------------------

// 更新 【周】 --rank
func UpdateGameWeekRank(openId string, value int64, pt int32) {
	zlog.Logger.Info().Msgf("更新排行榜数据%v %v ", openId, value)
	key := fmt.Sprintf(GmmeRank, pb.RankType_PlayerRankWeek, common.GameId, pt, untils.GetWeekZeroTimestamp())
	//加入排行榜 上限1000条
	_, err := GetRedisClient().ZIncrBy(context.Background(), key, float64(value), openId).Result()
	if err != nil {
		zlog.Logger.Error().Msgf("UpdateTbgCache err:%v [%v]", err, key)
		return
	}
}

// 获取  【周】 rank
// 获取主播通用排行榜 前100
func GetGameWeekRank(pt int32) []*RankVo {
	key := fmt.Sprintf(GmmeRank, pb.RankType_PlayerRankWeek, common.GameId, pt, untils.GetWeekZeroTimestamp())
	zs, err := GetRedisClient().ZRevRangeWithScores(context.Background(), key, 0, 99).Result()
	if err != nil {
		zlog.Logger.Error().Msgf("GetTbgAnchorCommRank err:%v [%v]", err, key)
		return nil
	}
	var res []*RankVo
	for _, z := range zs {
		rankVo := &RankVo{}
		rankVo.OpenId = z.Member.(string)
		// 值 用于外部显示的
		rankVo.Value = int64(z.Score)
		res = append(res, rankVo)
	}
	return res
}

//-------------------------------------------获取----------

// 获取周 排名 和 积分
// 获取主播通用排行榜 这个玩家的数据
func GetGameWeekRankByRankVo(openId string, pt int32) (int64, int64) {
	key := fmt.Sprintf(GmmeRank, pb.RankType_PlayerRankWeek, common.GameId, pt, untils.GetWeekZeroTimestamp())
	// 获取 openId的排名
	rank, err := GetRedisClient().ZRevRank(context.Background(), key, openId).Result()
	if err != nil {
		rank = -1
	}
	// 获取 openId的积分
	score, err := GetRedisClient().ZScore(context.Background(), key, openId).Result()
	if err != nil {
		score = 0
	}
	return int64(rank + 1), int64(score)

}

// 获取周 排名 和 积分
// 获取主播通用排行榜 这个玩家的数据
func GetGameMonthRankByRankVo(openId string, pt int32) (int64, int64) {
	key := fmt.Sprintf(GmmeRank, pb.RankType_PlayerRankMonth, common.GameId, pt, untils.GetFirstOfMonthTimestamp())
	// 获取 openId的排名
	rank, err := GetRedisClient().ZRevRank(context.Background(), key, openId).Result()
	if err != nil {
		rank = -1
	}
	// 获取 openId的积分
	score, err := GetRedisClient().ZScore(context.Background(), key, openId).Result()
	if err != nil {
		score = 0
	}
	return int64(rank + 1), int64(score)

}
