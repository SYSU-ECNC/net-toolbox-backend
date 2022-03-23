package config

import (
	"github.com/spf13/viper"
)

func init() {
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	viper.ReadInConfig()
}

func GetDBHost() string {
	return viper.Get("database.host").(string)
}

func GetDBPort() int {
	return viper.Get("database.port").(int)
}

func GetDBUser() string {
	return viper.Get("database.user").(string)
}

func GetDBPassword() string {
	return viper.Get("database.password").(string)
}

func GetDBName() string {
	return viper.Get("database.dbname").(string)
}

func GetAPPID() string {
	return viper.Get("app.app_id").(string)
}

func GetAPPSecret() string {
	return viper.Get("app.app_secret").(string)
}

func GetServerHost() string {
	return viper.Get("server.host").(string)
}

func GetServerPort() string {
	return viper.Get("server.port").(string)
}

func GetFrontendUrl() string {
	return viper.Get("publicUrl.frontentUrl").(string)
}

func GetLoginUrl() string {
	return viper.Get("publicUrl.loginUrl").(string)
}

func GetCallbackUrl() string {
	return viper.Get("publicUrl.callbackUrl").(string)
}

// // 实现config.yml对应的结构体，以便于读取

// package config

// import (
// 	"fmt"
// 	"io/ioutil"

// 	"gopkg.in/yaml.v2"
// )

// type Config struct {
// 	DATABASE *DATABASE `yaml:"DATABASE"`
// 	APP      *APP      `yaml:"APP"`
// 	SERVER   *SERVER   `yaml:"SERVER"`
// }

// type DATABASE struct {
// 	Host     string `yaml:"host"`
// 	Port     int    `yaml:"port"`
// 	User     string `yaml:"user"`
// 	Password string `yaml:"password"`
// 	Dbname   string `yaml:"dbname"`
// }

// type APP struct {
// 	App_id     string `yaml:"app_id"`
// 	App_secret string `yaml:"app_secret"`
// }

// type SERVER struct {
// 	Host string `yaml:"host"`
// 	Port string `yaml:"port"`
// }

// func getConfig() *Config {
// 	yamlFile, err := ioutil.ReadFile("config/config.yml")
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}
// 	var _config *Config
// 	err = yaml.Unmarshal(yamlFile, &_config)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}
// 	return _config
// }

// func GetDBInfo() (string, int, string, string, string) {
// 	_config := getConfig()
// 	return _config.DATABASE.Host, _config.DATABASE.Port, _config.DATABASE.User, _config.DATABASE.Password, _config.DATABASE.Dbname
// }

// func GetAPPInfo() (string, string) {
// 	_config := getConfig()
// 	return _config.APP.App_id, _config.APP.App_secret
// }

// func GetSERVERInfo() (string, string) {
// 	_config := getConfig()
// 	return _config.SERVER.Host, _config.SERVER.Port
// }
