/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-17 16:10:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-17 16:10:00
 * @FilePath: \go-sqlbuilder\cache\interface.go
 * @Description: 缓存接口定义
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package cache

import (
	"context"
	"time"
)

// Store 缓存存储接口
type Store interface {
	// Get 获取缓存值
	Get(ctx context.Context, key string) (string, error)

	// Set 设置缓存值
	Set(ctx context.Context, key string, value string, ttl time.Duration) error

	// Delete 删除缓存
	Delete(ctx context.Context, key string) error

	// Exists 检查缓存是否存在
	Exists(ctx context.Context, key string) (bool, error)

	// Clear 清除所有缓存（按前缀）
	Clear(ctx context.Context, prefix string) error
}

// Handler 缓存处理器接口（字节级别操作）
type Handler interface {
	Get(key []byte) ([]byte, error)
	Set(key, value []byte) error
	SetWithTTL(key, value []byte, ttl time.Duration) error
	Del(keys ...[]byte) error
	Exists(key []byte) (bool, error)
}
