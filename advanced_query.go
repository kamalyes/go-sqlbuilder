/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-11 00:00:00
 * @FilePath: \go-sqlbuilder\advanced_query.go
 * @Description: 高级查询参数 - 借鉴go-core/pkg/database设计模式
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package sqlbuilder

import (
	"fmt"
	"strings"
	"time"

	logger "github.com/kamalyes/go-logger"
)

// ==================== 过滤器相关 ====================

// FilterOperator 过滤操作符
type FilterOperator string

const (
	OP_EQ          FilterOperator = "="           // 等于
	OP_NEQ         FilterOperator = "!="          // 不等于
	OP_GT          FilterOperator = ">"           // 大于
	OP_GTE         FilterOperator = ">="          // 大于等于
	OP_LT          FilterOperator = "<"           // 小于
	OP_LTE         FilterOperator = "<="          // 小于等于
	OP_LIKE        FilterOperator = "LIKE"        // 模糊匹配
	OP_IN          FilterOperator = "IN"          // 包含
	OP_BETWEEN     FilterOperator = "BETWEEN"     // 范围
	OP_IS_NULL     FilterOperator = "IS NULL"     // 为空
	OP_FIND_IN_SET FilterOperator = "FIND_IN_SET" // MySQL FIND_IN_SET
)

// Filter 单个过滤条件
type Filter struct {
	Field    string         // 字段名
	Operator FilterOperator // 操作符
	Value    interface{}    // 值（可能是单个值或数组）
	Logic    string         // AND / OR
}

// FilterGroup 过滤组（支持嵌套）
type FilterGroup struct {
	Filters []*Filter
	Groups  []*FilterGroup
	Logic   string // AND / OR
}

// BaseInfoFilter 基础过滤信息 - 借鉴go-core设计
type BaseInfoFilter struct {
	DBField    string        // 数据库字段名
	Values     []interface{} // 过滤值列表
	ExactMatch bool          // 是否精确匹配（默认false表示LIKE）
	AllRegex   bool          // 是否全部匹配正则（所有值都包含）
}

// ==================== 高级查询参数 ====================

// AdvancedQueryParam 高级查询参数 - 支持复杂过滤、排序、分页
type AdvancedQueryParam struct {
	Filters       []*Filter                 // 过滤条件
	FilterGroups  []*FilterGroup            // 过滤组（支持AND/OR混合）
	TimeRanges    map[string][2]interface{} // 时间范围 field -> [startTime, endTime]
	FindInSets    map[string][]string       // FIND_IN_SET条件
	Orders        []OrderBy                 // 排序条件
	Page          int                       // 页码（从1开始）
	PageSize      int                       // 每页数量
	Offset        int                       // 偏移量（与Page互斥）
	Limit         int                       // 限制数量
	Distinct      bool                      // 是否去重
	SelectFields  []string                  // 指定查询字段
	HavingClauses []string                  // HAVING子句
}

// OrderBy 排序条件
type OrderBy struct {
	Field string // 字段名
	Order string // ASC / DESC
}

// NewAdvancedQueryParam 创建高级查询参数
func NewAdvancedQueryParam() *AdvancedQueryParam {
	return &AdvancedQueryParam{
		Filters:       make([]*Filter, 0),
		FilterGroups:  make([]*FilterGroup, 0),
		TimeRanges:    make(map[string][2]interface{}),
		FindInSets:    make(map[string][]string),
		Orders:        make([]OrderBy, 0),
		Page:          1,
		PageSize:      10,
		SelectFields:  make([]string, 0),
		HavingClauses: make([]string, 0),
	}
}

// AddFilter 添加单个过滤条件
func (aq *AdvancedQueryParam) AddFilter(field string, operator FilterOperator, value interface{}) *AdvancedQueryParam {
	aq.Filters = append(aq.Filters, &Filter{
		Field:    field,
		Operator: operator,
		Value:    value,
		Logic:    "AND",
	})
	return aq
}

// AddOrFilter 添加OR过滤条件
func (aq *AdvancedQueryParam) AddOrFilter(field string, operator FilterOperator, value interface{}) *AdvancedQueryParam {
	if len(aq.Filters) > 0 {
		aq.Filters[len(aq.Filters)-1].Logic = "OR"
	}
	return aq.AddFilter(field, operator, value)
}

// AddEQFilter 添加等于过滤
func (aq *AdvancedQueryParam) AddEQFilter(field string, value interface{}) *AdvancedQueryParam {
	return aq.AddFilter(field, OP_EQ, value)
}

