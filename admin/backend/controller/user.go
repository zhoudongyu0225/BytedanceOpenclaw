package controller

import (
	"admin-backend/model"
	"admin-backend/store"
	"admin-backend/utils"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// 初始化默认管理员账户
func InitDefaultAdmin() {
	collection := store.GetCollection("users")
	// 检查是否已有管理员
	var existingUser model.User
	err := collection.FindOne(context.TODO(), bson.M{"username": "admin"}).Decode(&existingUser)
	if err == nil {
		return // 已存在，不需要创建
	}

	// 加密密码
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123456"), bcrypt.DefaultCost)
	adminUser := model.User{
		ID:        primitive.NewObjectID(),
		Username:  "admin",
		Password:  string(hashedPassword),
		Nickname:  "超级管理员",
		Role:      model.RoleSuperAdmin,
		Status:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	collection.InsertOne(context.TODO(), adminUser)
}

// Login 登录接口
func Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "参数错误",
		})
		return
	}

	// 查询用户
	collection := store.GetCollection("users")
	var user model.User
	err := collection.FindOne(context.TODO(), bson.M{"username": req.Username}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "用户名或密码错误",
		})
		return
	}

	// 校验密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "用户名或密码错误",
		})
		return
	}

	// 校验状态
	if user.Status != 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "账号已被禁用",
		})
		return
	}

	// 生成token
	token, err := utils.GenerateToken(user.ID.Hex(), user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "生成token失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "登录成功",
		"data": model.LoginResponse{
			Token:    token,
			UserInfo: &user,
		},
	})
}

// GetUserInfo 获取当前用户信息
func GetUserInfo(c *gin.Context) {
	userID := c.GetString("user_id")
	objID, _ := primitive.ObjectIDFromHex(userID)

	collection := store.GetCollection("users")
	var user model.User
	err := collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "获取用户信息失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": user,
	})
}
