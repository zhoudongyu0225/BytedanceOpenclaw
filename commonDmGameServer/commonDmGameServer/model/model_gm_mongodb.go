package model

import (
	"context"
	"dmGameServer/untils"
	"dmGameServer/zlog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ---------------------------流水总览----------------------

// OverviewOfTransactions 流水总览
type OverviewOfTransactions struct {
	CurrTime int64 `json:"currTime"` // 当前时间
	// 新增玩家数
	NewPlayerNum int32 `json:"newPlayerNum"` // 新增玩家数
	// 礼物值
	GiftValue float64 `json:"giftValue"` // 礼物值
	// 活跃主播数
	ActiveAnchorNum int32 `json:"activeAnchorNum"` // 活跃主播数
	// 主播新增数
	NewAnchorNum int32 `json:"newAnchorNum"` // 主播新增数
	// 活跃玩家数
	ActivePlayerNum int32 `json:"activePlayerNum"` // 活跃玩家数
	// 直播时长
	LiveTime int64 `json:"liveTime"` // 直播时长(秒)
}

// UpdateOverviewOfTransactionsAll
func UpdateOverviewOfTransactionsAll(addOverviewOfTransactions *OverviewOfTransactions, CollectionKey string) bool {
	// 统计全部
	addOverviewOfTransactions.CurrTime = 1
	filter := bson.M{"currtime": addOverviewOfTransactions.CurrTime}
	update := bson.M{
		"$inc": bson.M{"newplayernum": addOverviewOfTransactions.NewPlayerNum,
			"giftvalue":       addOverviewOfTransactions.GiftValue,
			"activeanchornum": addOverviewOfTransactions.ActiveAnchorNum,
			"newanchornum":    addOverviewOfTransactions.NewAnchorNum,
			"activeplayernum": addOverviewOfTransactions.ActivePlayerNum,
			"livetime":        addOverviewOfTransactions.LiveTime},
		"$set": bson.M{"currtime": addOverviewOfTransactions.CurrTime}, // 指定不更新的字段
	}
	collection := AnchorDb.Collection(CollectionKey)
	// 设置选项，如果文档不存在，则创建一个新的文档
	opts := options.Update().SetUpsert(true)
	// 执行更新操作
	_, err := collection.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		zlog.Logger.Error().Msgf("UpdateOverviewOfTransactions err:%v", err)
		return false
	}
	return true
}

// UpdateOverviewOfTransactions
func UpdateOverviewOfTransactions(addOverviewOfTransactions *OverviewOfTransactions, CollectionKey string, isLastDay ...interface{}) bool {
	// 获取今天凌晨的时间戳(秒)
	midnightTimestamp := untils.GetMidnightTimestamp()

	if len(isLastDay) > 0 {
		// 上一天
		midnightTimestamp = midnightTimestamp - 24*60*60
	}

	// 获取时间
	addOverviewOfTransactions.CurrTime = midnightTimestamp
	filter := bson.M{"currtime": addOverviewOfTransactions.CurrTime}
	update := bson.M{
		"$inc": bson.M{"newplayernum": addOverviewOfTransactions.NewPlayerNum,
			"giftvalue":       addOverviewOfTransactions.GiftValue,
			"activeanchornum": addOverviewOfTransactions.ActiveAnchorNum,
			"newanchornum":    addOverviewOfTransactions.NewAnchorNum,
			"activeplayernum": addOverviewOfTransactions.ActivePlayerNum,
			"livetime":        addOverviewOfTransactions.LiveTime},
		"$set": bson.M{"currtime": addOverviewOfTransactions.CurrTime}, // 指定不更新的字段
	}
	collection := AnchorDb.Collection(CollectionKey)
	// 设置选项，如果文档不存在，则创建一个新的文档
	opts := options.Update().SetUpsert(true)
	// 执行更新操作
	_, err := collection.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		zlog.Logger.Error().Msgf("UpdateOverviewOfTransactions err:%v", err)
		return false
	}

	// 统计总
	UpdateOverviewOfTransactionsAll(addOverviewOfTransactions, CollectionKey)

	zlog.Logger.Debug().Msgf("增加流水总览当天的字段 success:%v", addOverviewOfTransactions)
	return true
}

