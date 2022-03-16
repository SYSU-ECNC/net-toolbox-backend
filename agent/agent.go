//尚无进展

package agent

import (
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/dgrijalva/jwt-go.v3"
)

type Claims struct {
	name       string `json:"name"`
	masterIP   string `json:"masterIP"`
	masterPort string `json:"masterPort"`
	jwt.StandardClaims
}

func AgentRegister(c *gin.Context) {
	// name := c.PostForm("name")
	// masterIP := c.PostForm("masterIP")
	// masterPort := c.PostForm(("masterPost"))
	claims := &Claims{
		name:       c.PostForm("name"),
		masterIP:   c.PostForm("masterIP"),
		masterPort: c.PostForm("masterPost"),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expire.Unix(),
			IssuedAt:  time.Now().Unix(),
			// Issuer:    "localhost",
			// Subject:   "user token",
		},
	}
	expire := time.Now().Add(2 * time.Hour)
}
