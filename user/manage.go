package user

import (
	"fmt"
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

type Tasks struct {
	AgentName []string `json:"name"`
	Command   string   `json:"command"`
}

func AddTasks(c *gin.Context) {
	var tasks Tasks
	if err := c.Bind(&tasks); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
		})
		return
	}
	fmt.Println(tasks)

	task_id := models.AddTaskToDB(GetNameFromCookie(c), tasks.Command)
	for _, agentName := range tasks.AgentName {
		models.AddExecutionToDB(task_id, agentName)
	}
	c.JSON(http.StatusOK, gin.H{
		"status": true,
	})
}
