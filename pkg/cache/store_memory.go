package cache

import (
	"time"

	"github.com/patrickmn/go-cache"

	"api/pkg/config"
	"api/pkg/logger"
)

// MemoryStore 内存驱动
type MemoryStore struct {
	MemoryClient *cache.Cache
	KeyPrefix    string
}

func MemoryNew() *MemoryStore {
	my := &MemoryStore{}
	//创建一个默认过期时间为 5 分钟的缓存，每 10 分钟清除一次
	my.MemoryClient = cache.New(5*time.Minute, 10*time.Minute)
	my.KeyPrefix = config.GetString("app.name") + ":cache:"
	return my
}

func (m *MemoryStore) Increment(parameters ...interface{}) {
	switch len(parameters) {
	case 1:
		key := parameters[0].(string)
		if err := m.MemoryClient.Increment(m.KeyPrefix+key, 1); err != nil {
			logger.ErrorString("Memory", "Increment", err.Error())
		}
	case 2:
		key := parameters[0].(string)
		value := parameters[1].(int64)
		if err := m.MemoryClient.Increment(m.KeyPrefix+key, value); err != nil {
			logger.ErrorString("Memory", "Increment", err.Error())
		}
	default:
		logger.ErrorString("Memory", "Increment", "参数过多")
	}
}

func (m *MemoryStore) Decrement(parameters ...interface{}) {
	switch len(parameters) {
	case 1:
		key := parameters[0].(string)
		if err := m.MemoryClient.Decrement(m.KeyPrefix+key, 1); err != nil {
			logger.ErrorString("Memory", "Increment", err.Error())
		}
	case 2:
		key := parameters[0].(string)
		value := parameters[1].(int64)
		if err := m.MemoryClient.Decrement(m.KeyPrefix+key, value); err != nil {
			logger.ErrorString("Memory", "Increment", err.Error())
		}
	default:
		logger.ErrorString("Memory", "Increment", "参数过多")
	}
}

func (m *MemoryStore) Set(key string, value string, expireTime time.Duration) {
	m.MemoryClient.Set(m.KeyPrefix+key, value, expireTime)
}

func (m *MemoryStore) Get(key string) string {
	data, found := m.MemoryClient.Get(m.KeyPrefix + key)
	if found {
		return data.(string)
	}
	return ""
}

func (m *MemoryStore) Has(key string) bool {
	return m.Get(m.KeyPrefix+key) != ""
}

func (m *MemoryStore) Forget(key string) {
	m.MemoryClient.Delete(m.KeyPrefix + key)
}

func (m *MemoryStore) Forever(key string, value string) {
	m.MemoryClient.Set(m.KeyPrefix+key, value, cache.NoExpiration)
}

func (m *MemoryStore) Flush() {
	m.MemoryClient.Flush()

}

func (m *MemoryStore) IsAlive() error {
	return nil
}
