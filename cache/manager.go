/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-11 00:00:00
 * @FilePath: \go-sqlbuilder\cache\manager.go
 * @Description: 缓存管理器 - 统计和管理缓存
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package cache

import (
	"context"
	"time"
)

// Stats 缓存统计信息
type Stats struct {
	TotalHits   int64         // 总命中数
	TotalMisses int64         // 总未命中数
	HitRate     float64       // 命中率
	AvgTTL      time.Duration // 平均 TTL
}

// Manager 缓存管理器 - 提供统计和管理功能
type Manager struct {
	store Store
	stats Stats
}

// NewManager 创建缓存管理器
func NewManager(store Store) *Manager {
	return &Manager{
		store: store,
		stats: Stats{},
	}
}

// InvalidatePattern 使匹配模式的缓存失效
func (cm *Manager) InvalidatePattern(ctx context.Context, pattern string) error {
	return cm.store.Clear(ctx, pattern)
}

// GetStats 获取缓存统计
func (cm *Manager) GetStats() Stats {
	return cm.stats
}

// RecordHit 记录缓存命中
func (cm *Manager) RecordHit() {
	cm.stats.TotalHits++
	cm.updateHitRate()
}

// RecordMiss 记录缓存未命中
func (cm *Manager) RecordMiss() {
	cm.stats.TotalMisses++
	cm.updateHitRate()
}

// updateHitRate 更新命中率
func (cm *Manager) updateHitRate() {
	total := cm.stats.TotalHits + cm.stats.TotalMisses
	if total == 0 {
		cm.stats.HitRate = 0
	} else {
		cm.stats.HitRate = float64(cm.stats.TotalHits) / float64(total)
	}
}

// ResetStats 重置统计
func (cm *Manager) ResetStats() {
	cm.stats = Stats{}
}
