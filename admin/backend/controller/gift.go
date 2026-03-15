package controller

import (
	"admin-backend/store"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Gift 礼物配置
type Gift struct {
	ID          int    `json:"id" bson:"id"`
	Name        string `json:"name" bson:"name"`
	Icon        string `json:"icon" bson:"icon"`
	Price       int    `json:"price" bson:"price"`       // 价格(币)
	Soldier     int    `json:"soldier" bson:"soldier"`   // 获得的士兵数
	Platform    int    `json:"platform" bson:"platform"` // 平台 1抖音 2快手 3TikTok
	Description string `json:"description" bson:"description"`
	Status      int    `json:"status" bson:"status"` // 状态 1启用 2禁用
	CreatedAt   int64  `json:"created_at" bson:"created_at"`
	UpdatedAt   int64  `json:"updated_at" bson:"updated_at"`
}

// GetGiftList 获取礼物列表
func GetGiftList(c *gin.Context) {
	platform := c.Query("platform")
	collection := store.GetCollection("gift_config")

	filter := bson.M{}
	if platform != "" {
		filter["platform"] = platform
	}
	filter["status"] = 1

	cursor, err := collection.Find(context.TODO(), filter, options.Find().SetSort(bson.M{"price": 1}))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "查询失败"})
		return
	}
	defer cursor.Close(context.TODO())

	var gifts []Gift
	if err := cursor.All(context.TODO(), &gifts); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "数据解析失败"})
		return
	}

	// 如果没有数据，返回默认配置
	if len(gifts) == 0 {
		gifts = getDefaultGifts(1)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gifts,
	})
}

// getDefaultGifts 获取默认礼物配置
func getDefaultGifts(platform int) []Gift {
	defaultGifts := map[int][]Gift{
		1: { // 抖音
			{ID: 1, Name: "小心心", Icon: "❤️", Price: 1, Soldier: 10, Platform: 1, Description: "基础礼物", Status: 1},
			{ID: 2, Name: "玫瑰花", Icon: "🌹", Price: 10, Soldier: 100, Platform: 1, Description: "浪漫礼物", Status: 1},
			{ID: 3, Name: "大啤酒", Icon: "🍺", Price: 30, Soldier: 280, Platform: 1, Description: "畅饮礼物", Status: 1},
			{ID: 4, Name: "仙女棒", Icon: "🪄", Price: 52, Soldier: 500, Platform: 1, Description: "魔法礼物", Status: 1},
			{ID: 5, Name: "大喇叭", Icon: "📢", Price: 100, Soldier: 950, Platform: 1, Description: "公告礼物", Status: 1},
			{ID: 6, Name: "跑车", Icon: "🚗", Price: 200, Soldier: 1900, Platform: 1, Description: "豪气礼物", Status: 1},
			{ID: 7, Name: "城堡", Icon: "🏰", Price: 500, Soldier: 4800, Platform: 1, Description: "梦幻礼物", Status: 1},
			{ID: 8, Name: "火箭", Icon: "🚀", Price: 1000, Soldier: 9900, Platform: 1, Description: "超级火箭", Status: 1},
			{ID: 9, Name: "嘉年华", Icon: "🎆", Price: 3000, Soldier: 29000, Platform: 1, Description: "顶级烟花", Status: 1},
			{ID: 10, Name: "为你点亮", Icon: "💖", Price: 5200, Soldier: 52000, Platform: 1, Description: "真爱特效", Status: 1},
			{ID: 11, Name: "求佛", Icon: "🙏", Price: 10000, Soldier: 100000, Platform: 1, Description: "许愿特效", Status: 1},
		},
		2: { // 快手
			{ID: 1, Name: "小心心", Icon: "❤️", Price: 1, Soldier: 10, Platform: 2, Description: "基础礼物", Status: 1},
			{ID: 2, Name: "玫瑰", Icon: "🌹", Price: 10, Soldier: 100, Platform: 2, Description: "鲜花礼物", Status: 1},
			{ID: 3, Name: "棒棒糖", Icon: "🍭", Price: 20, Soldier: 200, Platform: 2, Description: "甜蜜礼物", Status: 1},
			{ID: 4, Name: "大啤酒", Icon: "🍺", Price: 50, Soldier: 500, Platform: 2, Description: "畅饮礼物", Status: 1},
			{ID: 5, Name: "仙女棒", Icon: "🪄", Price: 100, Soldier: 1000, Platform: 2, Description: "魔法礼物", Status: 1},
			{ID: 6, Name: "跑车", Icon: "🚗", Price: 200, Soldier: 2000, Platform: 2, Description: "豪气礼物", Status: 1},
			{ID: 7, Name: "城堡", Icon: "🏰", Price: 500, Soldier: 5000, Platform: 2, Description: "梦幻礼物", Status: 1},
			{ID: 8, Name: "火箭", Icon: "🚀", Price: 1000, Soldier: 10000, Platform: 2, Description: "超级火箭", Status: 1},
			{ID: 9, Name: "为你点亮", Icon: "💖", Price: 3000, Soldier: 30000, Platform: 2, Description: "真爱特效", Status: 1},
		},
		3: { // TikTok
			{ID: 1, Name: "Rose", Icon: "🌹", Price: 1, Soldier: 10, Platform: 3, Description: "Basic Gift", Status: 1},
			{ID: 2, Name: "Coffee", Icon: "☕", Price: 5, Soldier: 50, Platform: 3, Description: "Coffee Gift", Status: 1},
			{ID: 3, Name: "Dragon", Icon: "🐉", Price: 30, Soldier: 300, Platform: 3, Description: "Dragon Gift", Status: 1},
			{ID: 4, Name: "Lambo", Icon: "🏎️", Price: 100, Soldier: 1000, Platform: 3, Description: "Luxury Car", Status: 1},
			{ID: 5, Name: "Castle", Icon: "🏰", Price: 500, Soldier: 5000, Platform: 3, Description: "Dream Castle", Status: 1},
			{ID: 6, Name: "Crown", Icon: "👑", Price: 1000, Soldier: 10000, Platform: 3, Description: "Royal Crown", Status: 1},
		},
	}

	if gifts, ok := defaultGifts[platform]; ok {
		return gifts
	}
	return defaultGifts[1]
}

