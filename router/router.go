package router

import (
	"toolBox/user"

	"github.com/gin-gonic/gin"
)

func SetupRouters() *gin.Engine {
	r := gin.Default()

	r.GET("/login", user.Login)
	r.GET("/callback/path", user.Callback)
	r.GET("/home", user.Home)
	r.GET("/getUserName", user.GetUserName)
	return r
}
