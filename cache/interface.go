/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-11 00:00:00
 * @FilePath: \go-sqlbuilder\cache\interface.go
 * @Description: 缓存存储接口定义
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package cache

import (
	"context"
	"time"
)

// Store 缓存存储接口 - 定义所有缓存操作的合约
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

// RedisClientInterface Redis 客户端接口 - 支持多种 Redis 库的适配
type RedisClientInterface interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Del(ctx context.Context, keys ...string) (int64, error)
	Exists(ctx context.Context, keys ...string) (int64, error)
	Keys(ctx context.Context, pattern string) ([]string, error)
}
