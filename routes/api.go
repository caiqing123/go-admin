// Package routes 注册路由
package routes

import (
	"github.com/gin-gonic/gin"

	controllers "api/app/http/controllers/api/v1"
	"api/app/http/controllers/api/v1/auth"
	"api/app/http/middlewares"
	"api/pkg/gogetssl"
	"api/pkg/limiter"
	"api/pkg/ws"
)

// RegisterAPIRoutes 注册 API 相关路由
func RegisterAPIRoutes(r *gin.Engine) {

	var v1 *gin.RouterGroup
	v1 = r.Group("/api/v1")

	wss := r.Group("")
	{
		wss.GET("/ws/:channel", ws.WebsocketManager.WsClient)
		wss.GET("/wslogout/:channel", ws.WebsocketManager.UnWsClient)
	}
	v1.GET("/demo", func(context *gin.Context) {
	})

	// 全局限流中间件：这里是所有 API （根据 IP）请求加起来。
	v1.Use(middlewares.LimitIP(limiter.Limiter("2000-M")))
	{

		noAuthGroup := v1.Group("", middlewares.GuestJWT())
		{
			// 登录
			lgc := new(auth.LoginController)
			noAuthGroup.POST("/login/using-phone", lgc.LoginByPhone)
			noAuthGroup.POST("/login/using-password", lgc.LoginByPassword)
			noAuthGroup.POST("/login/refresh-token", lgc.RefreshToken)

			vcc := new(auth.VerifyCodeController)
			noAuthGroup.GET("/captcha", middlewares.LimitPerRoute(limiter.Limiter("60-M")), vcc.ShowCaptcha)

		}

		// 需要token
		authGroup := v1.Group("", middlewares.AuthJWT())
		{

			lg := new(auth.LoginController)
			authGroup.POST("/logout", lg.Logout)

			uc := new(controllers.UsersController)
			authGroup.GET("/get-info", uc.CurrentUser)
			authGroup.GET("/user/profile", uc.GetProfile)
			authGroup.PUT("/user", uc.UpdateProfile)
			authGroup.POST("/user", uc.Add)
			authGroup.PUT("/user/save", uc.Update)
			authGroup.PUT("/user/password", uc.UpdatePassword)
			authGroup.POST("/user/avatar", uc.UpdateAvatar)
			authGroup.GET("/user", uc.Index)
			authGroup.GET("/user/info", uc.Info)
			authGroup.GET("/user/pwd/reset", uc.ResetByPassword)
			authGroup.DELETE("/user", uc.Delete)

			mu := new(controllers.MenusController)
			authGroup.GET("/menu-role", mu.GetMenuRole)
			authGroup.GET("/role-menu-tree-select", mu.GetMenuTreeSelect)
			authGroup.GET("/menu", mu.Index)
			authGroup.POST("/menu", mu.Add)
			authGroup.PUT("/menu", mu.Update)
			authGroup.GET("/menu/info", mu.GetMenus)
			authGroup.DELETE("/menu", mu.Delete)

			cf := new(controllers.ConfigsController)
			authGroup.GET("/app-config", cf.GetApp)
			authGroup.GET("/configKey", cf.ConfigKey)
			authGroup.GET("/set-config", cf.GetAll)
			authGroup.PUT("/set-config", cf.SetConfig)

			dcd := new(controllers.DictDataController)
			authGroup.GET("/dict-data/option-select", dcd.GetData)
			authGroup.GET("/dict/data", dcd.Index)
			authGroup.PUT("/dict/data", dcd.Update)
			authGroup.POST("/dict/data", dcd.Add)
			authGroup.DELETE("/dict/data", dcd.Delete)

			ro := new(controllers.RolesController)
			authGroup.GET("/role", ro.GetRole)
			authGroup.POST("/role", ro.Add)
			authGroup.PUT("/role", ro.Update)
			authGroup.GET("/roles", ro.GetRoles)
			authGroup.DELETE("/role", ro.Delete)

			file := new(controllers.FileController)
			authGroup.POST("/public/uploadFile", file.UploadFile)

			dit := new(controllers.DictTypeController)
			authGroup.GET("/dict/type", dit.Index)
			authGroup.POST("/dict/type", dit.Add)
			authGroup.PUT("/dict/type", dit.Update)
			authGroup.DELETE("/dict/type", dit.Delete)

			log := new(controllers.LogController)
			authGroup.GET("/login_log", log.LoginLog)
			authGroup.GET("/opera_log", log.OperaLog)
			authGroup.DELETE("/opera_log/clean", log.Clean)

			serve := new(controllers.ServerController)
			authGroup.GET("/server-monitor", serve.ServerInfo)
			authGroup.GET("/download-log", serve.DownloadLog)

			job := new(controllers.JobController)
			authGroup.GET("/job", job.Index)
			authGroup.POST("/job", job.Add)
			authGroup.PUT("/job", job.Update)
			authGroup.DELETE("/job", job.Delete)
			authGroup.GET("/job/start", job.StartJob)
			authGroup.GET("/job/remove", job.RemoveJob)
		}

		sslGroup := v1.Group("", middlewares.AuthSsl(), middlewares.AuthJWT())
		{
			goSsl := new(controllers.GoGetSllController)
			sslGroup.GET("/go_get_ssl", func(context *gin.Context) {
				gogetssl.VerificationToken()
			}, goSsl.Index)
			sslGroup.POST("/go_get_ssl", func(context *gin.Context) {
				gogetssl.VerificationToken()
			}, goSsl.Operation)
		}
	}
}
