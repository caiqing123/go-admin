package queue

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/adjust/rmq/v4"
	"github.com/go-redis/redis/v8"

	"api/pkg/logger"
)

var (
	CLI *redis.Client
)

// QueueClient 服务
type QueueClient struct {
	Client rmq.Connection
}

// once 确保全局的 Redis 对象只实例一次
var once sync.Once

// Queue 全局
var Queue *QueueClient

// ConnectRedis 连接 redis 数据库，设置全局的 Queue 对象
func ConnectRedis(address string, username string, db int, name string) {
	once.Do(func() {
		Queue = NewClient(address, username, db, name)
	})
}

func NewClient(address string, password string, db int, name string) *QueueClient {
	CLI = redis.NewClient(&redis.Options{
		Addr:     address,
		DB:       db,
		Password: password,
	})
	_, err := CLI.Ping(CLI.Context()).Result()
	logger.LogIf(err)
	errChan := make(chan error, 10)
	go logErrors(errChan)
	qu := &QueueClient{}
	qu.Client, err = rmq.OpenConnectionWithRedisClient(name, CLI, errChan)
	logger.LogIf(err)
	//清理器
	go clean(qu.Client)

	return qu
}

func clean(d rmq.Connection) {
	cleaner := rmq.NewCleaner(d)
	for range time.Tick(time.Second) {
		_, err := cleaner.Clean()
		if err != nil {
			logger.Error("failed to clean error: " + err.Error())
			continue
		}
	}
}

// 错误处理
func logErrors(errChan <-chan error) {
	for err := range errChan {
		switch err := err.(type) {
		case *rmq.HeartbeatError:
			if err.Count == rmq.HeartbeatErrorLimit {
				logger.Error("heartbeat error (limit): " + err.Error())
			} else {
				logger.Error("heartbeat error: " + err.Error())
			}
		case *rmq.ConsumeError:
			logger.Error("consume error: " + err.Error())
		case *rmq.DeliveryError:
			logger.Error("delivery error: " + err.Error())
		default:
			logger.Error("other error: " + err.Error())
		}
	}
}

//Producers 生产
func (que QueueClient) Producers(task interface{}, name string) error {
	q, err := que.Client.OpenQueue(name)
	taskBytes, err := json.Marshal(task)
	err = q.PublishBytes(taskBytes)
	return err
}

//Consumers 消费
func (que QueueClient) Consumers(taskConsumer rmq.Consumer, name string) error {
	q, err := que.Client.OpenQueue(name)
	err = q.StartConsuming(20, time.Millisecond)
	_, err = q.AddConsumer("task-consumer", taskConsumer)
	return err
}
