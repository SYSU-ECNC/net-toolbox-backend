// 数据库命名都是下划线风格，因为postgreSQL对字段名大小写不敏感
package models

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"toolBox/config"

	_ "github.com/lib/pq"
)

//models我认为是服务器内部的操作，所以一旦发生错误，就要fatal，
//而在manage中，相对来说是对于客户端的响应操作，可以是只输出错误日志而不fatal

// 命名为conf是为了避免和*gin.Context的变量重名
var conf config.Config = config.GetConfig()

func connectDB() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		conf.Database.Host, conf.Database.Port, conf.Database.User, conf.Database.Password, conf.Database.Dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalln(err)
	}

	return db
}

var DB = connectDB()

//将第一次登录的用户的信息写入数据库
func AddUser(name, UnionID string) {
	_, err := DB.Exec("insert into users(user_name,union_id) values($1,$2)", name, UnionID)
	if err != nil {
		log.Fatalln(err)
	}
}

// 判断用户是否已存在
func IsExist(UnionID string) bool {
	row, err := DB.Query(`SELECT count(*) from users where union_id = $1`, UnionID)
	if err == nil {
		var count int
		for row.Next() {
			err1 := row.Scan(&count)
			if err1 == nil && count == 1 {
				return true
			}
		}
	}
	return false
}

func GetUserNameFromDB(UnionID string) string {
	row, err := DB.Query(`SELECT user_name from users where union_id = $1`, UnionID)
	if err == nil {
		var name string
		for row.Next() {
			err1 := row.Scan(&name)
			if err1 == nil {
				return name
			}
		}
	}
	return ""
}

func AddTaskToDB(userName, command string) int {
	_, err := DB.Exec(`insert into tasks(submit_at, user_name, command) values(now() at time zone 'CCT', $1, $2)`, userName, command)
	if err != nil {
		log.Fatalln(err)
	}
	row, err := DB.Query(`SELECT last_value from tasks_id_seq;`)
	if err != nil {
		log.Fatalln(err)
	}
	var taskID int
	for row.Next() {
		if err := row.Scan(&taskID); err != nil {
			log.Fatalln(err)
		}
	}
	return taskID
}

func AddExecutionToDB(taskID int, agentName string) {
	_, err := DB.Exec(`insert into execution(task_id, agent_name) values($1, $2)`, taskID, agentName)
	if err != nil {
		log.Fatalln(err)
	}
}

type Agent struct {
	AgentName      string    `db:"agent_name"`
	Token          string    `db:"token"`
	LastTimeActive time.Time `db:"last_time_active"`
}

func GetAgentListFromDB() []Agent {

	row, err := DB.Query(`SELECT * from agents`)
	if err != nil {
		log.Fatalln(err)
	}

	var agentList []Agent
	for row.Next() {
		var agent Agent
		err = row.Scan(&agent.AgentName, &agent.Token, &agent.LastTimeActive)
		if err != nil {
			log.Fatalln(err)
		}
		agentList = append(agentList, agent)
	}
	return agentList
}

func DeleteAgentFromDB(agentName string) {
	_, err := DB.Exec(`DELETE FROM agents WHERE agent_name = $1`, agentName)
	if err != nil {
		log.Fatalln(err)
	}
}

type Execution struct {
	At        time.Time `db:"submit_at" json:"at"`
	AgentName string    `db:"agent_name" json:"agent_name"`
	Result    string    `db:"result" json:"result"`
}

func GetTaskByIDFromDB(ID string) ([]Execution, error) {
	row, err := DB.Query(`
	SELECT submit_at, agent_name, result
	FROM tasks, execution
	where tasks.id = execution.task_id and tasks.id = $1`, ID)

	if err != nil {
		return nil, err
	}

	var executionList []Execution
	for row.Next() {
		var execution Execution
		err = row.Scan(&execution.At, &execution.AgentName, &execution.Result)
		if err != nil {
			return nil, err
		}
		executionList = append(executionList, execution)
	}
	return executionList, err
}
