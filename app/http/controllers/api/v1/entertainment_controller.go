package v1

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"

	"api/pkg/book"
	"api/pkg/book/site"
	"api/pkg/cache"
	"api/pkg/file"
	"api/pkg/hotlist"
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

func init() {
	book.InitSites()
}

// Music 音乐
func (ctrl *BaseAPIController) Music(c *gin.Context) {
	option := c.Query("option")
	name := c.Query("name")
	if len(strings.TrimSpace(name)) == 0 {
		response.NormalVerificationError(c, "缺少参数")
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
		response.NormalVerificationError(c, "无效的参数")
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

// Book 阅读
func (ctrl *BaseAPIController) Book(c *gin.Context) {
	name := c.Query("name")
	if len(strings.TrimSpace(name)) == 0 {
		response.NormalVerificationError(c, "缺少参数")
		return
	}
	var list = make(map[string][]site.ChapterSearchResult)
	for _, s := range site.SitePool {
		if s.Search == nil {
			continue
		}
		result, err := s.Search(name)
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		if len(result) == 0 {
			continue
		}
		list[s.Name] = result
	}
	response.JSON(c, gin.H{
		"data": list,
	})
}

//BookInfo 详情
func (ctrl *BaseAPIController) BookInfo(c *gin.Context) {
	bookUrl := c.Query("url")
	visitorId := c.Query("visitorId")
	if len(strings.TrimSpace(bookUrl)) == 0 {
		response.NormalVerificationError(c, "缺少参数")
		return
	}
	result, err := site.BookInfo(bookUrl)
	if err != nil {
		logger.Error(err.Error())
		response.Abort500(c, "出错")
		return
	}
	if result.DownloadURL == "" {
		fileSrc := "public/uploads/book/" + result.BookName + "_" + visitorId + ".txt"
		if !file.CheckExist(fileSrc) {
			result.DownloadURL = "/uploads/book/" + result.BookName + "_" + visitorId + ".txt"
		}
	}
	response.JSON(c, gin.H{
		"data": result,
	})
}

//News 热门资讯
func (ctrl *BaseAPIController) News(c *gin.Context) {
	hotlist.All(false)
	data := cache.Get(c.Query("type"))
	response.JSON(c, gin.H{
		"data": data,
	})
}

// Download 下载
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
	_, _ = w.Write(data)
	c.Abort()
}
