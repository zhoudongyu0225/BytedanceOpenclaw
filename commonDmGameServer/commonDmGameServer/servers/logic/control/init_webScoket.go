package control

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var WsProcessMap map[int16]func(c *WebsocketAnchorClient, weMessage *WbMessageModel)

func init() {
	WsProcessMap = make(map[int16]func(c *WebsocketAnchorClient, weMessage *WbMessageModel))
	// 登录
	WsProcessMap[COMMAND_LOGIN_C2S] = WsLogin
	// 心跳s
	WsProcessMap[COMMAND_PING_C2S] = PingPong
	//获取排行榜
	WsProcessMap[COMMAND_RANK_C2S] = GetRank
}

// InitControl 初始化网路请求
func InitControl(router *gin.Engine) {
	router.GET("/ws", HandleConnection)
	// 公告
	router.POST("/Notice", Notice)
	// 野游的登录注册
	v1 := router.Group("/auth")
	{
		v1.POST("/login", Login)
		v1.POST("/register", Register)
	}
	// 测试api（需要压缩和加密）
	v2 := router.Group("/api")
	// 需要统一验证是否登录ls
	v2.Use(AuthMiddleware())
	{
		v2.POST("/chat", Chat)
		v2.POST("/like", Like)
		v2.POST("/gift", Gift)
	}
	// 快手api
	v3 := router.Group("/kuaiShouApi")
	{
		// 快手直播弹幕
		v3.POST("/gift", KuaishouBarrageHandle)
		v3.POST("/chat", KuaishouBarrageHandle)
		v3.POST("/like", KuaishouBarrageHandle)
		v3.POST("/follow", KuaishouBarrageHandle)
	}

	// 抖音apis
	dy := router.Group("/dyApi")
	{
		// 快手直播弹幕
		dy.POST("/gift", DyBarrageHandle)
		dy.POST("/chat", DyBarrageHandle)
		dy.POST("/like", DyBarrageHandle)
		dy.POST("/fst", DyBarrageHandle)
	}

	// Get测试
	router.GET("/get", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}
