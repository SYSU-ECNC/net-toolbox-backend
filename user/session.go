package user

import (
	"net/http"

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

func SetCookie(c *gin.Context, unionID, name string) {
	session, _ := store.Get(c.Request, "sessionID")
	session.Values["unionID"] = unionID
	session.Values["name"] = name

	// 将set-cookie写入到response里
	session.Save(c.Request, c.Writer)
}

// 验证cookie是否合法，不合法则返回401和重定向url
func CheckCookie(c *gin.Context) {
	session, err := store.Get(c.Request, "sessionID")
	if err != nil || session.Values["unionID"] == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"redirectUrl": conf.PublicUrl.LoginUrl})
		c.Abort()
	}
}

func GetIDFromCookie(c *gin.Context) string {
	session, _ := store.Get(c.Request, "sessionID")
	return session.Values["unionID"].(string)
}

func GetNameFromCookie(c *gin.Context) string {
	session, _ := store.Get(c.Request, "sessionID")
	return session.Values["name"].(string)
}
