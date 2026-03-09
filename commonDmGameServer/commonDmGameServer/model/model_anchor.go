package model

import (
	"context"
	"dmGameServer/common"
	pb "dmGameServer/pb"
	"dmGameServer/untils"
	"dmGameServer/zlog"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

var anchorMgr *AnchorMgr

// 主播的 信息管理集合
type AnchorMgr struct {
	cache sync.Map
}

// 直接初始化
func InitCache() {
	anchorMgr = &AnchorMgr{}
}

// 主播账号信息的key
func GetCollectionAnchorAccount() string {
	return fmt.Sprintf(CollectionAnchorAccountByGameId, common.GameId)
}

// Add主播数据
func AddAnchor(anchor *pb.AnchorDBInfo) error {
	if anchor.AccountId == "" {
		zlog.Logger.Error().Msgf("AccountId 为空:%v", anchor.AccountId)
		return errors.New("AccountId 为空")
	}
	// 保存
	updateAnchorDB(anchor)
	return nil
}

// 更新主播数据
func updateAnchorDB(anchor *pb.AnchorDBInfo) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if anchor.AccountId == "" {
		zlog.Logger.Error().Msgf("AccountId 为空:%v", anchor)
		return
	}
	filter := bson.M{"accountid": anchor.AccountId}
	update := bson.M{"$set": anchor}
	// 设置选项，如果文档不存在，则创建一个新的文档
	opts := options.Update().SetUpsert(true)
	_, err := AnchorDb.Collection(GetCollectionAnchorAccount()).UpdateOne(ctx, filter, update, opts)
	if err != nil {
		untils.TapErr(fmt.Sprintf("UpdateAnchor err:%v", err))
		return
	}
	return
}

// 根据号码查询主播  ir是用于获取的时候提示错误
func ModelGetAnchorById(AccountId string, ir ...interface{}) (*pb.AnchorDBInfo, error) {
	data, ok := anchorMgr.cache.Load(AccountId)
	var v *pb.AnchorDBInfo
	if ok {
		v = data.(*pb.AnchorDBInfo)
		return v, nil
	}
	filter := bson.M{"accountid": AccountId}
	collection := AnchorDb.Collection(GetCollectionAnchorAccount())
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	anchor := &pb.AnchorDBInfo{}
	err := collection.FindOne(ctx, filter).Decode(anchor)
	if err != nil {
		// len(ir)》0b表示注册 注册的肯定找不到
		if len(ir) == 0 {
			zlog.Logger.Error().Msgf("GetAnchorById err:%v %v", err, AccountId)
		}
		return nil, err
	}
	anchorMgr.cache.Store(AccountId, anchor)
	return anchor, nil
}
