/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-11 14:41:55
 * @FilePath: \go-sqlbuilder\advanced_query.go
 * @Description: [DEPRECATED] 高级查询参数向后兼容适配器
 *               请使用 github.com/kamalyes/go-sqlbuilder/query 包
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package sqlbuilder

import (
	"time"

	"github.com/kamalyes/go-sqlbuilder/query"
)

// ==================== 向后兼容适配 ====================
// 以下类型和函数为向后兼容，推荐使用 query 包中的对应类型

// Deprecated: 使用 query.Operator 替代
type FilterOperator = query.Operator

// Deprecated: 使用 query.OP_* 常量替代
const (
	OP_EQ          FilterOperator = query.OP_EQ          // 等于
	OP_NEQ         FilterOperator = query.OP_NEQ         // 不等于
	OP_GT          FilterOperator = query.OP_GT          // 大于
	OP_GTE         FilterOperator = query.OP_GTE         // 大于等于
	OP_LT          FilterOperator = query.OP_LT          // 小于
	OP_LTE         FilterOperator = query.OP_LTE         // 小于等于
	OP_LIKE        FilterOperator = query.OP_LIKE        // 模糊匹配
	OP_IN          FilterOperator = query.OP_IN          // 包含
	OP_BETWEEN     FilterOperator = query.OP_BETWEEN     // 范围
	OP_IS_NULL     FilterOperator = query.OP_IS_NULL     // 为空
	OP_FIND_IN_SET FilterOperator = query.OP_FIND_IN_SET // MySQL FIND_IN_SET
)

// Deprecated: 使用 query.Filter 替代
type Filter = query.Filter

// Deprecated: 使用 query.FilterGroup 替代
type FilterGroup = query.FilterGroup

// Deprecated: 使用 query.BaseInfoFilter 替代
type BaseInfoFilter = query.BaseInfoFilter

// Deprecated: 使用 query.OrderBy 替代
type OrderBy = query.OrderBy

// Deprecated: 使用 query.Param 替代
// AdvancedQueryParam 被保留用于向后兼容，新代码应使用 query.Param
type AdvancedQueryParam = query.Param

// Deprecated: 使用 query.NewParam() 替代
// NewAdvancedQueryParam 创建高级查询参数 - 向后兼容适配器
func NewAdvancedQueryParam() *AdvancedQueryParam {
	return query.NewParam()
}

// ==================== 查询选项 ====================

// FindOption 查询选项 - 借鉴go-core设计
type FindOption struct {
	BusinessId    string        // 业务ID
	ShopId        string        // 店铺ID
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

// WithBusinessId 设置业务ID
func (fo *FindOption) WithBusinessId(id string) *FindOption {
	fo.BusinessId = id
	return fo
}

// WithShopId 设置店铺ID
func (fo *FindOption) WithShopId(id string) *FindOption {
	fo.ShopId = id
	return fo
}

// WithTablePrefix 设置表名前缀
func (fo *FindOption) WithTablePrefix(prefix string) *FindOption {
	fo.TablePrefix = prefix
	return fo
}

// WithNoCahce 禁用缓存
func (fo *FindOption) WithNoCache() *FindOption {
	fo.NoCache = true
	return fo
}

// WithCacheTTL 设置缓存TTL
func (fo *FindOption) WithCacheTTL(ttl time.Duration) *FindOption {
	fo.CacheTTL = ttl
	return fo
}

// ==================== 分页响应 ====================

// PageBean 分页响应 - 借鉴go-core设计
type PageBean struct {
	Page     int         `json:"page"`      // 当前页码
	PageSize int         `json:"page_size"` // 每页数量
	Total    int64       `json:"total"`     // 总数
	Pages    int         `json:"pages"`     // 总页数
	Rows     interface{} `json:"rows"`      // 数据行
}

// NewPageBean 创建分页响应
func NewPageBean(page int, pageSize int, total int64, rows interface{}) *PageBean {
	pages := int((total + int64(pageSize) - 1) / int64(pageSize))
	if page < 1 {
		page = 1
	}
	if pages < 1 {
		pages = 1
	}
	return &PageBean{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
		Pages:    pages,
		Rows:     rows,
	}
}

// ==================== 查询参数常量 ====================

// 比较操作符
const (
	QP_LT  = "lt"  // 小于
	QP_GT  = "gt"  // 大于
	QP_LTE = "lte" // 小于等于
	QP_GTE = "gte" // 大于等于
	QP_EQ  = "eq"  // 等于
	QP_NEQ = "neq" // 不等于
	QP_LK  = "lk"  // LIKE
)

// 排序操作符
const (
	QP_PD = "pd" // 降序（descending）
	QP_PA = "pa" // 升序（ascending）
)

// 或条件操作符
const (
	QP_ORLK = "orlk" // OR LIKE
	QP_ORLT = "orlt" // OR LT
	QP_ORGT = "orgt" // OR GT
)
