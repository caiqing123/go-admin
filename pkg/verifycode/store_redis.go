package verifycode

import (
	"time"

	"api/pkg/cache"
	"api/pkg/config"
)

// RedisStore 实现 verifycode.Store interface
type RedisStore struct {
	RedisClient *cache.CacheService
	KeyPrefix   string
}

// Set 实现 verifycode.Store interface 的 Set 方法
func (s *RedisStore) Set(key string, value string) bool {
	ExpireTime := time.Second * time.Duration(config.GetInt("verifycode.expire_time")) * 60
	s.RedisClient.Store.Set(s.KeyPrefix+key, value, ExpireTime)
	return true
}

// Get 实现 verifycode.Store interface 的 Get 方法
func (s *RedisStore) Get(key string, clear bool) (value string) {
	key = s.KeyPrefix + key
	val := s.RedisClient.Store.Get(key)
	if clear {
		s.RedisClient.Store.Forget(key)
	}
	return val
}

// Verify 实现 verifycode.Store interface 的 Verify 方法
func (s *RedisStore) Verify(key, answer string, clear bool) bool {
	v := s.Get(key, clear)
	return v == answer
}
