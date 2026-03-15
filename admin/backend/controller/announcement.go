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

// Announcement 公告模型
type Announcement struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string            `json:"title" bson:"title"`           // 标题
	Content     string            `json:"content" bson:"content"`       // 内容
	Platform    int               `json:"platform" bson:"platform"`      // 平台 0全部 1抖音 2快手 3TikTok
	Type        int               `json:"type" bson:"type"`            // 类型 1系统 2活动
	Status      int               `json:"status" bson:"status"`        // 状态 1显示 2隐藏
	StartTime   time.Time         `json:"start_time" bson:"start_time"` // 开始时间
	EndTime    time.Time         `json:"end_time" bson:"end_time"`     // 结束时间
	Sort       int               `json:"sort" bson:"sort"`            // 排序
	Creator    string            `json:"creator" bson:"creator"`      // 创建人
	CreatedAt  time.Time         `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at" bson:"updated_at"`
}

// GetAnnouncementList 获取公告列表
func GetAnnouncementList(c *gin.Context) {
	var req struct {
		Page     int `form:"page"`
		PageSize int `form:"page_size"`
		Platform int `form:"platform"`
		Status   int `form:"status"`
	}
	c.ShouldBindQuery(&req)

	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	collection := store.GetCollection("announcements")

	filter := bson.M{}
	if req.Platform > 0 {
		filter["platform"] = req.Platform
	}
	if req.Status > 0 {
		filter["status"] = req.Status
	}

	skip := int64((req.Page - 1) * req.PageSize)
	opts := options.Find().SetSkip(skip).SetLimit(int64(req.PageSize)).SetSort(bson.M{"sort": -1, "created_at": -1})

	cursor, err := collection.Find(context.TODO(), filter, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "查询失败"})
		return
	}
	defer cursor.Close(context.TODO())

	var list []Announcement
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

// CreateAnnouncement 创建公告
func CreateAnnouncement(c *gin.Context) {
	var req Announcement
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	if req.Title == "" || req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "标题和内容不能为空"})
		return
	}

	now := time.Now()
	req.ID = primitive.NewObjectID()
	req.CreatedAt = now
	req.UpdatedAt = now
	if req.Status == 0 {
		req.Status = 1
	}
	if req.Platform == 0 {
		req.Platform = 0 // 全部平台
	}
	if req.StartTime.IsZero() {
		req.StartTime = now
	}

	collection := store.GetCollection("announcements")
	_, err := collection.InsertOne(context.TODO(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "创建成功",
		"data": req,
	})
}

// UpdateAnnouncement 更新公告
func UpdateAnnouncement(c *gin.Context) {
	var req Announcement
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	if req.ID.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "ID不能为空"})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"title":      req.Title,
			"content":    req.Content,
			"platform":   req.Platform,
			"type":       req.Type,
			"status":     req.Status,
			"start_time": req.StartTime,
			"end_time":   req.EndTime,
			"sort":       req.Sort,
			"updated_at":  time.Now(),
		},
	}

	collection := store.GetCollection("announcements")
	_, err := collection.UpdateOne(context.TODO(), bson.M{"_id": req.ID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "更新成功",
	})
}

// DeleteAnnouncement 删除公告
func DeleteAnnouncement(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "无效的ID"})
		return
	}

	collection := store.GetCollection("announcements")
	_, err = collection.DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}
