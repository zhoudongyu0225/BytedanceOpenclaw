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
			// 后续添加用户管理、主播管理等接口
		}
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "ok",
		})
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
