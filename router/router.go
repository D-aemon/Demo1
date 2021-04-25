package router

import (
	"Demo/controller"
	"Demo/middleware"
	"github.com/gin-gonic/gin"
)

func CollectRouter(r *gin.Engine) *gin.Engine {
	//处理跨域
	r.Use(middleware.Cors())
	r.POST("/api/auth/register", controller.Register)
	r.POST("/api/auth/login", controller.Login)
	//为获取用户信息添加jwt验证中间件
	r.POST("/api/auth/info", middleware.AuthMiddleware(), controller.Info)
	return r
}
