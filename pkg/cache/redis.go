package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(client *redis.Client) *RedisStorage {
	return &RedisStorage{
		client: client,
	}
}

// Set 设置缓存值
func (a *RedisStorage) Set(key string, value any, ttl time.Duration) error {
	return a.client.Set(context.TODO(), key, value, ttl).Err()
}

// Get 获取缓存值，如果不存在则返回nil
func (a *RedisStorage) Get(key string) (any, error) {
	val, err := a.client.Get(context.TODO(), key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return val, nil
}

// Delete 删除缓存
func (a *RedisStorage) Delete(keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	return a.client.Del(context.TODO(), keys...).Err()
}

// Exists 判断是否存在
func (a *RedisStorage) Exists(key string) bool {
	result, err := a.client.Exists(context.TODO(), key).Result()
	if err != nil {
		return false
	}
	return result > 0
}

// Expire 设置过期时间
func (a *RedisStorage) Expire(key string, ttl time.Duration) error {
	return a.client.Expire(context.TODO(), key, ttl).Err()
}

// Keys 获取匹配模式的所有键
func (a *RedisStorage) Keys(pattern string) ([]string, error) {
	ctx := context.TODO()
	var (
		cursor uint64
		result []string
	)
	for {
		keys, next, err := a.client.Scan(ctx, cursor, pattern, 1000).Result()
		if err != nil {
			return nil, err
		}
		if len(keys) > 0 {
			result = append(result, keys...)
		}
		cursor = next
		if cursor == 0 {
			break
		}
	}
	return result, nil
}

// Clear 清空所有缓存
func (a *RedisStorage) Clear() error {
	var cursor uint64
	ctx := context.TODO()
	for {
		keys, next, err := a.client.Scan(ctx, cursor, "*", 1000).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			// Use UNLINK for async non-blocking deletion
			if err := a.client.Unlink(ctx, keys...).Err(); err != nil {
				return err
			}
		}
		cursor = next
		if cursor == 0 {
			break
		}
	}
	return nil
}
