package user

import (
	"net/http"
	"toolBox/config"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

var store *sessions.CookieStore

func init() {
	authKeyOne := securecookie.GenerateRandomKey(64)
	encryptionKeyOne := securecookie.GenerateRandomKey(32)

	store = sessions.NewCookieStore(
		authKeyOne,
		encryptionKeyOne,
	)

	// session过期时间为2小时
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 120,
		HttpOnly: true,
	}
}

func SetCookie(c *gin.Context, UnionID string) {
	session, _ := store.Get(c.Request, "sessionID")
	session.Values["UnionID"] = UnionID

	// 将set-cookie写入到response里
	session.Save(c.Request, c.Writer)
}

// 生成重定向登录url
func spanLoginUrl() string {
	url := "http://" + config.GetServerHost() + ":" + config.GetServerPort() + "/login"
	return url
}

// 验证cookie是否合法，不合法重定向到登录界面
func VerifyCookie(c *gin.Context) bool {
	session, err := store.Get(c.Request, "sessionID")
	if err != nil || session.Values["UnionID"] == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"redirectUrl": spanLoginUrl()})
		return false
	}
	return true
}

// 使用此函数前默认已经使用了
func GetIDFromCookie(c *gin.Context) string {
	if !VerifyCookie(c) {
		return ""
	}
	session, _ := store.Get(c.Request, "sessionID")
	return session.Values["UnionID"].(string)
}