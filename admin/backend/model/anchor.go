package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	PlatformDouyin  = 1 // 抖音
	PlatformKuaishou = 2 // 快手
	PlatformTiktok  = 3 // TikTok
)

const (
	AnchorStatusNormal = 1 // 正常
	AnchorStatusBanned = 2 // 封禁
)

type Anchor struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name          string             `bson:"name" json:"name" binding:"required"`       // 主播姓名
	Nickname      string             `bson:"nickname" json:"nickname"`                   // 主播昵称
	Platform      int                `bson:"platform" json:"platform" binding:"required"` // 平台：1抖音 2快手 3TikTok
	RoomID        string             `bson:"room_id" json:"room_id" binding:"required"`   // 直播间ID
	PlatformUID   string             `bson:"platform_uid" json:"platform_uid"`            // 平台用户ID
	Status        int                `bson:"status" json:"status"`                        // 状态：1正常 2封禁
	Remark        string             `bson:"remark" json:"remark"`                        // 备注
	GameConfig    map[string]interface{} `bson:"game_config" json:"game_config"`          // 游戏配置参数
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
}

type AnchorListRequest struct {
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=100"`
	Keyword  string `form:"keyword"`
	Platform int    `form:"platform"`
	Status   int    `form:"status"`
}

// 设置默认值
func (req *AnchorListRequest) SetDefault() {
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}
}

type AnchorCreateRequest struct {
	Name         string                 `json:"name" binding:"required"`
	Nickname     string                 `json:"nickname"`
	Platform     int                    `json:"platform" binding:"required,oneof=1 2 3"`
	RoomID       string                 `json:"room_id" binding:"required"`
	PlatformUID  string                 `json:"platform_uid"`
	Status       int                    `json:"status" binding:"oneof=1 2"`
	Remark       string                 `json:"remark"`
	GameConfig   map[string]interface{} `json:"game_config"`
}

type AnchorUpdateRequest struct {
	ID           string                 `json:"id" binding:"required"`
	Name         string                 `json:"name" binding:"required"`
	Nickname     string                 `json:"nickname"`
	Platform     int                    `json:"platform" binding:"required,oneof=1 2 3"`
	RoomID       string                 `json:"room_id" binding:"required"`
	PlatformUID  string                 `json:"platform_uid"`
	Status       int                    `json:"status" binding:"oneof=1 2"`
	Remark       string                 `json:"remark"`
	GameConfig   map[string]interface{} `json:"game_config"`
}
