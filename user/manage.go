package user

import (
	"io/ioutil"
	"net/http"

	"github.com/bitly/go-simplejson"
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

// 需要提交的body为json，分别有agentName、cmd、target三个键
func AddTaskByUser(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)
	jsonObj, _ := simplejson.NewJson(body)
	agentName, _ := jsonObj.Get("agentName").String()
	cmd, _ := jsonObj.Get("cmd").String()
	target, _ := jsonObj.Get("target").String()
}
