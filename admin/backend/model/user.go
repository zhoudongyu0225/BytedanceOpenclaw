package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	RoleSuperAdmin = 1 // 超级管理员
	RoleOperator   = 2 // 运营
	RoleAnchor     = 3 // 主播
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username  string             `bson:"username" json:"username" binding:"required"`
	Password  string             `bson:"password" json:"-" binding:"required"`
	Nickname  string             `bson:"nickname" json:"nickname"`
	Role      int                `bson:"role" json:"role" binding:"required"` // 1:超级管理员 2:运营 3:主播
	Status    int                `bson:"status" json:"status"` // 1:正常 2:禁用
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token    string `json:"token"`
	UserInfo *User  `json:"user_info"`
}
