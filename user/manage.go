package user

import (
	"net/http"
	"toolBox/models"

	"github.com/gin-gonic/gin"
)

//前端接收到的要么是200状态码，要么是401状态码
func GetUserName(c *gin.Context) {
	UnionID := GetIDFromCookie(c)

	// UnionID为空说明cookie不合法，用户处于未登录状态，返回重定向登录url，前端实现重定向
	if UnionID == "" {
		return
	}

	name := models.GetUserNameFromDB(UnionID)
	c.JSON(http.StatusOK, gin.H{"username": name})
}

// //验证身份，无效直接返回响应重定向url，有效函数返回UnionID
// func VerifyIdentity(c *gin.Context) string {
// 	token := c.Request.Header.Get("Authorization")
// 	if !VerifyToken(token) {
// 		redirect_url := "http://127.0.0.1:8888/login"
// 		c.JSON(http.StatusFound, gin.H{"msg": "wrong token!", "redirect_url": redirect_url})
// 	}
// 	_, claims, _ := ParseToken(token)
// 	return claims.Id
// }

// func GetUserName(c *gin.Context) {
// 	UnionID := VerifyIdentity(c)
// 	name := models.GetUserName(UnionID)
// 	c.JSON(http.StatusOK, gin.H{"name": name})
// }

// func Home(c *gin.Context) {
// 	Union_id := VerifyIdentity(c)
// 	c.JSON(http.StatusOK, gin.H{"open_id": Union_id})
// }
