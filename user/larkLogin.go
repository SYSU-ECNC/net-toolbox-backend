package user

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"toolBox/config"
	"toolBox/models"

	"github.com/gin-gonic/gin"
)

func getAppAccessToken() string {
	app_id, app_secret := config.GetAPPID(), config.GetAPPSecret()
	data := url.Values{"app_id": {app_id}, "app_secret": {app_secret}}
	_url := "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal"
	resp, err := http.PostForm(_url, data)
	var match string
	if err == nil {
		data, _ := ioutil.ReadAll(resp.Body)
		r, _ := regexp.Compile(`[a-zA-Z0-9]-..{20,}[a-zA-Z0-9]`)
		match = r.FindString(string(data))
	}

	//避免内存泄露
	defer resp.Body.Close()
	return match
}

func getCode(c *gin.Context) string {
	ret := c.Query("code")
	return ret
}

type Payload struct {
	Grant_type string `json:"grant_type"`
	Code       string `json:"code"`
}

func getUserMessage(_code string) string {
	_url := "https://open.feishu.cn/open-apis/authen/v1/access_token"

	payload := Payload{
		Grant_type: "authorization_code",
		Code:       _code,
	}
	data, _ := json.Marshal(payload)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", _url, bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+getAppAccessToken())
	resp, err := client.Do(req)
	var respBody []byte
	if err == nil {
		defer resp.Body.Close()
		respBody, _ = ioutil.ReadAll(resp.Body)
	}

	return string(respBody)
}

func Login(c *gin.Context) {
	app_id := config.GetAPPID()
	callbackUrl := "http://" + config.GetServerHost() + ":" + config.GetServerPort() + "/callback/path"
	// state := "stateOK"
	redirectUrl := `https://open.feishu.cn/open-apis/authen/v1/index?redirect_uri=` + callbackUrl + `&app_id=` + app_id
	// c.String(http.StatusFound, redirectUrl)
	c.Redirect(http.StatusFound, redirectUrl)
}

type Data struct {
	Access_token       string `json:"access_token"`
	Token_type         string `json:"token_type"`
	Expires_in         int    `json:"expires_in"`
	Name               string `json:"name"`
	En_name            string `json:"en_name"`
	Avatar_url         string `json:"avatar_url"`
	Avatar_thumb       string `json:"avatar_thumb"`
	Avatar_middle      string `json:"avatar_middle"`
	Avatar_big         string `json:"avatar_big"`
	Open_id            string `json:"open_id"`
	Union_id           string `json:"union_id"`
	Email              string `json:"email"`
	User_id            string `json:"user_id"`
	Mobile             string `json:"mobile"`
	Tenant_key         string `json:"tenant_key"`
	Refresh_expires_in int    `json:"refresh_expires_in"`
	Refresh_token      string `json:"refresh_token"`
}

type User struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data Data   `json:"data"`
}

func Callback(c *gin.Context) {
	code := getCode(c)
	userMsg := getUserMessage(code)
	var user User
	err := json.Unmarshal([]byte(userMsg), &user)
	var name string
	var UnionID string
	if err == nil {
		name = user.Data.Name
		// Union_id是飞书账号的唯一标识
		UnionID = user.Data.Union_id
		if !models.IsExist(UnionID) {
			// 用户不存在，在数据库中写入用户信息
			models.DatabaseWrite(name, UnionID)
		}
		// token := generateToken(Union_id)
		SetCookie(c, UnionID)
		redirect_url := config.GetFrontendUrl()
		c.Redirect(http.StatusFound, redirect_url)
	} else {
		// c.Redirect(http.StatusInternalServerError, "callback error!")
		c.Redirect(http.StatusInternalServerError, "/login")
	}
}

// //token加密密钥
// var Key = []byte("")

// func generateToken(Union_id string) string {
// 	// 用户登录有效时间为2小时
// 	expireTime := time.Now().Add(2 * time.Hour)
// 	claims := &jwt.StandardClaims{
// 		Id:        Union_id,
// 		ExpiresAt: expireTime.Unix(),
// 		IssuedAt:  time.Now().Unix(),
// 	}
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	tokenStr, _ := token.SignedString(Key)
// 	return tokenStr
// }

// func VerifyToken(tokenStr string) bool {
// 	if tokenStr == "" {
// 		return false
// 	}
// 	token, _, err := ParseToken(tokenStr)
// 	if err != nil || !token.Valid {
// 		return false
// 	}
// 	return true
// }

// func ParseToken(tokenString string) (*jwt.Token, *jwt.StandardClaims, error) {
// 	Claims := &jwt.StandardClaims{}
// 	token, err := jwt.ParseWithClaims(tokenString, Claims, func(token *jwt.Token) (i interface{}, err error) {
// 		return Key, nil
// 	})
// 	return token, Claims, err
// }
