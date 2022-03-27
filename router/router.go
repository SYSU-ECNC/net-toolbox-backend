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
	userRouters.GET("/home", frontendHome)
	userRouters.GET("/registration/token", agent.GetToken)
	userRouters.POST("/add/task", agent.AddTaskByUser)

	// protectedRouters.GET("tasks", task.GetTasks)
	// protectedRouters.GET("/agents", agent.GetAgents)

	// agent api
	agentRouters := r.Group("/agent")
	agentRouters.Use(agent.CheckToken)

	agentRouters.POST("/task", agent.GetTask)
	return r
}
