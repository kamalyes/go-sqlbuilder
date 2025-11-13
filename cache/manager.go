/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-13 00:00:00
 * @FilePath: \go-sqlbuilder\cache\manager.go
 * @Description: 缓存管理器 - 线程安全的统计和管理
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package cache

import (
	"context"
	"sync/atomic"
	"time"
)

// Stats 缓存统计信息
type Stats struct {
	TotalHits   int64         // 总命中数
	TotalMisses int64         // 总未命中数
	HitRate     float64       // 命中率
	AvgTTL      time.Duration // 平均 TTL
}

// Manager 缓存管理器 - 提供线程安全的统计和管理功能
type Manager struct {
	store       Store
	totalHits   atomic.Int64
	totalMisses atomic.Int64
	hitRate     atomic.Value // float64
}

// NewManager 创建缓存管理器
func NewManager(store Store) *Manager {
	m := &Manager{
		store: store,
	}
	m.hitRate.Store(0.0)
	return m
}

// InvalidatePattern 使匹配模式的缓存失效 (线程安全)
func (cm *Manager) InvalidatePattern(ctx context.Context, pattern string) error {
	return cm.store.Clear(ctx, pattern)
}

// GetStats 获取缓存统计 (线程安全)
func (cm *Manager) GetStats() Stats {
	hits := cm.totalHits.Load()
	misses := cm.totalMisses.Load()
	hitRate := 0.0
	
	total := hits + misses
	if total > 0 {
		hitRate = float64(hits) / float64(total)
	}
	
	return Stats{
		TotalHits:   hits,
		TotalMisses: misses,
		HitRate:     hitRate,
	}
}

// RecordHit 记录缓存命中 (线程安全)
func (cm *Manager) RecordHit() {
	cm.totalHits.Add(1)
}

// RecordMiss 记录缓存未命中 (线程安全)
func (cm *Manager) RecordMiss() {
	cm.totalMisses.Add(1)
}

// ResetStats 重置统计 (线程安全)
func (cm *Manager) ResetStats() {
	cm.totalHits.Store(0)
	cm.totalMisses.Store(0)
	cm.hitRate.Store(0.0)
}
