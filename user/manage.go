package user

import (
	"log"
	"net/http"
	"time"
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

type AgentResp struct {
	AgentName string `json:"name"`
	Status    bool   `json:"status"`
}

func GetAgentList(c *gin.Context) {
	agentList := models.GetAgentListFromDB()
	var agentRespList []AgentResp
	for _, agent := range agentList {
		var agentResp AgentResp
		agentResp.AgentName = agent.AgentName

		// China Standard Time UTC + 8:00
		duration := time.Since(agent.LastTimeActive) + 8*time.Hour

		// 最后见到 Agent 的时间超过 30秒，Master 就认为 Agent 离线
		agentResp.Status = duration < 30*time.Second
		agentRespList = append(agentRespList, agentResp)
	}
	c.JSON(http.StatusOK, gin.H{
		"agents": agentRespList,
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

func GetTaskListByAgentName(c *gin.Context) {
	agentTaskList, err := models.GetTaskListByAgentNameFromDB(c.Query("agent_name"))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"task_list": agentTaskList,
	})
}

func TasksHandler(c *gin.Context) {
	agentName := c.Query("agent_name")
	if agentName == "" {
		GetTasksList(c)
	} else {
		GetTaskListByAgentName(c)
	}
}
