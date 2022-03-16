package user

import (
	"net/http"
	"toolBox/models"

	"github.com/gin-gonic/gin"
)

//验证身份，无效直接返回响应重定向url，有效函数返回UnionID
func VerifyIdentity(c *gin.Context) string {
	token := c.Request.Header.Get("Authorization")
	if !VerifyToken(token) {
		redirect_url := "http://127.0.0.1:8888/login"
		c.JSON(http.StatusFound, gin.H{"msg": "wrong token!", "redirect_url": redirect_url})
	}
	_, claims, _ := ParseToken(token)
	return claims.Id
}

func GetUserName(c *gin.Context) {
	UnionID := VerifyIdentity(c)
	name := models.GetUserName(UnionID)
	c.JSON(http.StatusOK, gin.H{"name": name})
}

func Home(c *gin.Context) {
	// token := c.Request.Header.Get("Authorization")
	// if !VerifyToken(token) {
	// 	redirect_url := "http://127.0.0.1:8888/login"
	// 	c.JSON(http.StatusFound, gin.H{"msg": "wrong token!", "redirect_url": redirect_url})
	// }
	// _, claims, _ := ParseToken(token)
	// open_id := claims.Id
	Union_id := VerifyIdentity(c)
	c.JSON(http.StatusOK, gin.H{"open_id": Union_id})
}
