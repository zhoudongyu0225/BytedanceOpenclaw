package model

import (
	"context"
	"dmGameServer/common"
	"dmGameServer/untils"
	"dmGameServer/zlog"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"strconv"
)

const (
	// RedisKey redis key
	DayAnchorInfo      = "DayAnchorInfo.%v"         // 当天主播账号集合 v是时间戳凌晨时间戳 hash (跨天删掉)
	DayPlayerInfo      = "DayPlayerInfo.%v"         // 当天玩家账号集合 v是时间戳凌晨时间戳 hash  (跨天删掉)
	AnchorActivePlayer = "AnchorActivePlayer.%v.%v" // v1是主播id v2是时间戳凌晨时间戳 hash //当日主播下活跃的玩家
	AnchorPlayer       = "AnchorPlayer.%v"          // v1是主播id hash //主播下玩过的的玩家(会很大) // 1亿的玩家大概1.5g
	AnchorPlayerRank   = "AnchorPlayerRank.%v"      // 主播下玩家排行榜 zSet
	CurrAnchorInfo     = "CurrAnchorInfo"           // 正在直播的主播账号集合
	RankGiftValue      = "RankGiftValue.%v.%v.%v"   // 主播排行榜礼物值 平台 游戏 时间
	AnchorLiveTime     = "AnchorLiveTime.%v.%v.%v"  // 主播直播时长 hash 平台 游戏 时间
	GameShouru         = "GameShouru.%v.%v"         // 游戏收入统计 hash 平台 时间  f才是gameID
)

func convertToMoney(value float64) string {
	if value >= 10000 {
		return fmt.Sprintf("%.2f w", value/10000)
	}
	return fmt.Sprintf("%v", value)
}

func convertToTime(value int) string {
	if value >= 86400 {
		return fmt.Sprintf("%.v天", value/86400)
	}
	if value >= 3600 {
		return fmt.Sprintf("%.v小时", value/3600)
	}
	if value >= 60 {
		return fmt.Sprintf("%.v分钟", value/60)
	}
	return fmt.Sprintf("%v秒", value)
}

