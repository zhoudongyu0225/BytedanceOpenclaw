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
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
	"time"
)

var playerMgr *PlayerMgr

// 玩家信息管理集合
type PlayerMgr struct {
	cache sync.Map
	// 公共玩家数据缓存
	commonCache sync.Map
}

// 直接初始化
func init() {
	playerMgr = &PlayerMgr{}
}

// 获取玩家信息管理集合
func GetPlayerMgr() *PlayerMgr {
	return playerMgr
}

// 散掉某个玩家的缓存
func (p *PlayerMgr) DelPlayerCache(openId string) {
	p.cache.Delete(openId)
}

// 保存
func (p *PlayerMgr) SavePlayerInfo(openId string, info *pb.OpenVo) {
	p.cache.Store(openId, info)
}

// 获取玩家信息   bool是否是新玩家
func (p *PlayerMgr) ModelGetOpenVo(openId string) (*pb.OpenVo, bool) {
	if openId == "" {
		zlog.Logger.Error().Msgf("GetOpenVo openId or accountId is 空 %v", openId)
		return nil, false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collectionKey := fmt.Sprintf(CollectionCommonPlayerInfo, common.GameId)
	data, ok := p.cache.Load(openId)
	var v *pb.OpenVo
	isNew := false
	if !ok {
		v = &pb.OpenVo{}
		filter := bson.M{"openid": openId}
		collection := AnchorDb.Collection(collectionKey)
		err := collection.FindOne(ctx, filter).Decode(v)
		if errors.Is(err, mongo.ErrNoDocuments) {
			// 玩家
			zlog.Logger.Info().Msgf("新玩家 玩家没有公共数据 %v %v", openId)
			v = &pb.OpenVo{
				OpenId: openId,
			}
			isNew = true
		} else if err != nil {
			untils.TapErr(fmt.Sprintf("数据库 err:%v", err))
		}
		// 记录非新玩家
		p.cache.Store(openId, v)
	} else {
		v = data.(*pb.OpenVo)
	}
	if v == nil {
		untils.TapErr(fmt.Sprintf("不应该出现 获取玩家数据 v is nil %v %v", openId))
	}
	return v, isNew
}

// 保存玩家的通用数据
func (p *PlayerMgr) updateOpenVoDB(tOpenVoList []*pb.OpenVo) bool {
	CollectionKey := fmt.Sprintf(CollectionCommonPlayerInfo, common.GameId)
	if len(tOpenVoList) <= 0 {
		return false
	}
	// 执行批量更新操作
	var updates []mongo.WriteModel
	for _, player := range tOpenVoList {
		if player == nil {
			zlog.Logger.Error().Msgf("保存玩家通用数据 player is nil")
			continue
		}
		filter := bson.M{"openid": player.OpenId}
		update := bson.M{"$set": player}
		updateModel := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(true)
		updates = append(updates, updateModel)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := AnchorDb.Collection(CollectionKey).BulkWrite(ctx, updates)
	if err != nil {
		untils.TapErr(fmt.Sprintf("err:%v", err))
		return false
	}
	return true

}
