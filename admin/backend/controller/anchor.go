package controller

import (
	"admin-backend/model"
	"admin-backend/store"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetAnchorList 获取主播列表
func GetAnchorList(c *gin.Context) {
	var req model.AnchorListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		req.SetDefault()
	}
	req.SetDefault()

	// 构建查询条件
	filter := bson.M{}
	if req.Keyword != "" {
		filter["$or"] = bson.A{
			bson.M{"name": bson.M{"$regex": req.Keyword, "$options": "i"}},
			bson.M{"nickname": bson.M{"$regex": req.Keyword, "$options": "i"}},
			bson.M{"room_id": bson.M{"$regex": req.Keyword, "$options": "i"}},
		}
	}
	if req.Platform > 0 {
		filter["platform"] = req.Platform
	}
	if req.Status > 0 {
		filter["status"] = req.Status
	}

	collection := store.GetCollection("anchors")
	// 分页
	skip := (req.Page - 1) * req.PageSize
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(req.PageSize)).SetSort(bson.M{"created_at": -1})

	// 查询数据
	cursor, err := collection.Find(context.TODO(), filter, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "查询失败",
		})
		return
	}
	defer cursor.Close(context.TODO())

	var anchors []model.Anchor
	if err := cursor.All(context.TODO(), &anchors); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "数据解析失败",
		})
		return
	}

	// 统计总数
	total, err := collection.CountDocuments(context.TODO(), filter)
	if err != nil {
		total = 0
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":  anchors,
			"total": total,
			"page":  req.Page,
			"page_size": req.PageSize,
		},
	})
}

// CreateAnchor 新增主播
func CreateAnchor(c *gin.Context) {
	var req model.AnchorCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	// 检查同平台同直播间ID是否已存在
	collection := store.GetCollection("anchors")
	var existing model.Anchor
	err := collection.FindOne(context.TODO(), bson.M{
		"platform": req.Platform,
		"room_id":  req.RoomID,
	}).Decode(&existing)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "该平台直播间ID已存在",
		})
		return
	}

	// 创建主播
	anchor := model.Anchor{
		ID:           primitive.NewObjectID(),
		Name:         req.Name,
		Nickname:     req.Nickname,
		Platform:     req.Platform,
		RoomID:       req.RoomID,
		PlatformUID:  req.PlatformUID,
		Status:       req.Status,
		Remark:       req.Remark,
		GameConfig:   req.GameConfig,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if anchor.Status == 0 {
		anchor.Status = model.AnchorStatusNormal
	}

	_, err = collection.InsertOne(context.TODO(), anchor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "创建失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "创建成功",
		"data": anchor,
	})
}

// UpdateAnchor 更新主播
func UpdateAnchor(c *gin.Context) {
	var req model.AnchorUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	objID, err := primitive.ObjectIDFromHex(req.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的ID",
		})
		return
	}

	// 检查同平台同直播间ID是否已存在（排除自己）
	collection := store.GetCollection("anchors")
	var existing model.Anchor
	err = collection.FindOne(context.TODO(), bson.M{
		"_id":      bson.M{"$ne": objID},
		"platform": req.Platform,
		"room_id":  req.RoomID,
	}).Decode(&existing)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "该平台直播间ID已存在",
		})
		return
	}

	// 更新
	update := bson.M{
		"$set": bson.M{
			"name":          req.Name,
			"nickname":      req.Nickname,
			"platform":      req.Platform,
			"room_id":       req.RoomID,
			"platform_uid":  req.PlatformUID,
			"status":        req.Status,
			"remark":        req.Remark,
			"game_config":   req.GameConfig,
			"updated_at":    time.Now(),
		},
	}

	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": objID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "更新失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "更新成功",
	})
}

// DeleteAnchor 删除主播
func DeleteAnchor(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "无效的ID",
		})
		return
	}

	collection := store.GetCollection("anchors")
	_, err = collection.DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "删除失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除成功",
	})
}

// BatchUpdateStatus 批量更新主播状态
func BatchUpdateAnchorStatus(c *gin.Context) {
	var req struct {
		IDs    []string `json:"ids" binding:"required"`
		Status int      `json:"status" binding:"required,oneof=1 2"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	objIDs := make([]primitive.ObjectID, 0, len(req.IDs))
	for _, id := range req.IDs {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			continue
		}
		objIDs = append(objIDs, objID)
	}

	if len(objIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "没有有效的ID",
		})
		return
	}

	collection := store.GetCollection("anchors")
	update := bson.M{
		"$set": bson.M{
			"status":     req.Status,
			"updated_at": time.Now(),
		},
	}

	_, err := collection.UpdateMany(context.TODO(), bson.M{"_id": bson.M{"$in": objIDs}}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "批量更新失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "批量更新成功",
	})
}
