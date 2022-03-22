package models

import (
	"database/sql"
	"fmt"
	"toolBox/config"

	_ "github.com/lib/pq"
)

func connectDB() *sql.DB {
	host, port, user, password, dbname := config.GetDBHost(), config.GetDBPort(), config.GetDBUser(), config.GetDBPassword(), config.GetDBName()
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

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
