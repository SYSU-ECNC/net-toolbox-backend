package config

import (
	"github.com/spf13/viper"
)

var c Config

func GetConfig() Config {
	return c
}

func init() {
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	viper.ReadInConfig()

	c.Database.Host = viper.Get("database.host").(string)
	c.Database.Port = viper.Get("database.port").(int)
	c.Database.User = viper.Get("database.user").(string)
	c.Database.Password = viper.Get("database.password").(string)
	c.Database.Dbname = viper.Get("database.dbname").(string)

	c.App.Id = viper.Get("app.app_id").(string)
	c.App.Secret = viper.Get("app.app_secret").(string)

	c.Server.Host = viper.Get("server.host").(string)
	c.Server.Port = viper.Get("server.port").(string)

	c.PublicUrl.CallbackUrl = viper.Get("publicUrl.frontentUrl").(string)
	c.PublicUrl.LoginUrl = viper.Get("publicUrl.loginUrl").(string)
	c.PublicUrl.FrontentUrl = viper.Get("publicUrl.callbackUrl").(string)
}

type Config struct {
	Database  Database
	App       App
	Server    Server
	PublicUrl PublicUrl
}

type Database struct {
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
}

type App struct {
	Id     string
	Secret string
}

type Server struct {
	Host string
	Port string
}

type PublicUrl struct {
	CallbackUrl string
	LoginUrl    string
	FrontentUrl string
}
