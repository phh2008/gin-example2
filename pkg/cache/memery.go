package cache

import (
	"time"

	"github.com/jellydator/ttlcache/v3"
)

type MemeryStorage struct {
	pool *ttlcache.Cache[string, any]
}

func NewMemeryStorage() *MemeryStorage {
	return &MemeryStorage{
		pool: ttlcache.New[string, any](),
	}
}

// Set 设置缓存值
func (a *MemeryStorage) Set(key string, value any, expiration time.Duration) error {

	// TODO
	return nil
}

// Get 获取缓存值，如果不存在则返回nil
func (a *MemeryStorage) Get(key string) (any, error) {
	// TODO
	return nil, nil
}

// Delete 删除缓存
func (a *MemeryStorage) Delete(keys ...string) error {

	return nil
}

// Exists 判断是否存在
func (a *MemeryStorage) Exists(key string) bool {

	return false
}

// Expire 设置过期时间
func (a *MemeryStorage) Expire(key string, expiration time.Duration) error {

	return nil
}

// Clear 清空所有缓存
func (a *MemeryStorage) Clear() error {

	return nil
}
