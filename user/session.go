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

func SetCookie(c *gin.Context, UnionID, name string) {
	session, _ := store.Get(c.Request, "sessionID")
	session.Values["UnionID"] = UnionID
	session.Values["name"] = name

	// 将set-cookie写入到response里
	session.Save(c.Request, c.Writer)
}

// 验证cookie是否合法，不合法重定向到登录界面
func VerifyCookie(c *gin.Context) bool {
	session, err := store.Get(c.Request, "sessionID")
	if err != nil || session.Values["UnionID"] == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"redirectUrl": config.GetLoginUrl()})
		return false
	}
	return true
}

func GetIDFromCookie(c *gin.Context) string {
	if !VerifyCookie(c) {
		return ""
	}
	session, _ := store.Get(c.Request, "sessionID")
	return session.Values["UnionID"].(string)
}

func GetNameFromCookie(c *gin.Context) string {
	if !VerifyCookie(c) {
		return ""
	}
	session, _ := store.Get(c.Request, "sessionID")
	return session.Values["name"].(string)
}