// GiftRecord 礼物发送记录
type GiftRecord struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	Platform    int       `json:"platform" bson:"platform"`
	AnchorID    string    `json:"anchor_id" bson:"anchor_id"`
	AnchorName  string    `json:"anchor_name" bson:"anchor_name"`
	GiftID      int       `json:"gift_id" bson:"gift_id"`
	GiftName    string    `json:"gift_name" bson:"gift_name"`
	GiftIcon    string    `json:"gift_icon" bson:"gift_icon"`
	Count       int       `json:"count" bson:"count"`
	Soldier     int       `json:"soldier" bson:"soldier"`
	UserName    string    `json:"user_name" bson:"user_name"`
	UserID      string    `json:"user_id" bson:"user_id"`
	TotalPrice  int       `json:"total_price" bson:"total_price"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
}

// SendGift 发送模拟礼物
func SendGift(c *gin.Context) {
	var req struct {
		Platform   int    `json:"platform" binding:"required"`
		AnchorID   string `json:"anchor_id" binding:"required"`
		GiftID     int    `json:"gift_id" binding:"required"`
		Count      int    `json:"count" binding:"required,min=1,max=100"`
		UserName   string `json:"user_name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	// 获取礼物信息
	gifts := getDefaultGifts(req.Platform)
	var giftInfo *Gift
	for _, g := range gifts {
		if g.ID == req.GiftID {
			giftInfo = &g
			break
		}
	}
	if giftInfo == nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "礼物不存在"})
		return
	}

	// 保存礼物记录
	record := GiftRecord{
		Platform:   req.Platform,
		AnchorID:   req.AnchorID,
		GiftID:     req.GiftID,
		GiftName:   giftInfo.Name,
		GiftIcon:   giftInfo.Icon,
		Count:      req.Count,
		Soldier:    giftInfo.Soldier * req.Count,
		UserName:   req.UserName,
		TotalPrice: giftInfo.Price * req.Count,
		CreatedAt:  time.Now(),
	}

	collection := store.GetCollection("gift_records")
	_, err := collection.InsertOne(context.TODO(), record)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "保存记录失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "礼物发送成功",
		"data": record,
	})
}

// GetGiftRecords 获取礼物记录
func GetGiftRecords(c *gin.Context) {
	platform := c.Query("platform")
	anchorID := c.Query("anchor_id")
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "20")

	collection := store.GetCollection("gift_records")

	filter := bson.M{}
	if platform != "" {
		filter["platform"] = platform
	}
	if anchorID != "" {
		filter["anchor_id"] = anchorID
	}

	skip := (parseInt(page) - 1) * parseInt(pageSize)
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(parseInt(pageSize))).SetSort(bson.M{"created_at": -1})

	cursor, err := collection.Find(context.TODO(), filter, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "查询失败"})
		return
	}
	defer cursor.Close(context.TODO())

	var records []GiftRecord
	if err := cursor.All(context.TODO(), &records); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "数据解析失败"})
		return
	}

	total, _ := collection.CountDocuments(context.TODO(), filter)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":  records,
			"total": total,
			"page":  parseInt(page),
			"page_size": parseInt(pageSize),
		},
	})
}

func parseInt(s string) int {
	var n int
	for _, c := range s {
		n = n*10 + int(c-'0')
	}
	return n
}
