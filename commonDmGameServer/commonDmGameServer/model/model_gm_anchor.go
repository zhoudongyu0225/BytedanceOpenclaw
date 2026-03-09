package model

import (
	"context"
	"dmGameServer/untils"
	"dmGameServer/zlog"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

//-------------------------------gm------------

// 通用gm
type CommonGm struct {
	Accountid string `json:"accountid"`
	GameId    int32  `json:"gameId"`
}

// 添加白名单
func AddWhite(AccountId string, GameId int32) error {
	collection := AnchorDb.Collection(CollectionWhiteList)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	commonGm := &CommonGm{
		Accountid: AccountId,
		GameId:    GameId,
	}
	_, err := collection.InsertOne(ctx, commonGm)
	if err != nil {
		zlog.Logger.Error().Msgf("AddAnchor err:%v", err)
		return err
	}
	return nil
}

// GetWhiteList 获取白名单
func GetWhiteList(startIndex, endIndex, GameId int32, accountId string) (list []*CommonGm, count int64, err error) {
	collection := AnchorDb.Collection(CollectionWhiteList)
	limit := endIndex - startIndex + 1
	filter := bson.M{"gameid": GameId, "accountid": accountId}
	if GameId <= 0 {
		filter = bson.M{"accountid": accountId}
		if accountId == "" {
			filter = bson.M{}
		}
	} else {
		if accountId == "" {
			filter = bson.M{"gameid": GameId}
		}
	}

	// 执行查询总数量
	count, err1 := collection.CountDocuments(context.Background(), filter)
	if err1 != nil {
		zlog.Logger.Error().Msgf("GetOverviewOfTransactionsByDate err:%v", err1)
	}
	// 执行查询
	cursor, err := collection.Find(context.Background(), filter, options.Find().SetSkip(int64(startIndex)).SetLimit(int64(limit)))
	if err != nil {
		untils.TapErr(fmt.Sprintf("err:%v", err))
		return nil, count, err
	}
	// 延迟关闭游标
	defer cursor.Close(context.Background())
	// 遍历结果
	for cursor.Next(context.Background()) {
		result := &CommonGm{}
		err = cursor.Decode(result)
		if err != nil {
			untils.TapErr(fmt.Sprintf("err:%v", err))

			continue
		}
		list = append(list, result)
	}
	// 检查游标是否出错
	if err = cursor.Err(); err != nil {
		untils.TapErr(fmt.Sprintf("err:%v", err))

		return nil, count, err
	}
	return list, count, err

}

// 删除白名单
func DeleteWhite(AccountId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"accountid": AccountId}
	collection := AnchorDb.Collection(CollectionWhiteList)
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		untils.TapErr(fmt.Sprintf("err:%v", err))

		return err
	}
	return nil
}

// IsWhite 是否在白名单
func IsWhite(AccountId string, GameId int32) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"accountid": AccountId, "gameid": GameId}
	collection := AnchorDb.Collection(CollectionWhiteList)
	commonGm := &CommonGm{}
	err := collection.FindOne(ctx, filter).Decode(commonGm)
	if err != nil {
		zlog.Logger.Info().Msgf("主播不在白名单 err:%v %v", err, AccountId)
		return false
	}
	zlog.Logger.Info().Msgf("主播在白名单 %v", AccountId)
	return true
}

// 是否黑名单 （封号）
func IsBlack(AccountId string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"accountid": AccountId}
	collection := AnchorDb.Collection(CollectionBlackList)
	commonGm := &CommonGm{}
	err := collection.FindOne(ctx, filter).Decode(commonGm)
	if err != nil {
		return false
	}
	zlog.Logger.Info().Msgf("主播在黑名单 %v", AccountId)
	return true
}

// AddBlack 添加黑名单
func AddBlack(AccountId string) error {
	collection := AnchorDb.Collection(CollectionBlackList)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	commonGm := &CommonGm{
		Accountid: AccountId,
	}
	_, err := collection.InsertOne(ctx, commonGm)
	if err != nil {
		zlog.Logger.Error().Msgf("AddAnchor err:%v", err)
		return err
	}
	return nil
}

// DeleteBlack 删除黑名单
func DeleteBlack(AccountId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"accountid": AccountId}
	collection := AnchorDb.Collection(CollectionBlackList)
	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		zlog.Logger.Error().Msgf("DeleteAnchor err:%v", err)
		return err
	}
	return nil
}
