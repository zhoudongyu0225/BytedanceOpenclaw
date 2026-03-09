package model

import (
	"context"
	"dmGameServer/common"
	"dmGameServer/zlog"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

const (
	DatabaseDm = "DatabaseDm" // 弹幕游戏数据库
)

const (
	CollectionAnchorAccountByGameId  = "CollectionAnchorAccount.%v"                // 主播账号集合 Anchor.GameId
	CollectionCommonPlayerInfo       = "CollectionCommonPlayerInfo.%v"             // 玩家账号集合  Anchor.GameId
	CollectionClientGameJson         = "CollectionClientGameJson"                  // 客服端游戏参数 xid gameId.platformId.modeId
	CollectionAppConfigJson          = "CollectionAppConfigJson"                   // app配置参数
	CollectionWhiteList              = "CollectionWhiteList"                       // 白名单
	CollectionBlackList              = "CollectionBlackList"                       // 黑名单
	CollectionOverviewOfTransactions = "CollectionOverviewOfTransactions.%v.%v.%v" // 流水总览  Anchor.GameId, Anchor.PlatformId, Anchor.ModeId
	CollectionAnchorOfTransactions   = "CollectionAnchorOfTransactions.%v.%v.%v"   // 主播流水  Anchor.GameId, Anchor.PlatformId, Anchor.ModeId,
	CollectionPlayerOfTransactions   = "CollectionPlayerOfTransactions.%v.%v.%v"   // 玩家流水  Anchor.GameId, Anchor.PlatformId,  Anchor.ModeId,
	AnDb                             = "ADyAccountInfo"
	// 通用服务器缓存
	CollectionCommonCache = "CollectionCommonCache.%v" // gameid
)

var (
	// AnchorDb 主播数据库
	AnchorDb *mongo.Database
)

// InitMongoDB 初始化MongoDB连接
func InitMongoDB() {
	if common.GetConfConfig().MongoUrl == "" {
		return
	}
	// 创建 MongoDB 客户端选项
	clientOptions := options.Client().ApplyURI(common.GetConfConfig().MongoUrl)

	// 创建一个上下文对象，用于在连接和操作时设置超时
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 连接到MongoDB服务器
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		zlog.Logger.Panic().Msgf("连接到MongoDB服务器失败", err)
	}

	// 检查连接是否成功
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		zlog.Logger.Panic().Msgf("连接到MongoDB服务器失败", err)
	}
	// 选择特定的数据库并将其存储在全局变量中以供应用程序使用
	// 注册的主播
	AnchorDb = client.Database(DatabaseDm)

	//// 执行 createCollection 命令以创建空集合（创建后会忽略）
	//err = AnchorDb.RunCommand(context.Background(), bson.D{{"create", CollectionAnchorAccount}}).Err()
	//if err != nil {
	//	zlog.Logger.Panic().Msgf("MongoDBRunCommand失败", err)
	//}
	err = AnchorDb.RunCommand(context.Background(), bson.D{{"create", CollectionClientGameJson}}).Err()
	if err != nil {
		zlog.Logger.Panic().Msgf("MongoDBRunCommand失败", err)
	}
	err = AnchorDb.RunCommand(context.Background(), bson.D{{"create", CollectionClientGameJson}}).Err()
	if err != nil {
		zlog.Logger.Panic().Msgf("MongoDBRunCommand失败", err)
	}
	err = AnchorDb.RunCommand(context.Background(), bson.D{{"create", CollectionAppConfigJson}}).Err()
	if err != nil {
		zlog.Logger.Panic().Msgf("MongoDBRunCommand失败", err)
	}
	err = AnchorDb.RunCommand(context.Background(), bson.D{{"create", CollectionWhiteList}}).Err()
	if err != nil {
		zlog.Logger.Panic().Msgf("MongoDBRunCommand失败", err)
	}
	err = AnchorDb.RunCommand(context.Background(), bson.D{{"create", CollectionBlackList}}).Err()
	if err != nil {
		zlog.Logger.Panic().Msgf("MongoDBRunCommand失败", err)
	}

	zlog.Logger.Info().Msg(" 初始化MongoDB成功")
}

// 主播信息存档的
type AnInfoDBInfo struct {
	// 唯一id
	AccountId string `protobuf:"bytes,1,opt,name=accountId,proto3" json:"accountId,omitempty"`
	// 预充值额度(砖石)
	Money float64 `protobuf:"fixed64,9,opt,name=money,proto3" json:"money,omitempty"`
	// 累计礼物值（砖石）
	GiftValue float64 `protobuf:"fixed64,10,opt,name=gift,proto3" json:"gift,omitempty"`
	// 密码
	Password string `protobuf:"bytes,19,opt,name=password,proto3" json:"password,omitempty"`
	// Ip发给的ip地址 没有填就默认
	Ip []string `protobuf:"bytes,20,opt,name=ip,proto3" json:"ip,omitempty"`
	// IP Index 可以配置游戏的地址
	Index int32 `protobuf:"varint,21,opt,name=index,proto3" json:"index,omitempty"`
	// 密码
	Name string `json:"name"`
}

// 根据 野游号码查询主播
func GetYeYouInfoById(AccountId string) (*AnInfoDBInfo, error) {
	filter := bson.M{"accountid": AccountId}
	collection := AnchorDb.Collection(AnDb)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	anchor := &AnInfoDBInfo{}
	err := collection.FindOne(ctx, filter).Decode(anchor)
	if err != nil {
		return nil, err
	}
	return anchor, nil
}

// 添加服务器通用缓存
// string  id =1;// 可以是用户id 也可以是主播id
//
//	int32  typeCache =2;// 数据类型 可以自由发挥
//	string  jsonData =3;// 客户端的json数据
func AddCommonCache(key string, jsonData string) error {
	collection := AnchorDb.Collection(fmt.Sprintf(CollectionCommonCache, common.GameId))
	// 更新覆盖
	_, err := collection.UpdateOne(context.Background(), bson.D{{"key", key}}, bson.D{{"$set", bson.D{{"jsondata", jsonData}}}}, options.Update().SetUpsert(true))
	if err != nil {
		zlog.Logger.Error().Msgf("AddCommonCache失败", err)
	}
	return nil
}

// 获取服务器通用缓存
func GetCommonCache(key string) string {
	collection := AnchorDb.Collection(fmt.Sprintf(CollectionCommonCache, common.GameId))
	var result bson.M
	err := collection.FindOne(context.Background(), bson.D{{"key", key}}).Decode(&result)
	if err != nil {
		zlog.Logger.Error().Msgf("GetCommonCache失败", err)
		return ""
	}
	return result["jsondata"].(string)
}
