/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-11 00:00:00
 * @FilePath: \go-sqlbuilder\cache.go
 * @Description: Redis缓存层集成 - 自动WithTTL管理
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package sqlbuilder

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"time"
)

// CacheConfig Redis缓存配置
type CacheConfig struct {
	Enabled   bool          // 是否启用缓存
	TTL       time.Duration // 缓存过期时间
	KeyPrefix string        // 缓存键前缀
}

// CacheStore 缓存存储接口
type CacheStore interface {
	// Get 获取缓存
	Get(ctx context.Context, key string) (string, error)

	// Set 设置缓存
	Set(ctx context.Context, key string, value string, ttl time.Duration) error

	// Delete 删除缓存
	Delete(ctx context.Context, key string) error

	// Exists 检查缓存是否存在
	Exists(ctx context.Context, key string) (bool, error)

	// Clear 清除所有缓存（按前缀）
	Clear(ctx context.Context, prefix string) error
}

// CachedBuilder 带缓存的查询构建器
type CachedBuilder struct {
	*Builder
	cacheStore  CacheStore
	cacheConfig CacheConfig
	queryHash   string
}

// NewCachedBuilder 创建带缓存的构建器
func NewCachedBuilder(dbInstance interface{}, cache CacheStore, config CacheConfig) (*CachedBuilder, error) {
	if config.TTL == 0 {
		config.TTL = 1 * time.Hour // 默认1小时
	}
	if config.KeyPrefix == "" {
		config.KeyPrefix = "sqlbuilder:"
	}

	builder, err := New(dbInstance)
	if err != nil {
		return nil, err
	}

	return &CachedBuilder{
		Builder:     builder,
		cacheStore:  cache,
		cacheConfig: config,
	}, nil
}

// WithTTL 设置缓存TTL
func (cb *CachedBuilder) WithTTL(ttl time.Duration) *CachedBuilder {
	cb.cacheConfig.TTL = ttl
	return cb
}

// DisableCache 禁用此次查询的缓存
func (cb *CachedBuilder) DisableCache() *CachedBuilder {
	cb.cacheConfig.Enabled = false
	return cb
}

// EnableCache 启用此次查询的缓存
func (cb *CachedBuilder) EnableCache() *CachedBuilder {
	cb.cacheConfig.Enabled = true
	return cb
}

// ClearCache 清除缓存
func (cb *CachedBuilder) ClearCache() error {
	if cb.cacheStore == nil {
		return fmt.Errorf("cache store not configured")
	}
	return cb.cacheStore.Clear(cb.ctx, cb.cacheConfig.KeyPrefix)
}

// generateCacheKey 生成缓存键
func (cb *CachedBuilder) generateCacheKey() string {
	sql, args := cb.ToSQL()

	// 创建SQL和参数的哈希
	key := fmt.Sprintf("%s%s_%v", cb.cacheConfig.KeyPrefix, sql, args)

	hash := md5.Sum([]byte(key))
	return fmt.Sprintf("%s%x", cb.cacheConfig.KeyPrefix, hash)
}

// GetCached 获取结果（带缓存）
func (cb *CachedBuilder) GetCached(dest interface{}) error {
	if !cb.cacheConfig.Enabled || cb.cacheStore == nil {
		return cb.Get(dest)
	}

	cacheKey := cb.generateCacheKey()

	// 尝试从缓存获取
	cached, err := cb.cacheStore.Get(cb.ctx, cacheKey)
	if err == nil && cached != "" {
		// 缓存命中
		return json.Unmarshal([]byte(cached), dest)
	}

	// 从数据库获取
	if err := cb.Get(dest); err != nil {
		return err
	}

	// 存入缓存
	if data, err := json.Marshal(dest); err == nil {
		_ = cb.cacheStore.Set(cb.ctx, cacheKey, string(data), cb.cacheConfig.TTL)
	}

	return nil
}

