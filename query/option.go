/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-11 00:00:00
 * @FilePath: \go-sqlbuilder\query\option.go
 * @Description: 查询选项结构
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package query

import "time"

// FindOption 查询选项 - 借鉴 go-core 设计
type FindOption struct {
	BusinessId    string        // 业务 ID
	ShopId        string        // 店铺 ID
	TablePrefix   string        // 表名前缀
	ExcludeField  bool          // 是否排除字段
	ExcludeFields []string      // 排除字段列表
	IncludeFields []string      // 包含字段列表
	NoCache       bool          // 是否不使用缓存
	CacheTTL      time.Duration // 缓存过期时间
}

// NewFindOption 创建查询选项
func NewFindOption() *FindOption {
	return &FindOption{
		CacheTTL: 1 * time.Hour,
	}
}

// WithBusinessId 设置业务 ID
func (fo *FindOption) WithBusinessId(id string) *FindOption {
	fo.BusinessId = id
	return fo
}

// WithShopId 设置店铺 ID
func (fo *FindOption) WithShopId(id string) *FindOption {
	fo.ShopId = id
	return fo
}

// WithTablePrefix 设置表名前缀
func (fo *FindOption) WithTablePrefix(prefix string) *FindOption {
	fo.TablePrefix = prefix
	return fo
}

// WithNoCache 禁用缓存
func (fo *FindOption) WithNoCache() *FindOption {
	fo.NoCache = true
	return fo
}

// WithCacheTTL 设置缓存 TTL
func (fo *FindOption) WithCacheTTL(ttl time.Duration) *FindOption {
	fo.CacheTTL = ttl
	return fo
}

// WithExcludeFields 设置排除字段
func (fo *FindOption) WithExcludeFields(fields ...string) *FindOption {
	fo.ExcludeFields = fields
	return fo
}

// WithIncludeFields 设置包含字段
func (fo *FindOption) WithIncludeFields(fields ...string) *FindOption {
	fo.IncludeFields = fields
	return fo
}
