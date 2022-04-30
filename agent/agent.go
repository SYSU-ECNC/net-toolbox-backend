package agent

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
	"toolBox/models"

	"github.com/gin-gonic/gin"
)

func init() {
	// 生成32位的token
	Token = generateToken(32)
	// 设置日志文件
	logFile, err := os.OpenFile("./agent.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("open log file failed, err:", err)
		return
	}
	log.SetOutput(logFile)
}

var Token string

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

func GetTaskByAgentNameFromDB(name string) (taskID int, command string) {
	row, err := DB.Query(`
	SELECT execution.task_id, tasks.command
	FROM tasks, execution
	WHERE tasks.id = execution.task_id and execution.agent_name = $1 and is_exec = FALSE
	LIMIT 1
	`, name)
	if err != nil {
		log.Println(err)
		return 0, ""
	}

	for row.Next() {
		err = row.Scan(&taskID, &command)
		if err != nil {
			log.Println(err)
		}
	}
	return taskID, command
}

func GetTaskByAgentName(c *gin.Context) {
	token := c.Request.Header["Authorization"][0]
	name := c.Query("agent_name")
	if !CheckAgentToken(name, token) {
		c.JSON(http.StatusUnauthorized, nil)
		return
	}
	taskID, command := GetTaskByAgentNameFromDB(name)
	c.JSON(http.StatusOK, gin.H{
		"task_id": taskID,
		"command": command,
	})
}

type ExecResult struct {
	Name   string `json:"name"`
	TaskID int    `json:"task_id"`
	Result string `json:"result"`
}

func UpdateResultToExecutionTable(execResult ExecResult) {
	_, err := DB.Exec(`
	UPDATE execution 
	SET result = $3, is_exec = TRUE 
	WHERE agent_name = $1 and task_id = $2
	`, execResult.Name, execResult.TaskID, execResult.Result)
	if err != nil {
		log.Println(err)
	}
}

func RecieveResultFromAgent(c *gin.Context) {
	token := c.Request.Header["Authorization"][0]

	var execResult ExecResult
	if err := c.Bind(&execResult); err != nil {
		c.JSON(http.StatusBadRequest, nil)
		log.Println(err)
		// fmt.Println(execResult)
		return
	}
	if !CheckAgentToken(execResult.Name, token) {
		c.JSON(http.StatusUnauthorized, nil)
		return
	}
	UpdateResultToExecutionTable(execResult)
	c.JSON(http.StatusOK, nil)
}
