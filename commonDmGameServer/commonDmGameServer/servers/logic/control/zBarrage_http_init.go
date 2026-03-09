package control

import (
	"dmGameServer/common"
	"dmGameServer/zlog"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var jwtKey = []byte("!@#qwe123")

type Claims struct {
	Id string // 用户的号
	// RoomId string
	jwt.StandardClaims
}

// 生成token id为用户的uid  int64结束时间戳
func GetToken(id string) (string, error, int64) {
	MinExpire := time.Hour * time.Duration(common.TokenOut)
	expirationTime := time.Now().Add(MinExpire)
	claims := &Claims{
		Id: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	// 返回未加密signature
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 利用secret签名对token加密
	if tokenString, err := token.SignedString(jwtKey); err != nil {
		zlog.Logger.Error().Msgf("生成token err:%v", err)
		return "", err, 0
	} else {
		return tokenString, nil, expirationTime.Unix()
	}
}

// ParseToken 从tokenString中解析出相关信息
func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(token *jwt.Token) (i interface{}, err error) {
			return jwtKey, nil
		})
	return token, claims, err
}

// 野游验证解析token
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		MinExpire := time.Hour * time.Duration(common.TokenOut)
		// 获取authorization header
		tokenString := ctx.GetHeader("Authorization")
		// 从tokenString中解析信息
		token, claims, err := ParseToken(tokenString)
		if err != nil {
			zlog.Logger.Error().Msgf("验证解析token err:%v tokenString:%v", err, tokenString)
			Fail(ctx, "", "权限不足")
			// 放弃后面的接口
			ctx.Abort()
			return
		}
		// 令牌有效刷新时间
		if token.Valid {
			// 在这里检查 JWT 是否需要刷新
			if time.Until(time.Unix(claims.ExpiresAt, 0)) <= 5*time.Minute {
				// 20-30分钟刷新一次
				// 刷新 JWT 过期时间
				claims.ExpiresAt = time.Now().Add(MinExpire).Unix()
				token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ = token.SignedString(jwtKey)
				// 更新到请求头
				ctx.Header("Authorization", tokenString)
				// 设置
				GetMgr().SetYeyouToken(claims.Id, tokenString)
			}
			// 若存在该用户则将用户信息写入上下文
			ctx.Set("AccountId", claims.Id)
			ctx.Next()
		} else {
			zlog.Logger.Error().Msgf("验证解析token err:%v", err)
			// ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			Fail(ctx, "", "权限不足")
			// 放弃后面的接口
			ctx.Abort()
			return
		}
	}
}

// Response 统一返回格式
func Response(ctx *gin.Context, httpStatus int, code int, data gin.H, msg string) {
	if data == nil {
		data = gin.H{
			"msg": msg,
		}
	}
	ctx.JSON(httpStatus, gin.H{
		"code": code,
		"data": data,
		"msg":  msg,
	})
}

// Success 统一成功返回
func SuccessLogin(ctx *gin.Context, code, roomId string, msg, Token string, ExpiredTime int64) {
	data := gin.H{
		"data": gin.H{
			"result":      1,
			"roomId":      roomId,
			"message":     msg,
			"token":       Token,
			"expiredTime": ExpiredTime,
		},
	}
	ctx.JSON(http.StatusOK, data)
}

// Success 统一成功返回
func Success(ctx *gin.Context, code, roomId string, msg string) {
	data := gin.H{
		"data": gin.H{
			"result":  1,
			"roomId":  roomId,
			"message": msg,
		},
	}

	ctx.JSON(http.StatusOK, data)
}

// Fail 统一失败
func Fail(ctx *gin.Context, roomId, msg string) {
	data := gin.H{
		"data": gin.H{
			"result":  0,
			"roomId":  roomId,
			"message": msg,
		},
	}
	ctx.JSON(http.StatusBadRequest, data)
}
