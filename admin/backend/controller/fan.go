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

// Fan 粉丝模型
type Fan struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Platform     int               `json:"platform" bson:"platform"`         // 平台 1抖音 2快手 3TikTok
	AnchorID     string            `json:"anchor_id" bson:"anchor_id"`     // 关联主播ID
	OpenID       string            `json:"open_id" bson:"open_id"`         // 平台用户ID
	Nickname     string            `json:"nickname" bson:"nickname"`        // 昵称
	Avatar       string            `json:"avatar" bson:"avatar"`            // 头像
	Gender       int               `json:"gender" bson:"gender"`           // 性别 0未知 1男 2女
	GiftTotal    int               `json:"gift_total" bson:"gift_total"`   // 礼物总额
	GiftCount    int               `json:"gift_count" bson:"gift_count"`   // 礼物数量
	BattleCount  int               `json:"battle_count" bson:"battle_count"` // 参与对战次数
	Status       int               `json:"status" bson:"status"`            // 状态 1正常 2禁言
	Remark       string            `json:"remark" bson:"remark"`            // 备注
	CreatedAt    time.Time         `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at" bson:"updated_at"`
}

// GetFanList 获取粉丝列表
func GetFanList(c *gin.Context) {
	var req struct {
		Page     int    `form:"page"`
		PageSize int    `form:"page_size"`
		Platform int    `form:"platform"`
		AnchorID string `form:"anchor_id"`
		Keyword  string `form:"keyword"`
		Status   int    `form:"status"`
	}
	c.ShouldBindQuery(&req)

	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	collection := store.GetCollection("fans")

	filter := bson.M{}
	if req.Platform > 0 {
		filter["platform"] = req.Platform
	}
	if req.AnchorID != "" {
		filter["anchor_id"] = req.AnchorID
	}
	if req.Status > 0 {
		filter["status"] = req.Status
	}
	if req.Keyword != "" {
		filter["$or"] = bson.A{
			bson.M{"nickname": bson.M{"$regex": req.Keyword, "$options": "i"}},
			bson.M{"open_id": bson.M{"$regex": req.Keyword, "$options": "i"}},
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

	var fans []Fan
	if err := cursor.All(context.TODO(), &fans); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "数据解析失败"})
		return
	}

	total, _ := collection.CountDocuments(context.TODO(), filter)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":      fans,
			"total":     total,
			"page":      req.Page,
			"page_size": req.PageSize,
		},
	})
}

// UpdateFanStatus 更新粉丝状态
func UpdateFanStatus(c *gin.Context) {
	var req struct {
		ID     string `json:"id" binding:"required"`
		Status int    `json:"status" binding:"required,oneof=1 2"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	objID, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "无效的ID"})
		return
	}

	collection := store.GetCollection("fans")
	update := bson.M{
		"$set": bson.M{
			"status":     req.Status,
			"updated_at": time.Now(),
		},
	}

	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "更新成功",
	})
}
