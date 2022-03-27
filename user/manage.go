package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//前端接收到的要么是200状态码，要么是401状态码
func GetUserName(c *gin.Context) {
	name := GetNameFromCookie(c)
	c.JSON(http.StatusOK, gin.H{"username": name})
}