func GetLiveInfo(id string, platformId int32, gameId int32) *untils.LivInfo {
	livInfo := &untils.LivInfo{
		AccountId: id,
	}
	client := GetRedisClient()
	ctx := context.Background()
	// 获取排名1
	key := fmt.Sprintf(RankGiftValue, platformId, gameId, untils.GetMidnightTimestamp())
	rank, err := client.ZRevRank(ctx, key, id).Result()
	if err != nil {
		zlog.Logger.Error().Msgf("获取排名err %v", err)
		livInfo.Rank1 = "未排名"
	} else {
		livInfo.Rank1 = fmt.Sprintf("%v", rank+1)
	}

	// 获取成员的收入1
	score1, err := client.ZScore(ctx, key, id).Result()
	if err != nil {
		score1 = 0
	}
	livInfo.ShouRu1 = convertToMoney(score1)

	// 获取排名2
	key = fmt.Sprintf(RankGiftValue, platformId, gameId, untils.GetMonthZeroTimestamp_1())
	rank, err = client.ZRevRank(ctx, key, id).Result()
	if err != nil {
		zlog.Logger.Error().Msgf("获取排名err %v", err)
		livInfo.Rank2 = "未排名"
	} else {
		livInfo.Rank2 = fmt.Sprintf("%v", rank+1)
	}

	// 获取成员的收入2
	score2, err := client.ZScore(ctx, key, id).Result()
	if err != nil {
		score2 = 0
	}
	livInfo.ShouRu2 = convertToMoney(score2)

	// 获取排名3
	key = fmt.Sprintf(RankGiftValue, platformId, gameId, 0)
	rank, err = client.ZRevRank(ctx, key, id).Result()
	if err != nil {
		zlog.Logger.Error().Msgf("获取排名err %v", err)
		livInfo.Rank3 = "未排名"
	} else {
		livInfo.Rank3 = fmt.Sprintf("%v", rank+1)
	}

	// 获取成员的收入3
	score3, err := client.ZScore(ctx, key, id).Result()
	if err != nil {
		score3 = 0
	}
	livInfo.ShouRu3 = convertToMoney(score3)

	// 获取时间1
	key = fmt.Sprintf(AnchorLiveTime, platformId, gameId, untils.GetMidnightTimestamp())
	timeStr, err := client.HGet(ctx, key, id).Result()
	if err != nil || timeStr == "" {
		timeStr = "0"
	}
	time1, err := strconv.Atoi(timeStr)
	if err != nil {
		zlog.Logger.Error().Msgf("strconv.Atoi err %v", timeStr)
	}
	if time1 <= 0 {
		time1 = 1
	}
	livInfo.Time1 = convertToTime(time1)

	// 获取时间2
	key = fmt.Sprintf(AnchorLiveTime, platformId, gameId, untils.GetMonthZeroTimestamp_1())
	timeStr, err = client.HGet(ctx, key, id).Result()
	if err != nil || timeStr == "" {
		timeStr = "0"
	}
	time2, err := strconv.Atoi(timeStr)
	if err != nil {
		zlog.Logger.Error().Msgf("strconv.Atoi err %v", timeStr)
	}
	if time2 <= 0 {
		time2 = 1
	}
	livInfo.Time2 = convertToTime(time2)

	// 获取时间3
	key = fmt.Sprintf(AnchorLiveTime, platformId, gameId, 0)
	timeStr, err = client.HGet(ctx, key, id).Result()
	if err != nil || timeStr == "" {
		timeStr = "0"
	}
	time3, err := strconv.Atoi(timeStr)
	if err != nil {
		zlog.Logger.Error().Msgf("strconv.Atoi err %v", timeStr)
	}
	if time3 <= 0 {
		time3 = 1
	}
	livInfo.Time3 = convertToTime(time3)

	livInfo.Dps1 = fmt.Sprintf("%.2f/h", (score1/float64(time1))*3600)
	livInfo.Dps2 = fmt.Sprintf("%.2f/h", (score2/float64(time2))*3600)
	livInfo.Dps3 = fmt.Sprintf("%.2f/h", (score3/float64(time3))*3600)
	// 游戏的收入
	todaySrRedisKey := fmt.Sprintf(GameShouru, platformId, untils.GetMidnightTimestamp())
	weekSrRedisKey := fmt.Sprintf(GameShouru, platformId, untils.GetMonthZeroTimestamp_1())
	allSrRedisKey := fmt.Sprintf(GameShouru, platformId, 0)
	// 收入1
	srStr1, err := client.HGet(ctx, todaySrRedisKey, fmt.Sprintf("%v", gameId)).Result()
	if err != nil || srStr1 == "" {
		srStr1 = "0"
	}

	sr1, err := strconv.Atoi(srStr1)
	if err != nil {
		zlog.Logger.Error().Msgf("strconv.Atoi err %v", srStr1)
	}
	livInfo.GameSR1 = convertToMoney(float64(sr1))

	// 收入2
	srStr2, err := client.HGet(ctx, weekSrRedisKey, fmt.Sprintf("%v", gameId)).Result()
	if err != nil || srStr2 == "" {
		srStr2 = "0"
	}
	sr2, err := strconv.Atoi(srStr2)
	if err != nil {
		zlog.Logger.Error().Msgf("strconv.Atoi err %v", srStr2)
	}
	livInfo.GameSR2 = convertToMoney(float64(sr2))

	// 收入3
	srStr3, err := client.HGet(ctx, allSrRedisKey, fmt.Sprintf("%v", gameId)).Result()
	if err != nil || srStr3 == "" {
		srStr3 = "0"
	}
	sr3, err := strconv.Atoi(srStr3)
	if err != nil {
		zlog.Logger.Error().Msgf("strconv.Atoi err %v", srStr3)
	}
	livInfo.GameSR3 = convertToMoney(float64(sr3))

	return livInfo
}

// 添加主播直播时间
func UpdateAnchorLiveTime(id string, value int64, platformId int32, gameId int32) bool {
	todayRedisKey := fmt.Sprintf(AnchorLiveTime, platformId, gameId, untils.GetMidnightTimestamp())

	weekRedisKey := fmt.Sprintf(AnchorLiveTime, platformId, gameId, untils.GetMonthZeroTimestamp_1())

	allRedisKey := fmt.Sprintf(AnchorLiveTime, platformId, gameId, 0)

	client := GetRedisClient()
	ctx := context.Background()
	// 今天
	_, err := client.HIncrBy(ctx, todayRedisKey, id, value).Result()
	if err != nil {
		zlog.Logger.Error().Msgf("添加直播时间err %v", err)
		return false
	}
	// 当月
	_, err = client.HIncrBy(ctx, weekRedisKey, id, value).Result()
	if err != nil {
		zlog.Logger.Error().Msgf("添加直播时间err %v", err)
		return false
	}

	// 总
	_, err = client.HIncrBy(ctx, allRedisKey, id, value).Result()
	if err != nil {
		zlog.Logger.Error().Msgf("添加直播时间err %v", err)
		return false
	}

	return true
}