// GetOverviewOfTransactionsByDateAll
func GetOverviewOfTransactionsByDateAll(CollectionKey string) (list *OverviewOfTransactions) {
	// 获取集合
	collection := AnchorDb.Collection(CollectionKey)
	// 查询 currtime = 1
	filter := bson.M{"currtime": 1}
	// 执行查询
	err := collection.FindOne(context.Background(), filter).Decode(&list)
	if err != nil {
		zlog.Logger.Error().Msgf("GetOverviewOfTransactionsByDate err:%v", err)
		return nil
	}
	return list
}

// GetOverviewOfTransactionsByDate 获取流水总览
func GetOverviewOfTransactionsByDate(startIndex, endIndex int32, CollectionKey string) (list []*OverviewOfTransactions, err error, count int64) {
	zlog.Logger.Debug().Msgf("获取流水总览 [%v] %v %v", startIndex, endIndex, CollectionKey)
	// 获取集合
	collection := AnchorDb.Collection(CollectionKey)

	// 执行查询总数量
	count, err1 := collection.CountDocuments(context.Background(), bson.D{})
	if err1 != nil {
		zlog.Logger.Error().Msgf("GetOverviewOfTransactionsByDate err:%v", err1)
	}

	// 定义排序的字段和顺序，这里按照timestamp字段降序排序
	sortOptions := options.Find().SetSort(bson.D{{Key: "currtime", Value: -1}})
	limit := endIndex - startIndex + 1
	filter := bson.D{}
	// 执行查询
	cursor, err := collection.Find(context.Background(), filter, sortOptions, options.Find().SetSkip(int64(startIndex)).SetLimit(int64(limit)))
	if err != nil {
		zlog.Logger.Error().Msgf("GetOverviewOfTransactionsByDate err:%v", err)
		return nil, err, count
	}
	// 延迟关闭游标
	defer cursor.Close(context.Background())
	// 遍历结果
	for cursor.Next(context.Background()) {
		result := &OverviewOfTransactions{}
		err = cursor.Decode(result)
		if err != nil {
			zlog.Logger.Error().Msgf("Decode err:%v", err)
			continue
		}
		list = append(list, result)
	}
	// 检查游标是否出错
	if err = cursor.Err(); err != nil {
		zlog.Logger.Error().Msgf("cursor err:%v", err)
		return nil, err, count
	}
	return list, nil, count
}

// ---------------------------主播流水----------------------

// AnchorOfTransactions 主播流水
type AnchorOfTransactions struct {
	AccountId       string  `json:"accountId"`       // 唯一id
	Name            string  `json:"name"`            // 名字
	UnsealCnt       int32   `json:"unsealCnt"`       // 解封次数
	Rank3PlayerInfo string  `json:"rank3PlayerInfo"` // 前三名玩家信息
	CurrTime        int64   `json:"currTime"`        // 当前时间
	LiveTime        int64   `json:"liveTime"`        // 直播时长(秒)
	DayGiftValue    float64 `json:"dayGiftValue"`    // 当天礼物值
	AllGiftValue    float64 `json:"allGiftValue"`    // 总礼物值
	NewPlayerNum    int32   `json:"newPlayerNum"`    // 新增玩家数
	ActivePlayerNum int32   `json:"activePlayerNum"` // 活跃玩家数W
}

// GetAnchorOfTransactionsByDateAll 获取主播当日的礼物值和当日直播时长(秒)  // GetCollectionAnchorOfTransactionsKey(ws.AccountId)
func GetAnchorOfTransactionsByDateAll(AccountId string, CollectionKey string) (list *AnchorOfTransactions) {
	// 获取集合
	collection := AnchorDb.Collection(CollectionKey)
	// 查询 currtime = 1
	filter := bson.M{"currtime": untils.GetMidnightTimestamp(),
		"accountid": AccountId}
	// 执行查询
	err := collection.FindOne(context.Background(), filter).Decode(&list)
	if err != nil {
		zlog.Logger.Error().Msgf("GetAnchorOfTransactionsByDate err:%v", err)
		return &AnchorOfTransactions{}
	}
	return list
}

