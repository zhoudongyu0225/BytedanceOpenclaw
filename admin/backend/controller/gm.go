package controller

import (
	"admin-backend/store"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GMCommand GM命令模型
type GMCommand struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Platform    int               `json:"platform" bson:"platform"`        // 平台
	AnchorID   string            `json:"anchor_id" bson:"anchor_id"`    // 主播ID
	AnchorName string            `json:"anchor_name" bson:"anchor_name"` // 主播名称
	FanID      string            `json:"fan_id" bson:"fan_id"`        // 粉丝ID
	FanName    string            `json:"fan_name" bson:"fan_name"`      // 粉丝名称
	CmdType    string            `json:"cmd_type" bson:"cmd_type"`    // 命令类型: reward/mute/kick/custom
	Command    string            `json:"command" bson:"command"`     // 命令内容
	Reward     int               `json:"reward" bson:"reward"`      // 奖励值
	Description string          `json:"description" bson:"description"` // 描述
	Result     string            `json:"result" bson:"result"`       // 执行结果
	Operator   string            `json:"operator" bson:"operator"`   // 操作人
	CreatedAt  time.Time         `json:"created_at" bson:"created_at"`
}

// ExecuteGMCommand 执行GM命令
func ExecuteGMCommand(c *gin.Context) {
	var req struct {
		Platform    int    `json:"platform" binding:"required"`
		AnchorID    string `json:"anchor_id" binding:"required"`
		FanID       string `json:"fan_id" binding:"required"`
		CmdType     string `json:"cmd_type" binding:"required"`
		Reward      int    `json:"reward"`
		Command     string `json:"command"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	// 获取用户名
	username := "admin"
	if userInfo, exists := c.Get("user_info"); exists {
		if u, ok := userInfo.(map[string]interface{}); ok {
			if name, ok := u["username"].(string); ok {
				username = name
			}
		}
	}

	// 保存命令记录
	gm := GMCommand{
		Platform:    req.Platform,
		AnchorID:   req.AnchorID,
		CmdType:    req.CmdType,
		Command:    req.Command,
		Reward:     req.Reward,
		Description: req.Description,
		FanID:      req.FanID,
		Result:     "成功",
		Operator:   username,
		CreatedAt:  time.Now(),
	}

	collection := store.GetCollection("gm_commands")
	_, err := collection.InsertOne(context.TODO(), gm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "保存记录失败"})
		return
	}

	// TODO: 这里可以添加实际的GM命令执行逻辑，比如通过WebSocket发送命令到游戏服务器

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "命令执行成功",
		"data": gm,
	})
}

// GetGMHistory 获取GM命令历史
func GetGMHistory(c *gin.Context) {
	var req struct {
		Page     int    `form:"page"`
		PageSize int    `form:"page_size"`
		Platform int    `form:"platform"`
		Keyword  string `form:"keyword"`
	}
	c.ShouldBindQuery(&req)

	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	collection := store.GetCollection("gm_commands")

	filter := bson.M{}
	if req.Platform > 0 {
		filter["platform"] = req.Platform
	}
	if req.Keyword != "" {
		filter["$or"] = bson.A{
			bson.M{"anchor_name": bson.M{"$regex": req.Keyword, "$options": "i"}},
			bson.M{"fan_name": bson.M{"$regex": req.Keyword, "$options": "i"}},
			bson.M{"command": bson.M{"$regex": req.Keyword, "$options": "i"}},
		}
	}

	skip := int64((req.Page - 1) * req.PageSize)
	opts := options.Find().SetSkip(skip).SetLimit(int64(req.PageSize)).SetSort(bson.M{"created_at": -1})

	cursor, err := collection.Find(context.TODO(), filter, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "查询失败"})
		return
	}
	defer cursor.Close(context.TODO())

	var list []GMCommand
	if err := cursor.All(context.TODO(), &list); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "数据解析失败"})
		return
	}

	total, _ := collection.CountDocuments(context.TODO(), filter)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":      list,
			"total":     total,
			"page":      req.Page,
			"page_size": req.PageSize,
		},
	})
}

// SystemConfig 系统配置模型
type SystemConfig struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Key       string           `json:"key" bson:"key"`         // 配置键
	Value    string           `json:"value" bson:"value"`    // 配置值
	Type     string           `json:"type" bson:"type"`      // 类型: string/int/json/bool
	Category string           `json:"category" bson:"category"` // 分类: game/gift/platform
	Remark   string           `json:"remark" bson:"remark"`    // 备注
	Status   int              `json:"status" bson:"status"`   // 状态 1启用 2禁用
	Creator  string           `json:"creator" bson:"creator"`  // 创建人
	CreatedAt time.Time        `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time        `json:"updated_at" bson:"updated_at"`
}

// GetSystemConfig 获取系统配置
func GetSystemConfig(c *gin.Context) {
	category := c.Query("category")

	collection := store.GetCollection("system_config")

	filter := bson.M{"status": 1}
	if category != "" {
		filter["category"] = category
	}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "查询失败"})
		return
	}
	defer cursor.Close(context.TODO())

	var list []SystemConfig
	if err := cursor.All(context.TODO(), &list); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "数据解析失败"})
		return
	}

	// 返回默认配置
	if len(list) == 0 {
		list = getDefaultConfigs()
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": list,
	})
}

// UpdateSystemConfig 更新系统配置
func UpdateSystemConfig(c *gin.Context) {
	var req struct {
		Key   string `json:"key" binding:"required"`
		Value string `json:"value" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	collection := store.GetCollection("system_config")

	update := bson.M{
		"$set": bson.M{
			"value":     req.Value,
			"updated_at": time.Now(),
		},
	}

	_, err := collection.UpdateOne(context.TODO(), bson.M{"key": req.Key}, update, options.Update().SetUpsert(true))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "更新成功",
	})
}

func getDefaultConfigs() []SystemConfig {
	return []SystemConfig{
		{Key: "game_duration", Value: "180", Type: "int", Category: "game", Remark: "游戏时长(秒)"},
		{Key: "match_wait_time", Value: "30", Type: "int", Category: "game", Remark: "匹配等待时间(秒)"},
		{Key: "max_battles", Value: "50", Type: "int", Category: "game", Remark: "最大同时对局数"},
		{Key: "danmaku_interval", Value: "500", Type: "int", Category: "game", Remark: "弹幕间隔(ms)"},
	}
}
