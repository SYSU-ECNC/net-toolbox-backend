package router

import (
	"net/http"
	"toolBox/agent"
	"toolBox/user"

	"github.com/gin-gonic/gin"
)

// 模拟前端主页页面
func frontendHome(c *gin.Context) {
	c.String(http.StatusOK, "Hello, here is the home!")
}

func SetupRouters() *gin.Engine {
	r := gin.Default()

	r.GET("/login", user.Login)
	r.GET("/callback", user.Callback)

	// user api
	userRouters := r.Group("/user")
	userRouters.Use(user.CheckCookie)

	userRouters.GET("/name", user.GetUserName)
	userRouters.POST("/ping", user.AddTask)
	userRouters.GET("/agent/token", agent.GetToken)
	userRouters.GET("/agent/token/new", agent.ReSetToken)
	userRouters.GET("/agents", user.GetAgentList)
	userRouters.DELETE("/agent", user.DeleteAgent)
	userRouters.GET("/task/:id", user.GetTaskByID)
	userRouters.GET("/tasks", user.TasksHandler)

	userRouters.GET("/home", frontendHome)

	// agent api
	r.POST("/agent/new", agent.AgentRegister)
	r.GET("/agent/task", agent.GetTaskByAgentName)
	r.POST("/agent/task", agent.RecieveResultFromAgent)
	return r
}
