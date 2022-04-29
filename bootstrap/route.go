// Package bootstrap 处理程序初始化逻辑
package bootstrap

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/rakyll/statik/fs"

	"api/app/http/middlewares"
	"api/pkg/ws"
	"api/routes"
	_ "api/statik"
)

// SetupRoute 路由初始化
func SetupRoute(router *gin.Engine) {

	// websocket处理
	go ws.WebsocketManager.Start()
	go ws.WebsocketManager.SendService()
	go ws.WebsocketManager.SendAllService()

	// 注册全局中间件
	registerGlobalMiddleWare(router)

	//  注册 API 路由
	routes.RegisterAPIRoutes(router)

	//  配置 404 路由
	setup404Handler(router)

	//静态文件
	g := router.Group("")
	StaticFileRouter(g)
}

type GinFS struct {
	FS http.FileSystem
}

// Open 打开文件
func (b *GinFS) Open(name string) (http.File, error) {
	return b.FS.Open(name)
}

// Exists 文件是否存在
func (b *GinFS) Exists(prefix string, filepath string) bool {
	if _, err := b.FS.Open(filepath); err != nil {
		return false
	}
	return true
}

func registerGlobalMiddleWare(router *gin.Engine) {
	router.Use(
		middlewares.CorsMiddleware(),
		middlewares.Logger(),
		middlewares.Recovery(),
		middlewares.ForceUA(),
	)

	//静态页面处理
	var StaticFS static.ServeFileSystem
	StaticFS = &GinFS{}
	StaticFS.(*GinFS).FS, _ = fs.New()
	router.Use(static.Serve("/", StaticFS))
}

func setup404Handler(router *gin.Engine) {
	// 处理 404 请求
	router.NoRoute(func(c *gin.Context) {
		// 获取标头信息的 Accept 信息
		acceptString := c.Request.Header.Get("Accept")
		if strings.Contains(acceptString, "text/html") {
			// 静态页面路由直接到首页
			c.Redirect(http.StatusMovedPermanently, "/")
			// 如果是 HTML 的话
			//c.String(http.StatusNotFound, "404")
		} else {
			// 默认返回 JSON
			c.JSON(http.StatusNotFound, gin.H{
				"error_code": 404,
				"message":    "路由未定义，请确认 url 和请求方法是否正确。",
			})
		}
	})
}

func StaticFileRouter(r *gin.RouterGroup) {
	r.Static("/uploads", "./public/uploads")
}
