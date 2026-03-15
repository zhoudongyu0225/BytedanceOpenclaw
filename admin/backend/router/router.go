package router

import (
	"admin-backend/controller"
	"admin-backend/middleware"
	"admin-backend/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	// 初始化默认管理员
	controller.InitDefaultAdmin()

	// 跨域中间件
	r.Use(CORS())

	// 公开接口
	v1 := r.Group("/api/v1")
	{
		v1.POST("/login", controller.Login)
	}

	// 需要认证的接口
	authV1 := v1.Group("")
	authV1.Use(middleware.JWTAuth())
	{
		authV1.GET("/user/info", controller.GetUserInfo)

		// 管理员权限接口
		adminGroup := authV1.Group("")
		adminGroup.Use(middleware.RoleAuth(model.RoleSuperAdmin))
		{
			// 主播管理接口
			anchorGroup := adminGroup.Group("/anchor")
			{
				anchorGroup.GET("/list", controller.GetAnchorList)
				anchorGroup.POST("/create", controller.CreateAnchor)
				anchorGroup.POST("/update", controller.UpdateAnchor)
				anchorGroup.DELETE("/delete/:id", controller.DeleteAnchor)
				anchorGroup.POST("/batch-update-status", controller.BatchUpdateAnchorStatus)
			}

			// 礼物管理接口
			giftGroup := adminGroup.Group("/gift")
			{
				giftGroup.GET("/list", controller.GetGiftList)
				giftGroup.POST("/send", controller.SendGift)
				giftGroup.GET("/records", controller.GetGiftRecords)
			}

			// 粉丝管理接口
			fanGroup := adminGroup.Group("/fan")
			{
				fanGroup.GET("/list", controller.GetFanList)
				fanGroup.POST("/update-status", controller.UpdateFanStatus)
			}

			// 公告管理接口
			announcementGroup := adminGroup.Group("/announcement")
			{
				announcementGroup.GET("/list", controller.GetAnnouncementList)
				announcementGroup.POST("/create", controller.CreateAnnouncement)
				announcementGroup.POST("/update", controller.UpdateAnnouncement)
				announcementGroup.DELETE("/delete/:id", controller.DeleteAnnouncement)
			}

			// GM命令接口
			gmGroup := adminGroup.Group("/gm")
			{
				gmGroup.POST("/execute", controller.ExecuteGMCommand)
				gmGroup.GET("/history", controller.GetGMHistory)
			}

			// 系统配置接口
			configGroup := adminGroup.Group("/config")
			{
				configGroup.GET("/list", controller.GetSystemConfig)
				configGroup.POST("/update", controller.UpdateSystemConfig)
			}
		}
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "ok",
		})
	})

	// 文档静态服务
	r.StaticFS("/docs", http.Dir("./static/docs"))
	// 静态资源服务
	r.Static("/assets", "./static/assets")
	r.StaticFile("/vite.svg", "./static/vite.svg")
	// 前端页面路由
	r.NoRoute(func(c *gin.Context) {
		c.File("./static/index.html")
	})

	return r
}

// CORS 跨域中间件
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
