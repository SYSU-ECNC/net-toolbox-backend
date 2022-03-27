// 数据库命名都是下划线风格，因为postgreSQL对字段名大小写不敏感
package models

import (
	"database/sql"
	"fmt"
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
		panic(err)
	}

	return db
}

var DB = connectDB()

//将第一次登录的用户的信息写入数据库
func AddUser(name, UnionID string) {
	_, err := DB.Exec("insert into ecncer(name,union_id) values($1,$2)", name, UnionID)
	if err != nil {
		panic(err)
	}
}

// 判断用户是否已存在
func IsExist(UnionID string) bool {
	row, err := DB.Query(`SELECT count(*) from ecncer where union_id = $1`, UnionID)
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
	row, err := DB.Query(`SELECT name from ecncer where union_id = $1`, UnionID)
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

func GetAgentTokenFromDB(name string) string {
	row, err := DB.Query(`SELECT token from agent where name = $1`, name)
	if err == nil {
		var token string
		for row.Next() {
			err1 := row.Scan(&token)
			if err1 == nil {
				return token
			}
		}
	}
	return ""
}

func AddAgent(name, token string) {
	_, err := DB.Exec("insert into agent(name,token) values($1,$2)", name, token)
	if err != nil {
		panic(err)
	}
}