// UpdateAnchorOfTransactionsLiveTime 更新主播流水直播时长
func UpdateAnchorOfTransactionsLiveTime(anchorOfTransactions *AnchorOfTransactions, addLiveTime int64, CollectionKey string, isLastDay ...interface{}) bool {
	anchorOfTransactions.CurrTime = untils.GetMidnightTimestamp()
	if len(isLastDay) > 0 {
		anchorOfTransactions.CurrTime = anchorOfTransactions.CurrTime - 24*60*60
	}
	if anchorOfTransactions.AccountId == "" {
		return false
	}
	// 根据时间
	filter := bson.M{"currtime": anchorOfTransactions.CurrTime,
		"accountid": anchorOfTransactions.AccountId,
	}
	update := bson.M{
		"$set": bson.M{"currtime": anchorOfTransactions.CurrTime, // 时间
			"accountid": anchorOfTransactions.AccountId, // 唯一id(需要传)
			"name":      anchorOfTransactions.Name,      // 名字(需要传)
		},
		"$inc": bson.M{"livetime": addLiveTime},
	}
	collection := AnchorDb.Collection(CollectionKey)
	// 设置选项，如果文档不存在，则创建一个新的文档
	opts := options.Update().SetUpsert(true)
	// 执行更新操作
	_, err := collection.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		zlog.Logger.Error().Msgf("UpdateAnchorOfTransactions err:%v", err)
		return false
	}
	zlog.Logger.Debug().Msgf("更新主播流水直播时长 success:%v  添加 %v", anchorOfTransactions, addLiveTime)
	return true
}

// 更新主播流水礼物值
func UpdateAnchorOfTransactionsGiftValue(anchorOfTransactions *AnchorOfTransactions, allGiftValue, addGiftValue float64, CollectionKey string) bool {
	if anchorOfTransactions.AccountId == "" {
		return false
	}
	anchorOfTransactions.CurrTime = untils.GetMidnightTimestamp()
	// 根据时间
	filter := bson.M{"currtime": anchorOfTransactions.CurrTime,
		"accountid": anchorOfTransactions.AccountId,
	}
	update := bson.M{
		"$set": bson.M{"currtime": anchorOfTransactions.CurrTime, // 时间
			"accountid":    anchorOfTransactions.AccountId, // 唯一id(需要传)
			"name":         anchorOfTransactions.Name,      // 名字(需要传)
			"allgiftvalue": allGiftValue,                   // 总礼物值
		},
		"$inc": bson.M{"daygiftvalue": addGiftValue},
	}
	collection := AnchorDb.Collection(CollectionKey)
	// 设置选项，如果文档不存在，则创建一个新的文档
	opts := options.Update().SetUpsert(true)
	// 执行更新操作
	_, err := collection.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		zlog.Logger.Error().Msgf("UpdateAnchorOfTransactions err:%v", err)
		return false
	}

	zlog.Logger.Debug().Msgf("更新主播流水礼物值 success:%v  添加 %v", anchorOfTransactions, addGiftValue)

	return true
}

