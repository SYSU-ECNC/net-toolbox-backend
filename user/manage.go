package user

import (
	"net/http"
	"toolBox/models"

	"github.com/gin-gonic/gin"
)

//前端接收到的要么是200状态码，要么是401状态码
func GetUserName(c *gin.Context) {
	name := GetNameFromCookie(c)
	c.JSON(http.StatusOK, gin.H{
		"code":      200,
		"user_name": name,
	})
}

type Task struct {
	AgentName []string `json:"name"`
	Command   string   `json:"command"`
}

func AddTask(c *gin.Context) {
	var task Task
	if err := c.Bind(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
		})
		return
	}

	taskID := models.AddTaskToDB(GetNameFromCookie(c), task.Command)
	for _, agentName := range task.AgentName {
		models.AddExecutionToDB(taskID, agentName)
	}
	c.JSON(http.StatusOK, gin.H{
		"status": true,
	})
}