// AddEQ 是 AddEQFilter 的简洁包装，便于链式调用
func (aq *AdvancedQueryParam) AddEQ(field string, value interface{}) *AdvancedQueryParam {
	return aq.AddEQFilter(field, value)
}

// AddOrEQ 添加 OR 等于 过滤条件
func (aq *AdvancedQueryParam) AddOrEQ(field string, value interface{}) *AdvancedQueryParam {
	return aq.AddOrFilter(field, OP_EQ, value)
}

// AddOrLike 添加 OR LIKE 过滤条件
func (aq *AdvancedQueryParam) AddOrLike(field string, value string) *AdvancedQueryParam {
	return aq.AddOrFilter(field, OP_LIKE, "%"+value+"%")
}

// AddOrIn 添加 OR IN 过滤条件
func (aq *AdvancedQueryParam) AddOrIn(field string, values ...interface{}) *AdvancedQueryParam {
	return aq.AddOrFilter(field, OP_IN, values)
}

// AddLikeFilter 添加LIKE过滤
func (aq *AdvancedQueryParam) AddLikeFilter(field string, value string) *AdvancedQueryParam {
	return aq.AddFilter(field, OP_LIKE, "%"+value+"%")
}

// AddLike 是 AddLikeFilter 的简洁包装（全模糊）
func (aq *AdvancedQueryParam) AddLike(field string, value string) *AdvancedQueryParam {
	return aq.AddLikeFilter(field, value)
}

// AddLikeStartFilter 添加LIKE过滤（前缀匹配）
func (aq *AdvancedQueryParam) AddLikeStartFilter(field string, value string) *AdvancedQueryParam {
	return aq.AddFilter(field, OP_LIKE, value+"%")
}

// AddStartsWith 前缀匹配（简洁包装）
func (aq *AdvancedQueryParam) AddStartsWith(field string, value string) *AdvancedQueryParam {
	return aq.AddLikeStartFilter(field, value)
}

// AddEndsWith 后缀匹配（LIKE %value）
func (aq *AdvancedQueryParam) AddEndsWith(field string, value string) *AdvancedQueryParam {
	return aq.AddFilter(field, OP_LIKE, "%"+value)
}

// AddInFilter 添加IN过滤
func (aq *AdvancedQueryParam) AddInFilter(field string, values ...interface{}) *AdvancedQueryParam {
	return aq.AddFilter(field, OP_IN, values)
}

// AddIn 是 AddInFilter 的简洁包装
func (aq *AdvancedQueryParam) AddIn(field string, values ...interface{}) *AdvancedQueryParam {
	return aq.AddInFilter(field, values...)
}

// AddGT 添加大于过滤
func (aq *AdvancedQueryParam) AddGT(field string, value interface{}) *AdvancedQueryParam {
	return aq.AddFilter(field, OP_GT, value)
}

// AddGTE 添加大于等于过滤
func (aq *AdvancedQueryParam) AddGTE(field string, value interface{}) *AdvancedQueryParam {
	return aq.AddFilter(field, OP_GTE, value)
}

// AddLT 添加小于过滤
func (aq *AdvancedQueryParam) AddLT(field string, value interface{}) *AdvancedQueryParam {
	return aq.AddFilter(field, OP_LT, value)
}

// AddLTE 添加小于等于过滤
func (aq *AdvancedQueryParam) AddLTE(field string, value interface{}) *AdvancedQueryParam {
	return aq.AddFilter(field, OP_LTE, value)
}

// AddNEQ 添加不等于过滤
func (aq *AdvancedQueryParam) AddNEQ(field string, value interface{}) *AdvancedQueryParam {
	return aq.AddFilter(field, OP_NEQ, value)
}

// AddOrGT 添加 OR 大于 过滤
func (aq *AdvancedQueryParam) AddOrGT(field string, value interface{}) *AdvancedQueryParam {
	return aq.AddOrFilter(field, OP_GT, value)
}

// AddOrGTE 添加 OR 大于等于 过滤
func (aq *AdvancedQueryParam) AddOrGTE(field string, value interface{}) *AdvancedQueryParam {
	return aq.AddOrFilter(field, OP_GTE, value)
}

// AddOrLT 添加 OR 小于 过滤
func (aq *AdvancedQueryParam) AddOrLT(field string, value interface{}) *AdvancedQueryParam {
	return aq.AddOrFilter(field, OP_LT, value)
}

// AddOrLTE 添加 OR 小于等于 过滤
func (aq *AdvancedQueryParam) AddOrLTE(field string, value interface{}) *AdvancedQueryParam {
	return aq.AddOrFilter(field, OP_LTE, value)
}

