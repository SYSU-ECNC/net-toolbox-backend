package agent

import (
	"log"
	"math/rand"
	"net/http"
	"time"
	"toolBox/models"

	"github.com/gin-gonic/gin"
)

var Token string

func init() {
	// 生成32位的token
	Token = generateToken(32)
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
	c.JSON(http.StatusOK, gin.H{"token": Token})
}

func ReSetToken(c *gin.Context) {
	Token = generateToken(32)
	GetToken(c)
}

var DB = models.DB

func GetTokenFromDB(agentName string) string {
	var token string
	row, err := DB.Query(`SELECT token FROM agents WHERE agent_name = $1`, agentName)
	if err != nil {
		log.Println(err)
	} else {
		for row.Next() {
			err = row.Scan(&token)
			if err != nil {
				log.Println(err)
			}
		}
	}
	return token
}

func RefreshTokenInDB(name, token string) {
	_, err := DB.Exec(`
	INSERT INTO agents(agent_name, token, last_time_active)
	VALUES ($1, $2, now() at time zone 'CCT')
	ON CONFLICT(agent_name)
	DO NOTHING;
	`, name, token)

	if err != nil {
		log.Println(err)
	}
}

func CheckAgentToken(name, token string) bool {
	if token == GetTokenFromDB(name) && token != "" {
		return true
	} else if token == Token {
		RefreshTokenInDB(name, token)
		return true
	}
	return false
}

type AgentInfo struct {
	Name  string `json:"name"`
	Token string `json:"token"`
}

func AgentRegister(c *gin.Context) {
	var agentInfo AgentInfo
	if err := c.Bind(&agentInfo); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
		})
		return
	}
	if !CheckAgentToken(agentInfo.Name, agentInfo.Token) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": false,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": true,
	})

}