// UpdateAnchorOfTransactionsNewPlayerNum 当日新增玩家数
func UpdateAnchorOfTransactionsNewPlayerNum(anchorOfTransactions *AnchorOfTransactions, addNewPlayerNum int32, CollectionKey string) bool {
	if anchorOfTransactions.AccountId == "" {
		return false
	}
	anchorOfTransactions.CurrTime = untils.GetMidnightTimestamp()
	// 根据时间
	filter := bson.M{"currtime": anchorOfTransactions.CurrTime,
		"accountid": anchorOfTransactions.AccountId,
	}
	update := bson.M{
		"$set": bson.M{"currtime": anchorOfTransactions.CurrTime, // 时间
			"accountid": anchorOfTransactions.AccountId, // 唯一id(需要传)
			"name":      anchorOfTransactions.Name,      // 名字(需要传)
		},
		"$inc": bson.M{"newplayernum": addNewPlayerNum},
	}
	collection := AnchorDb.Collection(CollectionKey)
	// 设置选项，如果文档不存在，则创建一个新的文档
	opts := options.Update().SetUpsert(true)
	// 执行更新操作
	_, err := collection.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		zlog.Logger.Error().Msgf("UpdateAnchorOfTransactions err:%v", err)
		return false
	}
	zlog.Logger.Debug().Msgf("当日新增玩家数 success:%+v  添加 %v", anchorOfTransactions, addNewPlayerNum)

	return true
}

// UpdateAnchorOfTransactionsActivePlayerNum 当日活跃玩家数
func UpdateAnchorOfTransactionsActivePlayerNum(anchorOfTransactions *AnchorOfTransactions, addActivePlayerNum int32, CollectionKey string) bool {
	if anchorOfTransactions.AccountId == "" {
		return false
	}
	anchorOfTransactions.CurrTime = untils.GetMidnightTimestamp()
	// 根据时间
	filter := bson.M{"currtime": anchorOfTransactions.CurrTime,
		"accountid": anchorOfTransactions.AccountId,
	}
	update := bson.M{
		"$set": bson.M{"currtime": anchorOfTransactions.CurrTime, // 时间
			"accountid": anchorOfTransactions.AccountId, // 唯一id(需要传)
			"name":      anchorOfTransactions.Name,      // 名字(需要传)
		},
		"$inc": bson.M{"activeplayernum": addActivePlayerNum},
	}
	collection := AnchorDb.Collection(CollectionKey)
	// 设置选项，如果文档不存在，则创建一个新的文档
	opts := options.Update().SetUpsert(true)
	// 执行更新操作
	_, err := collection.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		zlog.Logger.Error().Msgf("UpdateAnchorOfTransactions err:%v", err)
		return false
	}
	zlog.Logger.Debug().Msgf("当日活跃玩家数 success:%v  添加 %v", anchorOfTransactions, addActivePlayerNum)
	return true
}

// GetAnchorOfTransactionsByDate 获取主播流水
func GetAnchorOfTransactionsByDate(startIndex, endIndex, sortType int32, currTime int64, CollectionKey string) (list []*AnchorOfTransactions, err error, count int64) {
	zlog.Logger.Debug().Msgf("获取主播流水 [%v] %v %v %v %v", startIndex, endIndex, sortType, currTime, CollectionKey)
	// 获取集合
	collection := AnchorDb.Collection(CollectionKey)
	// 定义排序的字段和顺序，这里按照timestamp字段降序排序
	sortOptions := options.Find().SetSort(bson.D{{Key: "allgiftvalue", Value: -1}})
	if sortType == 1 {
		// 每日
		sortOptions = options.Find().SetSort(bson.D{{Key: "daygiftvalue", Value: -1}})
	}
	limit := endIndex - startIndex + 1
	filter := bson.M{"currtime": currTime}

	// 执行查询总数量
	count, err1 := collection.CountDocuments(context.Background(), filter)
	if err1 != nil {
		zlog.Logger.Error().Msgf("GetOverviewOfTransactionsByDate err:%v", err1)
	}

	// 执行查询
	cursor, err := collection.Find(context.Background(), filter, sortOptions, options.Find().SetSkip(int64(startIndex)).SetLimit(int64(limit)))
	if err != nil {
		zlog.Logger.Error().Msgf("GetAnchorOfTransactionsByDate err:%v", err)
		return nil, err, count
	}
	// 延迟关闭游标
	defer cursor.Close(context.Background())
	// 遍历结果
	for cursor.Next(context.Background()) {
		result := &AnchorOfTransactions{}
		err = cursor.Decode(result)
		if err != nil {
			zlog.Logger.Error().Msgf("Decode err:%v", err)
			continue
		}
		list = append(list, result)
	}
	// 检查游标是否出错
	if err = cursor.Err(); err != nil {
		zlog.Logger.Error().Msgf("cursor err:%v", err)
		return nil, err, count
	}
	return list, nil, count
}

