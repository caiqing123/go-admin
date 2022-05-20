package v1

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"

	"api/pkg/logger"
	"api/pkg/music"
	"api/pkg/music/kugou"
	"api/pkg/music/kuwo"
	"api/pkg/music/migu"
	"api/pkg/music/netease"
	"api/pkg/music/qqmusic"
	"api/pkg/response"
)

type EntertainmentController struct {
	BaseAPIController
}

// Music 音乐
func (ctrl *BaseAPIController) Music(c *gin.Context) {
	option := c.Query("option")
	name := c.Query("name")
	if len(strings.TrimSpace(name)) == 0 {
		response.Abort500(c, "缺少参数")
		return
	}

	var sourcer func(string) ([]music.Result, error)
	switch option {
	case "qq":
		sourcer = qqmusic.QQMusic
	case "netease":
		sourcer = netease.Netease
	case "migu":
		sourcer = migu.Migu
	case "kugou":
		sourcer = kugou.Kugou
	case "kuwo":
		sourcer = kuwo.Kuwo
	default:
		response.Abort500(c, "无效的参数")
		return
	}
	ret, err := sourcer(name)
	if err != nil {
		logger.Error(err.Error())
		response.Abort500(c, "出错")
		return
	}
	response.JSON(c, gin.H{
		"data": ret,
	})
}

func (ctrl *BaseAPIController) Download(c *gin.Context) {
	url1 := c.Query("url")
	name := c.DefaultQuery("name", "download")
	// 转发处理
	//remote, err := url.Parse(url1)
	//if err != nil {
	//	panic(err)
	//}
	//
	//proxy := httputil.NewSingleHostReverseProxy(remote)
	////Define the director func
	////This is a good place to log, for example
	//proxy.Director = func(req *http.Request) {
	//	req.Header = c.Request.Header
	//	req.Host = remote.Host
	//	req.URL.Scheme = remote.Scheme
	//	req.URL.Host = remote.Host
	//	req.URL.Path = remote.Path
	//}
	//
	//proxy.ServeHTTP(c.Writer, c.Request)
	//c.Abort()

	//中转请求url处理
	w, r := c.Writer, c.Request

	client := &http.Client{}
	req, err := http.NewRequest(r.Method, url1, nil)

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	w.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename*=UTF-8''%s", url.PathEscape(name)))
	data, _ := ioutil.ReadAll(resp.Body)
	w.Write(data)
	c.Abort()
}
