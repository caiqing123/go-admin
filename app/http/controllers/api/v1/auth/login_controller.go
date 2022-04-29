package auth

import (
	"time"

	v1 "api/app/http/controllers/api/v1"
	"api/app/models/login_log"
	"api/app/requests"
	"api/pkg/auth"
	"api/pkg/ip"
	"api/pkg/jwt"
	"api/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"
)

// LoginController 用户控制器
type LoginController struct {
	v1.BaseAPIController
}

// LoginByPhone 手机登录
func (lc *LoginController) LoginByPhone(c *gin.Context) {

	// 1. 验证表单
	request := requests.LoginByPhoneRequest{}
	if ok := requests.Validate(c, &request, requests.LoginByPhone); !ok {
		return
	}

	// 2. 尝试登录
	user, err := auth.LoginByPhone(request.Phone)
	if err != nil {
		// 失败，显示错误提示
		response.Error(c, err, "账号不存在或密码错误")
	} else {
		// 登录成功
		token := jwt.NewJWT().IssueToken(user.GetStringID(), user.Name)

		response.JSON(c, gin.H{
			"token": token,
		})
	}
}

// LoginByPassword 多种方法登录，支持手机号、email 和用户名
func (lc *LoginController) LoginByPassword(c *gin.Context) {
	var status = "2"
	var msg = "登录成功"
	var username = ""
	defer func() {
		LoginLogToDB(c, status, msg, username)
	}()
	// 1. 验证表单
	request := requests.LoginByPasswordRequest{}
	if ok := requests.Validate(c, &request, requests.LoginByPassword); !ok {
		username = request.LoginID
		msg = "验证失败"
		status = "1"
		return
	}

	// 2. 尝试登录
	user, err := auth.Attempt(request.LoginID, request.Password)
	if err != nil {
		username = request.LoginID
		msg = "账号或密码错误"
		status = "1"
		// 失败，显示错误提示
		response.Error(c, err, "登录失败")
	} else {
		username = request.LoginID
		token := jwt.NewJWT().IssueToken(user.GetStringID(), user.Name)
		response.JSON(c, gin.H{
			"token": token,
		})
	}
}

// LoginLogToDB Write log to database
func LoginLogToDB(c *gin.Context, status string, msg string, username string) {
	ua := user_agent.New(c.Request.UserAgent())
	browserName, browserVersion := ua.Browser()
	loginModel := login_log.LoginLog{
		Username:      username,
		Status:        status,
		Ipaddr:        ip.GetClientIP(c),
		LoginLocation: ip.GetLocation(ip.GetClientIP(c)),
		Browser:       browserName + " " + browserVersion,
		Os:            ua.OS(),
		Platform:      ua.Platform(),
		LoginTime:     time.Now(),
		Remark:        c.Request.UserAgent(),
		Msg:           msg,
	}
	loginModel.Create()
}

func (lc *LoginController) Logout(c *gin.Context) {
	currentUser := auth.CurrentUser(c)
	LoginLogToDB(c, "2", "退出成功", currentUser.Name)
	response.Success(c)
}

// RefreshToken 刷新 Access Token
func (lc *LoginController) RefreshToken(c *gin.Context) {
	token, err := jwt.NewJWT().RefreshToken(c)
	if err != nil {
		response.Error(c, err, "令牌刷新失败")
	} else {
		response.JSON(c, gin.H{
			"token": token,
		})
	}
}
