/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-11 00:00:00
 * @FilePath: \go-sqlbuilder\cache\config.go
 * @Description: 缓存配置结构
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package cache

import "time"

// Config 缓存配置
type Config struct {
	Enabled   bool          // 是否启用缓存
	TTL       time.Duration // 默认缓存过期时间
	KeyPrefix string        // 缓存键前缀
}

// NewConfig 创建默认缓存配置
func NewConfig() *Config {
	return &Config{
		Enabled:   true,
		TTL:       1 * time.Hour,
		KeyPrefix: "sqlbuilder:",
	}
}

// WithEnabled 设置是否启用
func (c *Config) WithEnabled(enabled bool) *Config {
	c.Enabled = enabled
	return c
}

// WithTTL 设置 TTL
func (c *Config) WithTTL(ttl time.Duration) *Config {
	c.TTL = ttl
	return c
}

// WithKeyPrefix 设置键前缀
func (c *Config) WithKeyPrefix(prefix string) *Config {
	c.KeyPrefix = prefix
	return c
}
