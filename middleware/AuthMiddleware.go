package middleware

import (
	"Demo/common"
	"Demo/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

//gin的中间件就是一个函数，返回一个gin的HandlerFunc
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//获取authorization header
		tokenString := ctx.GetHeader("Authorization")

		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg": "权限不足",
			})
			ctx.Abort()
			return
		}
		//去除前缀
		tokenString = tokenString[7:]

		//通过json的方式，如有前缀自行去除
		//json := make(map[string]string)
		//ctx.BindJSON(&json)
		//tokenString := json["token"]
		//if tokenString == "" {
		//	ctx.JSON(http.StatusUnauthorized, gin.H{
		//		"code": 401,
		//		"msg": "权限不足",
		//	})
		//	ctx.Abort()
		//	return
		//}
		token, claims, err := common.ParseToken(tokenString)
		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg": "权限不足",
			})
			ctx.Abort()
			return
		}

		//验证通过后获取claims 中的userId
		userId := claims.UserId
		DB := common.GetDb()
		var user model.User
		DB.First(&user, userId)

		//用户是否存在
		if user.ID == 0 {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg": "权限不足",
			})
			ctx.Abort()
			return
		}

		//用户存在，将user信息写入上下文
		ctx.Set("user", user)

		ctx.Next()
	}
}
