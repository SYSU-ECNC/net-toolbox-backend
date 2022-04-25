package agent

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var token string

func init() {
	// 生成32位的token
	token = generateToken(32)
}

// 生成token
func generateToken(n int) string {
	rand.Seed(time.Now().UnixNano())
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func GetToken(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// func CheckToken(c *gin.Context) {
// 	request := c.Request
// 	request.ParseForm()
// 	name, agentToken := request.Form["name"], request.Form["token"]
// 	dbToken := models.GetAgentTokenFromDB(name[0])
// 	if agentToken[0] == token || agentToken[0] == dbToken {
// 		// 如果dbToken为空，说明是新注册的agent
// 		if dbToken == "" {
// 			models.AddAgent(name[0], token)
// 		}

// 	} else {
// 		// token非法,阻止访问agentApi
// 		c.Abort()
// 	}
// }

// 需要提交的body为json，分别有agentName、cmd、target三个键
// func AddTaskByUser(c *gin.Context) {
// 	body, _ := ioutil.ReadAll(c.Request.Body)
// 	jsonObj, _ := simplejson.NewJson(body)
// 	agentName, _ := jsonObj.Get("agentName").String()
// 	cmd, _ := jsonObj.Get("cmd").String()
// 	target, _ := jsonObj.Get("target").String()

// 	// return
// }
