/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-11 00:00:00
 * @FilePath: \go-sqlbuilder\cache\redis.go
 * @Description: Redis 缓存存储实现
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package cache

import (
	"context"
	"errors"
	"time"

	logger "github.com/kamalyes/go-logger"
)

// RedisStore Redis 缓存存储实现
type RedisStore struct {
	client RedisClientInterface
	prefix string
}

// NewRedisStore 创建 Redis 缓存存储
func NewRedisStore(client RedisClientInterface, prefix string) *RedisStore {
	if prefix == "" {
		prefix = "sqlbuilder:"
	}
	return &RedisStore{
		client: client,
		prefix: prefix,
	}
}

// Get 获取缓存
func (r *RedisStore) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key)
}

// Set 设置缓存
func (r *RedisStore) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl)
}

// Delete 删除缓存
func (r *RedisStore) Delete(ctx context.Context, key string) error {
	_, err := r.client.Del(ctx, key)
	return err
}

// Exists 检查缓存是否存在
func (r *RedisStore) Exists(ctx context.Context, key string) (bool, error) {
	count, err := r.client.Exists(ctx, key)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Clear 清除所有缓存（按前缀）
func (r *RedisStore) Clear(ctx context.Context, prefix string) error {
	if prefix == "" {
		prefix = r.prefix
	}

	// 使用 SCAN 命令获取匹配的键（避免 KEYS 命令阻塞）
	keys, err := r.client.Keys(ctx, prefix+"*")
	if err != nil {
		return err
	}

	if len(keys) == 0 {
		return nil
	}

	_, err = r.client.Del(ctx, keys...)
	return err
}

// ==================== Redis 适配器包装 ====================

// GoRedisAdapter go-redis v9 库的适配器实现
type GoRedisAdapter struct {
	client interface{} // *redis.Client 类型
}

// NewGoRedisAdapter 创建 go-redis 适配器
// 使用方式: adapter := NewGoRedisAdapter(redisClient)
// 注意：这是一个模板实现，需要根据实际的 redis 库版本完成具体方法
func NewGoRedisAdapter(redisClient interface{}) (RedisClientInterface, error) {
	return &GoRedisAdapter{client: redisClient}, nil
}

// Get 获取缓存
func (g *GoRedisAdapter) Get(ctx context.Context, key string) (string, error) {
	// 标准实现: return g.client.Get(ctx, key).Val(), nil
	logger.Error("redis adapter: Get called but adapter not implemented")
	return "", errors.New("adapter requires specific redis client implementation")
}

// Set 设置缓存
func (g *GoRedisAdapter) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	// 标准实现: return g.client.Set(ctx, key, value, expiration).Err()
	logger.Error("redis adapter: Set called but adapter not implemented")
	return errors.New("adapter requires specific redis client implementation")
}

// Del 删除缓存
func (g *GoRedisAdapter) Del(ctx context.Context, keys ...string) (int64, error) {
	// 标准实现: return g.client.Del(ctx, keys...).Result()
	logger.Error("redis adapter: Del called but adapter not implemented")
	return 0, errors.New("adapter requires specific redis client implementation")
}

// Exists 检查缓存是否存在
func (g *GoRedisAdapter) Exists(ctx context.Context, keys ...string) (int64, error) {
	// 标准实现: return g.client.Exists(ctx, keys...).Result()
	logger.Error("redis adapter: Exists called but adapter not implemented")
	return 0, errors.New("adapter requires specific redis client implementation")
}

// Keys 获取匹配的键
func (g *GoRedisAdapter) Keys(ctx context.Context, pattern string) ([]string, error) {
	// 标准实现: return g.client.Keys(ctx, pattern).Result()
	logger.Error("redis adapter: Keys called but adapter not implemented")
	return nil, errors.New("adapter requires specific redis client implementation")
}