// AddOrNEQ 添加 OR 不等于 过滤
func (aq *AdvancedQueryParam) AddOrNEQ(field string, value interface{}) *AdvancedQueryParam {
	return aq.AddOrFilter(field, OP_NEQ, value)
}

// AddTimeRange 添加时间范围过滤
func (aq *AdvancedQueryParam) AddTimeRange(field string, startTime, endTime interface{}) *AdvancedQueryParam {
	aq.TimeRanges[field] = [2]interface{}{startTime, endTime}
	return aq
}

// AddFindInSet 添加FIND_IN_SET过滤（MySQL特定）
func (aq *AdvancedQueryParam) AddFindInSet(field string, values ...string) *AdvancedQueryParam {
	aq.FindInSets[field] = values
	return aq
}

// AddOrder 添加排序
func (aq *AdvancedQueryParam) AddOrder(field string, order string) *AdvancedQueryParam {
	if order == "" {
		order = "ASC"
	}
	aq.Orders = append(aq.Orders, OrderBy{Field: field, Order: strings.ToUpper(order)})
	return aq
}

// AddOrderDesc 添加降序排序
func (aq *AdvancedQueryParam) AddOrderDesc(field string) *AdvancedQueryParam {
	return aq.AddOrder(field, "DESC")
}

// AddOrderAsc 添加升序排序
func (aq *AdvancedQueryParam) AddOrderAsc(field string) *AdvancedQueryParam {
	return aq.AddOrder(field, "ASC")
}

// SetPage 设置分页
func (aq *AdvancedQueryParam) SetPage(page int, pageSize int) *AdvancedQueryParam {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	aq.Page = page
	aq.PageSize = pageSize
	aq.Offset = (page - 1) * pageSize
	return aq
}

// SetDistinct 设置去重
func (aq *AdvancedQueryParam) SetDistinct(distinct bool) *AdvancedQueryParam {
	aq.Distinct = distinct
	return aq
}

// SetSelectFields 设置查询字段
func (aq *AdvancedQueryParam) SetSelectFields(fields ...string) *AdvancedQueryParam {
	aq.SelectFields = fields
	return aq
}

// AddHaving 添加HAVING子句
func (aq *AdvancedQueryParam) AddHaving(clause string) *AdvancedQueryParam {
	aq.HavingClauses = append(aq.HavingClauses, clause)
	return aq
}

// BuildWhereClause 构建WHERE子句 - 返回WHERE SQL片段和参数
func (aq *AdvancedQueryParam) BuildWhereClause() (string, []interface{}) {
	var whereClauses []string
	var args []interface{}

	// 处理普通过滤
	for _, filter := range aq.Filters {
		sql, filterArgs := aq.buildFilterSQL(filter)
		whereClauses = append(whereClauses, sql)
		args = append(args, filterArgs...)
	}

	// 处理时间范围
	for field, timeRange := range aq.TimeRanges {
		whereClauses = append(whereClauses, fmt.Sprintf("%s BETWEEN ? AND ?", field))
		args = append(args, timeRange[0], timeRange[1])
	}

	// 处理FIND_IN_SET
	for field, values := range aq.FindInSets {
		for _, value := range values {
			whereClauses = append(whereClauses, fmt.Sprintf("FIND_IN_SET(?, %s) > 0", field))
			args = append(args, value)
		}
	}

	if len(whereClauses) == 0 {
		// 使用全局日志记录没有过滤条件的情况（便于调试）
		logger.Debug("AdvancedQueryParam: no where clauses built")
		return "", args
	}

	return "WHERE " + strings.Join(whereClauses, " AND "), args
}

// buildFilterSQL 构建单个过滤SQL
func (aq *AdvancedQueryParam) buildFilterSQL(filter *Filter) (string, []interface{}) {
	var args []interface{}
	var sql string

	switch filter.Operator {
	case OP_IN:
		values := filter.Value.([]interface{})
		placeholders := make([]string, len(values))
		for i, v := range values {
			placeholders[i] = "?"
			args = append(args, v)
		}
		sql = fmt.Sprintf("%s IN (%s)", filter.Field, strings.Join(placeholders, ","))

	case OP_LIKE:
		sql = fmt.Sprintf("%s LIKE ?", filter.Field)
		args = append(args, filter.Value)

	case OP_BETWEEN:
		values := filter.Value.([2]interface{})
		sql = fmt.Sprintf("%s BETWEEN ? AND ?", filter.Field)
		args = append(args, values[0], values[1])

	case OP_IS_NULL:
		sql = fmt.Sprintf("%s IS NULL", filter.Field)

	default:
		sql = fmt.Sprintf("%s %s ?", filter.Field, filter.Operator)
		args = append(args, filter.Value)
	}

	return sql, args
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
