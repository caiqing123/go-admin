// Package bootstrap 启动程序功能
package bootstrap

import (
	"fmt"

	"api/pkg/cache"
	"api/pkg/config"
)

// SetupCache 缓存
func SetupCache() {
	var rds interface{ cache.Store }
	cached := config.GetString("redis.default_cache")
	if !config.GetBool("redis.enable") && cached != "file" {
		cached = "memory"
	}
	switch cached {
	case "redis":
		// 初始化缓存专用的 redis client, 使用专属缓存 DB
		rds = cache.NewRedisStore(
			fmt.Sprintf("%v:%v", config.GetString("redis.host"), config.GetString("redis.port")),
			config.GetString("redis.username"),
			config.GetString("redis.password"),
			config.GetInt("redis.database_cache"),
		)
	case "memory":
		rds = cache.MemoryNew()
	case "file":
		rds = cache.New(
			"storage/cache",
		)
	default:
		rds = cache.New(
			"storage/cache",
		)
	}

	cache.InitWithCacheStore(rds)
}
