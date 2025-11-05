package cache

import (
	"strings"
	"time"

	"github.com/jellydator/ttlcache/v3"
)

type MemeryStorage struct {
	pool *ttlcache.Cache[string, any]
}

func NewMemeryStorage() Storage {
	return &MemeryStorage{
		pool: ttlcache.New[string, any](),
	}
}

// Set 设置缓存值
func (a *MemeryStorage) Set(key string, value any, ttl time.Duration) error {
	a.pool.Set(key, value, ttl)
	return nil
}

// Get 获取缓存值，如果不存在则返回nil
func (a *MemeryStorage) Get(key string) (any, error) {
	item := a.pool.Get(key)
	if item != nil {
		return item.Value(), nil
	}
	return nil, nil
}

// Delete 删除缓存
func (a *MemeryStorage) Delete(keys ...string) error {
	for _, key := range keys {
		a.pool.Delete(key)
	}
	return nil
}

// Exists 判断是否存在
func (a *MemeryStorage) Exists(key string) bool {
	item := a.pool.Get(key)
	return item != nil && item.Value() != nil
}

// Expire 设置过期时间
func (a *MemeryStorage) Expire(key string, ttl time.Duration) error {
	item := a.pool.Get(key)
	if item != nil && item.Value() != nil {
		a.pool.Set(key, item.Value(), ttl)
	}
	return nil
}

func (a *MemeryStorage) Keys(pattern string) ([]string, error) {
	var keys []string
	pattern = strings.TrimRight(pattern, "*")
	for _, key := range a.pool.Keys() {
		if strings.HasPrefix(key, pattern) {
			keys = append(keys, key)
		}
	}
	return keys, nil
}

// Clear 清空所有缓存
func (a *MemeryStorage) Clear() error {
	a.pool.DeleteAll()
	return nil
}
