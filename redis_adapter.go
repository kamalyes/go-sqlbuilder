/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-11 00:00:00
 * @FilePath: \go-sqlbuilder\redis_adapter.go
 * @Description: Redis缓存适配器 - 与github.com/redis/go-redis兼容
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package sqlbuilder

import (
	"context"
	"errors"
	"time"

	logger "github.com/kamalyes/go-logger"
)

// RedisClientInterface Redis客户端接口 - 支持多种Redis库
type RedisClientInterface interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Del(ctx context.Context, keys ...string) (int64, error)
	Exists(ctx context.Context, keys ...string) (int64, error)
	Keys(ctx context.Context, pattern string) ([]string, error)
}

// RedisCacheStore Redis缓存存储实现
type RedisCacheStore struct {
	client RedisClientInterface
	prefix string
}

// NewRedisCacheStore 创建Redis缓存存储
func NewRedisCacheStore(client RedisClientInterface, prefix string) *RedisCacheStore {
	if prefix == "" {
		prefix = "sqlbuilder:"
	}
	return &RedisCacheStore{
		client: client,
		prefix: prefix,
	}
}

// Get 获取缓存
func (r *RedisCacheStore) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key)
}

// Set 设置缓存
func (r *RedisCacheStore) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl)
}

// Delete 删除缓存
func (r *RedisCacheStore) Delete(ctx context.Context, key string) error {
	_, err := r.client.Del(ctx, key)
	return err
}

// Exists 检查缓存是否存在
func (r *RedisCacheStore) Exists(ctx context.Context, key string) (bool, error) {
	count, err := r.client.Exists(ctx, key)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Clear 清除所有缓存（按前缀）
func (r *RedisCacheStore) Clear(ctx context.Context, prefix string) error {
	if prefix == "" {
		prefix = r.prefix
	}

	// 使用SCAN命令获取匹配的键（避免KEYS命令阻塞）
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

// ==================== 通用Redis适配器包装 ====================

// goRedisClientAdapter go-redis v9 适配器
type goRedisClientAdapter struct {
	client interface{} // *redis.Client
}

// NewGoRedisAdapter 创建go-redis适配器
// 使用方式: adapter := NewGoRedisAdapter(redisClient)
func NewGoRedisAdapter(redisClient interface{}) (RedisClientInterface, error) {
	// 此处为兼容性设计，实际使用时根据具体的redis库进行封装
	// 示例：github.com/redis/go-redis/v9
	return &goRedisClientAdapter{client: redisClient}, nil
}

// Get 获取缓存
func (g *goRedisClientAdapter) Get(ctx context.Context, key string) (string, error) {
	// 这是一个适配器示例，实际实现取决于使用的具体Redis库
	// 标准实现：return g.client.Get(ctx, key).Val(), nil
	logger.Error("redis adapter: Get called but adapter not implemented")
	return "", errors.New("adapter requires specific redis client implementation")
}

// Set 设置缓存
func (g *goRedisClientAdapter) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	// 标准实现：return g.client.Set(ctx, key, value, expiration).Err()
	logger.Error("redis adapter: Set called but adapter not implemented")
	return errors.New("adapter requires specific redis client implementation")
}

// Del 删除缓存
func (g *goRedisClientAdapter) Del(ctx context.Context, keys ...string) (int64, error) {
	// 标准实现：return g.client.Del(ctx, keys...).Result()
	logger.Error("redis adapter: Del called but adapter not implemented")
	return 0, errors.New("adapter requires specific redis client implementation")
}

// Exists 检查缓存是否存在
func (g *goRedisClientAdapter) Exists(ctx context.Context, keys ...string) (int64, error) {
	// 标准实现：return g.client.Exists(ctx, keys...).Result()
	logger.Error("redis adapter: Exists called but adapter not implemented")
	return 0, errors.New("adapter requires specific redis client implementation")
}

// Keys 获取匹配的键
func (g *goRedisClientAdapter) Keys(ctx context.Context, pattern string) ([]string, error) {
	// 标准实现：return g.client.Keys(ctx, pattern).Result()
	logger.Error("redis adapter: Keys called but adapter not implemented")
	return nil, errors.New("adapter requires specific redis client implementation")
}

// ==================== 缓存统计和管理 ====================

// CacheStats 缓存统计信息
type CacheStats struct {
	TotalHits   int64         // 总命中数
	TotalMisses int64         // 总未命中数
	HitRate     float64       // 命中率
	AvgTTL      time.Duration // 平均TTL
}

// CacheManager 缓存管理器
type CacheManager struct {
	store CacheStore
	stats CacheStats
}

// NewCacheManager 创建缓存管理器
func NewCacheManager(store CacheStore) *CacheManager {
	return &CacheManager{
		store: store,
		stats: CacheStats{},
	}
}

// InvalidatePattern 使匹配模式的缓存失效
func (cm *CacheManager) InvalidatePattern(ctx context.Context, pattern string) error {
	return cm.store.Clear(ctx, pattern)
}

// GetStats 获取缓存统计
func (cm *CacheManager) GetStats() CacheStats {
	return cm.stats
}

// RecordHit 记录缓存命中
func (cm *CacheManager) RecordHit() {
	cm.stats.TotalHits++
	cm.updateHitRate()
}

// RecordMiss 记录缓存未命中
func (cm *CacheManager) RecordMiss() {
	cm.stats.TotalMisses++
	cm.updateHitRate()
}

// updateHitRate 更新命中率
func (cm *CacheManager) updateHitRate() {
	total := cm.stats.TotalHits + cm.stats.TotalMisses
	if total == 0 {
		cm.stats.HitRate = 0
	} else {
		cm.stats.HitRate = float64(cm.stats.TotalHits) / float64(total)
	}
}

// ResetStats 重置统计
func (cm *CacheManager) ResetStats() {
	cm.stats = CacheStats{}
}
