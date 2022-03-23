package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//前端接收到的要么是200状态码，要么是401状态码
func GetUserName(c *gin.Context) {
	name := GetNameFromCookie(c)

	// name为空说明cookie不合法，用户处于未登录状态，返回重定向登录url(这个由Verify实现)，前端实现重定向
	if name == "" {
		return
	}
	c.JSON(http.StatusOK, gin.H{"username": name})
}
