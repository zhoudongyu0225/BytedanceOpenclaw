package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var ServiceSign string // 服务器标签 (如 logic_101 )
var ServiceType string // 服务器类型 (如 logic )
var GameId int

const (
	AdminPassWard string = "admin12345678"
)

const (
	ADMIN_SERVER_TYPE   string = "admin"   // 管理后台服
	BARRAGE_SERVER_TYPE string = "barrage" // 接受弹幕服
	LOGIC_SERVER_TYPE   string = "logic"   // 游戏逻辑服
	QD_SERVER_TYPE      string = "qd"      // 战斗服
	MATCH_SERVER_TYPE   string = "match"   // 匹配服
)

const (
	WsRedOut   = 120 // ws读 秒
	WsWriteOut = 5   // ws写 秒
	TokenOut   = 24  // token超时时间(小时)
	// PlayerPkNum pk玩家
	PlayerPkNum = 2
)

// Cors 解决跨域
func Cors() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method
		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token, x-token")
		context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PATCH, PUT")
		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		context.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			context.AbortWithStatus(http.StatusNoContent)
		}

	}
}