// 批量查询
func GetAnchorOfTransactionsByDateList(idArray []string, CollectionKey string) (list []*AnchorOfTransactions, err error, count int64) {
	// 获取集合
	collection := AnchorDb.Collection(CollectionKey)
	// 构建查询条件
	filter := bson.M{"accountid": bson.M{"$in": idArray}, "currtime": untils.GetMidnightTimestamp()}

	ctx := context.Background()
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		zlog.Logger.Error().Msgf("GetAnchorOfTransactionsByDate err:%v", err)
		return nil, err, count
	}

	// 遍历查询结果
	var results []*AnchorOfTransactions
	for cursor.Next(ctx) {
		result := &AnchorOfTransactions{}
		if err = cursor.Decode(&result); err != nil {
			zlog.Logger.Error().Msgf("Decode err:%v", err)
			continue
		}
		results = append(results, result)
	}
	return results, nil, count
}

// ---------------------------主播详细流水----------------------

// AnchorDetailOfTransactions 主播详细流水
type AnchorDetailOfTransactions struct {
	CurrTime        int64   `json:"currTime"`        // 当前时间
	AccountId       string  `json:"accountId"`       // 唯一id
	Name            string  `json:"name"`            // 名字
	DayGiftValue    float64 `json:"dayGiftValue"`    // 当天礼物值
	NewPlayerNum    int32   `json:"newPlayerNum"`    // 新增玩家数
	ActivePlayerNum int32   `json:"activePlayerNum"` // 活跃玩家数
}

// GetAnchorDetailOfTransactions 主播详细流水
func GetAnchorDetailOfTransactions(startIndex, endIndex int32, id, CollectionKey string) (list []*AnchorDetailOfTransactions, err error, count int64) {
	zlog.Logger.Debug().Msgf("主播详细流水 [%v] %v %v %v ", startIndex, endIndex, id, CollectionKey)
	// 获取集合
	collection := AnchorDb.Collection(CollectionKey)
	// 定义排序的字段和顺序，这里按照timestamp字段降序排序
	sortOptions := options.Find().SetSort(bson.D{{Key: "currtime", Value: -1}})
	limit := endIndex - startIndex + 1
	filter := bson.M{"accountid": id}

	// 执行查询总数量
	count, err1 := collection.CountDocuments(context.Background(), filter)
	if err1 != nil {
		zlog.Logger.Error().Msgf("GetOverviewOfTransactionsByDate err:%v", err1)
	}

	// 执行查询
	cursor, err := collection.Find(context.Background(), filter, sortOptions, options.Find().SetSkip(int64(startIndex)).SetLimit(int64(limit)))
	if err != nil {
		zlog.Logger.Error().Msgf("GetAnchorOfTransactionsByDate err:%v", err)
		return nil, err, count
	}
	// 延迟关闭游标
	defer cursor.Close(context.Background())
	// 遍历结果
	for cursor.Next(context.Background()) {
		result := &AnchorOfTransactions{}
		err = cursor.Decode(result)
		if err != nil {
			zlog.Logger.Error().Msgf("Decode err:%v", err)
			continue
		}
		ac := &AnchorDetailOfTransactions{
			CurrTime:        result.CurrTime,
			AccountId:       result.AccountId,
			Name:            result.Name,
			DayGiftValue:    result.DayGiftValue,
			NewPlayerNum:    result.NewPlayerNum,
			ActivePlayerNum: result.ActivePlayerNum,
		}
		list = append(list, ac)
	}
	// 检查游标是否出错
	if err = cursor.Err(); err != nil {
		zlog.Logger.Error().Msgf("cursor err:%v", err)
		return nil, err, count
	}
	return list, nil, count
}

