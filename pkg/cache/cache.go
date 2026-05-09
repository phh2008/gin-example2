package cache

import (
	"encoding/json"
	"errors"
	"time"
)

var ErrNotFound = errors.New("cache: key not found")

type Cache interface {
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

// Get 获取缓存值并反序列化为指定类型
// MemeryCache: 直接类型断言；RedisCache: JSON反序列化
func Get[T any](c Cache, key string) (T, error) {
	var zero T
	val, err := c.Get(key)
	if err != nil {
		return zero, err
	}
	if val == nil {
		return zero, ErrNotFound
	}
	// 直接类型断言（MemeryCache 原始值，或 RedisCache string→[]byte）
	if v, ok := val.(T); ok {
		return v, nil
	}
	// []byte 特殊处理：从 Redis 返回的 string 直接转 []byte，避免 JSON Base64 编解码
	if _, ok := any(zero).([]byte); ok {
		if s, ok := val.(string); ok {
			return any([]byte(s)).(T), nil
		}
	}
	var data []byte
	if s, ok := val.(string); ok {
		data = []byte(s)
	} else {
		data, _ = json.Marshal(val)
	}
	var result T
	if err := json.Unmarshal(data, &result); err != nil {
		return zero, err
	}
	return result, nil
}

// Set 序列化值并写入缓存
// MemeryCache: 直接存储原始值；RedisCache: JSON序列化后存储
func Set[T any](c Cache, key string, value T, expiration time.Duration) error {
	return c.Set(key, value, expiration)
}
