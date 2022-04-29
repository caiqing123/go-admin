// Package limiter 处理限流逻辑
package limiter

import (
	"strings"

	"github.com/gin-gonic/gin"
	limiterlib "github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"

	"api/pkg/config"
	"api/pkg/logger"
	"api/pkg/redis"
)

// GetKeyIP 获取 Limitor 的 Key，IP
func GetKeyIP(c *gin.Context) string {
	return c.ClientIP()
}

// GetKeyRouteWithIP Limitor 的 Key，路由+IP，针对单个路由做限流
func GetKeyRouteWithIP(c *gin.Context) string {
	return routeToKeyString(c.FullPath()) + c.ClientIP()
}

// Limiter 获取限流实例
func Limiter(formatted string) *limiterlib.Limiter {

	// 实例化依赖的 limiter 包的 limiter.Rate 对象
	rate, err := limiterlib.NewRateFromFormatted(formatted)
	if err != nil {
		logger.LogIf(err)
	}

	// 初始化存储，使用我们程序里共用的
	var store limiterlib.Store
	if config.GetBool("redis.enable") {
		store, err = sredis.NewStoreWithOptions(redis.Redis.Client, limiterlib.StoreOptions{
			// 为 limiter 设置前缀，保持 redis 里数据的整洁
			Prefix: config.GetString("app.name") + ":limiter",
		})
		if err != nil {
			logger.LogIf(err)
		}
	} else {
		store = memory.NewStoreWithOptions(limiterlib.StoreOptions{
			// 为 limiter 设置前缀，保持 redis 里数据的整洁
			Prefix: config.GetString("app.name") + ":limiter",
		})
	}

	// 使用上面的初始化的 limiter.Rate 对象和存储对象
	limiterObj := limiterlib.New(store, rate)

	return limiterObj
}

// routeToKeyString 辅助方法，将 URL 中的 / 格式为 -
func routeToKeyString(routeName string) string {
	routeName = strings.ReplaceAll(routeName, "/", "-")
	routeName = strings.ReplaceAll(routeName, ":", "_")
	return routeName
}
