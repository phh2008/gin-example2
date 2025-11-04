package cache

import "time"

type Storage interface {
	// Set 设置缓存值
	Set(key string, value any, expiration time.Duration) error
	// Get 获取缓存值，如果不存在则返回nil
	Get(key string) (any, error)
	// Delete 删除缓存
	Delete(keys ...string) error
	// Exists 判断是否存在
	Exists(key string) bool
	// Expire 设置过期时间
	Expire(key string, expiration time.Duration) error
	// Clear 清空所有缓存
	Clear() error
}
