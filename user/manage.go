package user

import (
	"log"
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

type FormAgentName struct {
	AgentName string `json:"name"`
}

func DeleteAgent(c *gin.Context) {
	var agentName FormAgentName
	if err := c.Bind(&agentName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": false,
		})
		log.Println(err)
		return
	}

	models.DeleteAgentFromDB(agentName.AgentName)

	c.JSON(http.StatusOK, gin.H{
		"code": true,
	})
}

func GetTaskByID(c *gin.Context) {
	ID := c.Param("id")
	executionList, err := models.GetTaskByIDFromDB(ID)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"results": executionList,
	})
}

func GetTasksList(c *gin.Context) {
	taskRespList, err := models.GetTasksListFromDB()
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"task_list": taskRespList,
	})
}
