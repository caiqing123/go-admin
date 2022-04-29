package config

import (
	"api/pkg/config"
)

func init() {

	config.Add("redis", func() map[string]interface{} {
		return map[string]interface{}{

			"enable":   config.Env("REDIS_SERVE", false),
			"host":     config.Env("REDIS_HOST", "127.0.0.1"),
			"port":     config.Env("REDIS_PORT", "6379"),
			"password": config.Env("REDIS_PASSWORD", ""),

			// 业务类存储使用 1 (图片验证码、短信验证码、会话)
			"database": config.Env("REDIS_MAIN_DB", 1),

			// 缓存 cache 包使用 0 ，缓存清空理应当不影响业务
			"database_cache": config.Env("REDIS_CACHE_DB", 0),

			//缓存方式 redis file memory
			"default_cache": config.Env("CACHE_DRIVER", "redis"),

			// 队列类存储使用 2
			"queue_database": config.Env("REDIS_QUEUE_DB", 2),
		}
	})
}
