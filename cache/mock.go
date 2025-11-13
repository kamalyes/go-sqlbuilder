/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-11 00:00:00
 * @FilePath: \go-sqlbuilder\cache\mock.go
 * @Description: 模拟缓存存储 - 用于开发和测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package cache

import (
	"context"
	"time"

	"github.com/kamalyes/go-sqlbuilder/errors"
)

// MockStore 模拟缓存存储实现（用于开发和测试）
type MockStore struct {
	data map[string]cacheEntry
}

type cacheEntry struct {
	value      string
	expireTime time.Time
}

// NewMockStore 创建模拟缓存存储
func NewMockStore() *MockStore {
	return &MockStore{
		data: make(map[string]cacheEntry),
	}
}

// Get 获取缓存
func (m *MockStore) Get(ctx context.Context, key string) (string, error) {
	entry, exists := m.data[key]
	if !exists {
		return "", errors.NewError(errors.ErrorCodeCacheKeyNotFound, errors.MsgKeyNotFound)
	}

	// 检查是否过期
	if time.Now().After(entry.expireTime) {
		delete(m.data, key)
		return "", errors.NewError(errors.ErrorCodeCacheExpired, errors.MsgKeyNotFound)
	}

	return entry.value, nil
}

// Set 设置缓存
func (m *MockStore) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	m.data[key] = cacheEntry{
		value:      value,
		expireTime: time.Now().Add(ttl),
	}
	return nil
}

// Delete 删除缓存
func (m *MockStore) Delete(ctx context.Context, key string) error {
	delete(m.data, key)
	return nil
}

// Exists 检查缓存是否存在
func (m *MockStore) Exists(ctx context.Context, key string) (bool, error) {
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
func (m *MockStore) Clear(ctx context.Context, prefix string) error {
	for key := range m.data {
		if len(prefix) == 0 || (len(key) > len(prefix) && key[:len(prefix)] == prefix) {
			delete(m.data, key)
		}
	}
	return nil
}

// GetStats 获取缓存统计信息
func (m *MockStore) GetStats() map[string]interface{} {
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
