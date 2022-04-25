// 数据库命名都是下划线风格，因为postgreSQL对字段名大小写不敏感
package models

import (
	"database/sql"
	"fmt"
	"log"
	"toolBox/config"

	_ "github.com/lib/pq"
)

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
	var task_id int
	for row.Next() {
		if err := row.Scan(&task_id); err != nil {
			log.Fatalln(err)
		}
	}
	return task_id
}

func AddExecutionToDB(task_id int, agent_name string) {
	_, err := DB.Exec(`insert into execution(task_id, agent_name) values($1, $2)`, task_id, agent_name)
	if err != nil {
		log.Fatalln(err)
	}
}
