package router

import (
	"net/http"
	"toolBox/user"

	"github.com/gin-gonic/gin"
)

// 模拟前端主页页面
func frontendHome(c *gin.Context) {
	c.String(http.StatusOK, "Hello, here is the home!")
}

func SetupRouters() *gin.Engine {
	r := gin.Default()

	// 登录页
	r.GET("/login", user.Login)
	// 飞书跳转后端url
	r.GET("/callback/path", user.Callback)
	// 模拟前端主页页面
	r.GET("/home", frontendHome)
	// r.GET("/home", user.Home)
	// 返回json对象
	r.GET("/getUserName", user.GetUserName)
	return r
}
