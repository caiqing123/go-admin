package captcha

import (
	"time"

	"api/pkg/cache"
	"api/pkg/config"
)

// CacheStore 实现 base64Captcha.Store interface
type CacheStore struct {
	CacheClient *cache.CacheService
	KeyPrefix   string
}

// Set 实现 base64Captcha.Store interface 的 Set 方法
func (s *CacheStore) Set(key string, value string) error {
	ExpireTime := time.Second * time.Duration(config.GetInt("captcha.expire_time")) * 60
	s.CacheClient.Store.Set(s.KeyPrefix+key, value, ExpireTime)
	return nil
}

// Get 实现 base64Captcha.Store interface 的 Get 方法
func (s *CacheStore) Get(key string, clear bool) string {
	key = s.KeyPrefix + key
	val := s.CacheClient.Store.Get(key)
	if clear {
		s.CacheClient.Store.Forget(key)
	}
	return val
}

// Verify 实现 base64Captcha.Store interface 的 Verify 方法
func (s *CacheStore) Verify(key, answer string, clear bool) bool {
	v := s.Get(key, clear)
	return v == answer
}
