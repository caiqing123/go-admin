// Package middlewares Gin 中间件
package middlewares

import (
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/chenhg5/collection"

	"api/pkg/cache"
	"api/pkg/http"
	"api/pkg/ip"
	"api/pkg/response"

	"github.com/gin-gonic/gin"
)

var Whitelist = []string{"json", "css", "js", "jpg", "svg", "png", "jpeg", "ico", "gif"}

// ForceUA 中间件，强制请求必须附带 User-Agent 标头
func ForceUA() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 获取 User-Agent 标头信息 请求头验证修改；
		if len(c.Request.Header["User-Agent"]) == 0 {
			response.BadRequest(c, errors.New("User-Agent 标头未找到"), "请求必须附带 User-Agent 标头")
			return
		}

		// 隐藏地址处理
		ipd := ip.GetClientIP(c)
		if strings.Index(c.Request.URL.String(), "975e38ed00fdbeadFS") != -1 && !collection.Collect(Whitelist).Contains(c.Request.URL.String()[strings.LastIndex(c.Request.URL.String(), ".")+1:]) && cache.Get("ma:"+ipd) == "" {
			cache.Set("ma:"+ipd, ipd, time.Minute*5)
			str := "ip:" + ipd + " ip地址:" + ip.GetLocation(ipd) + " url:" + strings.Replace(c.Request.URL.String(), "975e38ed00fdbeadFS", "url", -1) + " UserAgent:" + c.Request.UserAgent()
			_, err := http.Get("https://api2.pushdeer.com/message/push?pushkey=PDU17771TVdkeRyM3UJQJtvWG9fkYRFpprYJ2wSF5&text=" + url.QueryEscape(str))
			if err != nil {
				response.Abort404(c)
				return
			}
		}

		c.Next()
	}
}