// 主播添加收入排行榜
func UpdateAnchorRankList(id string, value int64, platformId int32, gameId int32) bool {
	todayRedisKey := fmt.Sprintf(RankGiftValue, platformId, gameId, untils.GetMidnightTimestamp())

	weekRedisKey := fmt.Sprintf(RankGiftValue, platformId, gameId, untils.GetMonthZeroTimestamp_1())

	allRedisKey := fmt.Sprintf(RankGiftValue, platformId, gameId, 0)

	// --------------当天-----------------
	client := GetRedisClient()
	ctx := context.Background()
	// 先将玩家添加到有序集合
	_, err := client.ZIncrBy(ctx, todayRedisKey, float64(value), id).Result()
	if err != nil {
		zlog.Logger.Error().Msgf("添加到排行榜err %v", err)
		return false
	}
	// 1000
	maxNum := int64(1000)
	// 如果有超过10000个玩家，删除排名最后的玩家
	count := GetRedisClient().ZCard(ctx, todayRedisKey).Val()
	if count > maxNum {
		GetRedisClient().ZRemRangeByRank(ctx, todayRedisKey, 0, count-maxNum-1)
	}

	// ---------------当天----------------
	client = GetRedisClient()
	ctx = context.Background()
	// 先将玩家添加到有序集合
	_, err = client.ZIncrBy(ctx, weekRedisKey, float64(value), id).Result()
	if err != nil {
		zlog.Logger.Error().Msgf("添加到排行榜err %v", err)
		return false
	}
	// 1000
	maxNum = int64(1000)
	// 如果有超过10000个玩家，删除排名最后的玩家
	count = GetRedisClient().ZCard(ctx, weekRedisKey).Val()
	if count > maxNum {
		GetRedisClient().ZRemRangeByRank(ctx, weekRedisKey, 0, count-maxNum-1)
	}

	// --------------全部排名-----------------
	client = GetRedisClient()
	ctx = context.Background()
	// 先将玩家添加到有序集合
	_, err = client.ZIncrBy(ctx, allRedisKey, float64(value), id).Result()
	if err != nil {
		zlog.Logger.Error().Msgf("添加到排行榜err %v", err)
		return false
	}
	// 1000
	maxNum = int64(1000)
	// 如果有超过10000个玩家，删除排名最后的玩家
	count = GetRedisClient().ZCard(ctx, allRedisKey).Val()
	if count > maxNum {
		GetRedisClient().ZRemRangeByRank(ctx, allRedisKey, 0, count-maxNum-1)
	}

	// -----------------游戏收入-------------------------
	todaySrRedisKey := fmt.Sprintf(GameShouru, platformId, untils.GetMidnightTimestamp())
	weekSrRedisKey := fmt.Sprintf(GameShouru, platformId, untils.GetMonthZeroTimestamp_1())
	allSrRedisKey := fmt.Sprintf(GameShouru, platformId, 0)
	// 今天
	_, err = client.HIncrBy(ctx, todaySrRedisKey, fmt.Sprintf("%v", gameId), value).Result()
	if err != nil {
		zlog.Logger.Error().Msgf("添加游戏收入err %v", err)
		return false
	}
	// 当月
	_, err = client.HIncrBy(ctx, weekSrRedisKey, fmt.Sprintf("%v", gameId), value).Result()
	if err != nil {
		zlog.Logger.Error().Msgf("添加游戏收入err %v", err)
		return false
	}
	// 总
	_, err = client.HIncrBy(ctx, allSrRedisKey, fmt.Sprintf("%v", gameId), value).Result()
	if err != nil {
		zlog.Logger.Error().Msgf("添加游戏收入err %v", err)
		return false
	}

	return true
}

var rdb *redis.Client
var Ctx context.Context = context.Background()

// GetRedisClient 获取redis
func GetRedisClient() *redis.Client {
	if rdb == nil {
		InitRedis()
	}
	return rdb
}

// InitRedis 初始化redis
func InitRedis() {
	if common.GetConfConfig().RedisUrl == "" {
		return
	}
	rdb = redis.NewClient(&redis.Options{
		Addr:     common.GetConfConfig().RedisUrl,
		Password: common.GetConfConfig().RedisPassword,
		DB:       0,
		PoolSize: 20,
	})
	_, err := rdb.Ping(Ctx).Result()
	if err != nil {
		log.Fatalln("redis err", err)
		return
	}
	zlog.Logger.Info().Msg(" 初始化redis成功")
}