// ---------------------------玩家流水----------------------

// PlayerDetailOfTransactions 玩家流水
type PlayerDetailOfTransactions struct {
	CurrTime     int64   `json:"currTime"`     // 当前时间
	Id           string  `json:"id"`           // 唯一玩家id
	Name         string  `json:"name"`         // 名字
	DayGiftValue float64 `json:"dayGiftValue"` // 当天礼物值
	AllGiftValue float64 `json:"allGiftValue"` // 总礼物值
}

// UpdatePlayerDetailOfTransactions 更新玩玩家礼物值
func UpdatePlayerDetailOfTransactions(playerDetailOfTransactions *PlayerDetailOfTransactions, allGiftValue, addGiftValue int64, CollectionKey string) bool {

	playerDetailOfTransactions.CurrTime = untils.GetMidnightTimestamp()
	filter := bson.M{"currtime": playerDetailOfTransactions.CurrTime,
		"id": playerDetailOfTransactions.Id,
	}
	update := bson.M{
		"$set": bson.M{"currtime": playerDetailOfTransactions.CurrTime, // 时间
			"id":           playerDetailOfTransactions.Id,   // 唯一id(需要传)
			"name":         playerDetailOfTransactions.Name, // 名字(需要传)
			"allgiftvalue": allGiftValue,                    // 总礼物值
		},
		"$inc": bson.M{"daygiftvalue": addGiftValue},
	}
	// 获取集合
	collection := AnchorDb.Collection(CollectionKey)
	// 设置选项，如果文档不存在，则创建一个新的文档
	opts := options.Update().SetUpsert(true)
	// 执行更新操作
	_, err := collection.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		zlog.Logger.Error().Msgf("UpdateAnchorOfTransactions err:%v", err)
		return false
	}
	zlog.Logger.Debug().Msgf("UpdateAnchorOfTransactions success:%v", playerDetailOfTransactions)
	return true
}

// GetPlayerDetailOfTransactions 获取  //  sortType=1每日  2总
func GetPlayerDetailOfTransactions(startIndex, endIndex, sortType int32, ViewTime int64, CollectionKey string) (list []*PlayerDetailOfTransactions, err error, count int64) {
	zlog.Logger.Debug().Msgf("获取主播流水 [%v] %v %v %v ", startIndex, endIndex, sortType, CollectionKey)
	// 获取集合
	collection := AnchorDb.Collection(CollectionKey)
	// 定义排序的字段和顺序，这里按照timestamp字段降序排序
	sortOptions := options.Find().SetSort(bson.D{{Key: "allgiftvalue", Value: -1}})
	if sortType == 1 {
		// 每日
		sortOptions = options.Find().SetSort(bson.D{{Key: "daygiftvalue", Value: -1}})
	}
	limit := endIndex - startIndex + 1
	filter := bson.M{"currtime": ViewTime}

	// 执行查询总数量
	count, err1 := collection.CountDocuments(context.Background(), filter)
	if err1 != nil {
		zlog.Logger.Error().Msgf("GetOverviewOfTransactionsByDate err:%v", err1)
	}

	// 执行查询
	cursor, err := collection.Find(context.Background(), filter, sortOptions, options.Find().SetSkip(int64(startIndex)).SetLimit(int64(limit)))
	if err != nil {
		zlog.Logger.Error().Msgf("GetAnchorOfTransactionsByDate err:%v", err)
		return nil, err, count
	}
	// 延迟关闭游标
	defer cursor.Close(context.Background())
	// 遍历结果
	for cursor.Next(context.Background()) {
		result := &PlayerDetailOfTransactions{}
		err = cursor.Decode(result)
		if err != nil {
			zlog.Logger.Error().Msgf("Decode err:%v", err)
			continue
		}
		list = append(list, result)
	}
	// 检查游标是否出错
	if err = cursor.Err(); err != nil {
		zlog.Logger.Error().Msgf("cursor err:%v", err)
		return nil, err, count
	}
	return list, nil, count
}