// FirstCached 获取第一条记录（带缓存）
func (cb *CachedBuilder) FirstCached(dest interface{}) error {
	if !cb.cacheConfig.Enabled || cb.cacheStore == nil {
		return cb.First(dest)
	}

	cacheKey := cb.generateCacheKey()

	// 尝试从缓存获取
	cached, err := cb.cacheStore.Get(cb.ctx, cacheKey)
	if err == nil && cached != "" {
		return json.Unmarshal([]byte(cached), dest)
	}

	// 从数据库获取
	if err := cb.First(dest); err != nil {
		return err
	}

	// 存入缓存
	if data, err := json.Marshal(dest); err == nil {
		_ = cb.cacheStore.Set(cb.ctx, cacheKey, string(data), cb.cacheConfig.TTL)
	}

	return nil
}

// CountCached 获取计数（带缓存）
func (cb *CachedBuilder) CountCached() (int64, error) {
	if !cb.cacheConfig.Enabled || cb.cacheStore == nil {
		return cb.Count()
	}

	cacheKey := cb.generateCacheKey() + ":count"

	// 尝试从缓存获取
	cached, err := cb.cacheStore.Get(cb.ctx, cacheKey)
	if err == nil && cached != "" {
		var count int64
		if err := json.Unmarshal([]byte(cached), &count); err == nil {
			return count, nil
		}
	}

	// 从数据库获取
	count, err := cb.Count()
	if err != nil {
		return 0, err
	}

	// 存入缓存
	if data, err := json.Marshal(count); err == nil {
		_ = cb.cacheStore.Set(cb.ctx, cacheKey, string(data), cb.cacheConfig.TTL)
	}

	return count, nil
}

// InvalidateCache 使特定查询的缓存失效
func (cb *CachedBuilder) InvalidateCache() error {
	if cb.cacheStore == nil {
		return fmt.Errorf("cache store not configured")
	}

	cacheKey := cb.generateCacheKey()
	return cb.cacheStore.Delete(cb.ctx, cacheKey)
}

// ==================== MockCacheStore 用于测试的模拟缓存实现 ====================

// MockCacheStore 模拟缓存存储（用于开发和测试）
type MockCacheStore struct {
	data map[string]cacheEntry
}

type cacheEntry struct {
	value      string
	expireTime time.Time
}

// NewMockCacheStore 创建模拟缓存存储
func NewMockCacheStore() *MockCacheStore {
	return &MockCacheStore{
		data: make(map[string]cacheEntry),
	}
}

// Get 获取缓存
func (m *MockCacheStore) Get(ctx context.Context, key string) (string, error) {
	entry, exists := m.data[key]
	if !exists {
		return "", fmt.Errorf("key not found")
	}

	// 检查是否过期
	if time.Now().After(entry.expireTime) {
		delete(m.data, key)
		return "", fmt.Errorf("key expired")
	}

	return entry.value, nil
}

// Set 设置缓存
func (m *MockCacheStore) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	m.data[key] = cacheEntry{
		value:      value,
		expireTime: time.Now().Add(ttl),
	}
	return nil
}

// Delete 删除缓存
func (m *MockCacheStore) Delete(ctx context.Context, key string) error {
	delete(m.data, key)
	return nil
}

// Exists 检查缓存是否存在
func (m *MockCacheStore) Exists(ctx context.Context, key string) (bool, error) {
	entry, exists := m.data[key]
	if !exists {
		return false, nil
	}

	// 检查是否过期
	if time.Now().After(entry.expireTime) {
		delete(m.data, key)
		return false, nil
	}

	return true, nil
}

// Clear 清除所有缓存（按前缀）
func (m *MockCacheStore) Clear(ctx context.Context, prefix string) error {
	for key := range m.data {
		if len(prefix) == 0 || (len(key) > len(prefix) && key[:len(prefix)] == prefix) {
			delete(m.data, key)
		}
	}
	return nil
}

// GetStats 获取缓存统计
func (m *MockCacheStore) GetStats() map[string]interface{} {
	validCount := 0
	for _, entry := range m.data {
		if time.Now().Before(entry.expireTime) {
			validCount++
		}
	}

	return map[string]interface{}{
		"total":   len(m.data),
		"valid":   validCount,
		"expired": len(m.data) - validCount,
	}
}
