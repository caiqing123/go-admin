package bootstrap

import (
	"fmt"

	"api/app/models/opera_log"
	"api/pkg/config"
	"api/pkg/logger"
	"api/pkg/queue"
	"api/pkg/redis"
)

// SetupRedis 初始化 Redis
func SetupRedis() {
	if config.GetBool("redis.enable") {
		// 建立 Redis 连接
		redis.ConnectRedis(
			fmt.Sprintf("%v:%v", config.GetString("redis.host"), config.GetString("redis.port")),
			config.GetString("redis.username"),
			config.GetString("redis.password"),
			config.GetInt("redis.database"),
		)

		// 建立 Redis queue 连接
		queue.ConnectRedis(
			fmt.Sprintf("%v:%v", config.GetString("redis.host"), config.GetString("redis.port")),
			config.GetString("redis.username"),
			config.GetInt("redis.queue_database"),
			"queue",
		)
		//添加操作log消费处理
		err := queue.Queue.Consumers(&opera_log.OpLogConsumer{}, "oplog")
		logger.LogIf(err)
	} else {
		queue.ConnectNsq("queue")
		//添加消费处理
		queue.QueueFile.Handler("oplog", queue.Oplog)
		queue.QueueFile.Start()
	}
}
