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

// 命名为conf是为了避免和*gin.Context的变量重名
var conf config.Config = config.GetConfig()

func getAppAccessToken() string {
	appId, appSecret := conf.App.Id, conf.App.Secret
	// 此处的key命名为app_id是飞书接口的要求
	data := url.Values{"app_id": {appId}, "app_secret": {appSecret}}
	accessUrl := "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal"
	resp, err := http.PostForm(accessUrl, data)
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

// 此处结构体成员变量命名首字母不大写就会有警告
// 因为结构体转json存储时，结构体成员变量首字母必须大写才能成功输出
// 所以只能写成大驼峰形式
type payload struct {
	GrantType string `json:"grant_type"`
	Code      string `json:"code"`
}

func getUserMessage(userCode string) string {
	url := "https://open.feishu.cn/open-apis/authen/v1/access_token"

	payloadData := payload{
		GrantType: "authorization_code",
		Code:      userCode,
	}

	data, _ := json.Marshal(payloadData)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	// 暂时还是改为每次登录再请求飞书获得AppAccessToken
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
	appId := conf.App.Id
	callbackUrl := conf.PublicUrl.CallbackUrl
	redirectUrl := `https://open.feishu.cn/open-apis/authen/v1/index?redirect_uri=` + callbackUrl + `&app_id=` + appId
	c.Redirect(http.StatusFound, redirectUrl)
}

// 此处结构体成员变量命名首字母不大写就会有警告
// 因为结构体转json存储时，结构体成员变量首字母必须大写才能成功输出
// 所以只能写成大驼峰形式
// UserData是通过用户的code向飞书请求的该用户的信息
type UserData struct {
	AcessToken       string `json:"access_token"`
	TokenType        string `json:"token_type"`
	ExpiresIn        int    `json:"expires_in"`
	Name             string `json:"name"`
	EnName           string `json:"en_name"`
	AvatarUrl        string `json:"avatar_url"`
	AvatarThumb      string `json:"avatar_thumb"`
	AvatarMiddle     string `json:"avatar_middle"`
	AvatarBig        string `json:"avatar_big"`
	OpenId           string `json:"open_id"`
	UnionId          string `json:"union_id"`
	Email            string `json:"email"`
	UserId           string `json:"user_id"`
	Mobile           string `json:"mobile"`
	TenantKey        string `json:"tenant_key"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	RefreshToken     string `json:"refresh_token"`
}

type User struct {
	Code     int      `json:"code"`
	Msg      string   `json:"msg"`
	UserData UserData `json:"data"`
}

func Callback(c *gin.Context) {
	code := getCode(c)

	userMsg := getUserMessage(code)
	var user User
	err := json.Unmarshal([]byte(userMsg), &user)
	var name string
	var unionID string
	if err == nil {
		name = user.UserData.Name
		// unionID是飞书账号的唯一标识
		unionID = user.UserData.UnionId
		if !models.IsExist(unionID) {
			// 用户不存在，在数据库中写入用户信息
			models.AddUser(name, unionID)
		}
		SetCookie(c, unionID, name)
		redirect_url := conf.PublicUrl.FrontentUrl
		c.Redirect(http.StatusFound, redirect_url)
	} else {
		c.Redirect(http.StatusInternalServerError, conf.PublicUrl.LoginUrl)
	}
}
