/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-23 08:15:15
 * @FilePath: \go-sqlbuilder\repository\filter.go
 * @Description: 过滤器和查询构建 - Filter、FilterGroup、Query等查询条件
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package repository

import "github.com/kamalyes/go-sqlbuilder/constants"

// Filter 过滤条件
type Filter struct {
	Field    string             // 字段名
	Operator constants.Operator // 操作符
	Value    interface{}        // 值
}

// FilterGroup 过滤条件组，支持逻辑操作
type FilterGroup struct {
	Filters []*Filter          // 过滤条件列表
	Groups  []*FilterGroup     // 嵌套条件组
	LogicOp constants.Operator // 逻辑操作符：AND/OR
}

// NewFilterGroup 创建新的过滤条件组
func NewFilterGroup(logicOp constants.Operator) *FilterGroup {
	if logicOp != constants.LOGIC_AND && logicOp != constants.LOGIC_OR {
		logicOp = constants.LOGIC_AND // 默认为 AND
	}
	return &FilterGroup{
		Filters: make([]*Filter, 0),
		Groups:  make([]*FilterGroup, 0),
		LogicOp: logicOp,
	}
}

// AddFilter 向条件组添加过滤条件
func (fg *FilterGroup) AddFilter(filter *Filter) *FilterGroup {
	if filter != nil {
		fg.Filters = append(fg.Filters, filter)
	}
	return fg
}

// AddFilters 向条件组批量添加过滤条件
func (fg *FilterGroup) AddFilters(filters ...*Filter) *FilterGroup {
	for _, f := range filters {
		fg.AddFilter(f)
	}
	return fg
}

// AddGroup 向条件组添加嵌套条件组
func (fg *FilterGroup) AddGroup(group *FilterGroup) *FilterGroup {
	if group != nil {
		fg.Groups = append(fg.Groups, group)
	}
	return fg
}

// IsEmpty 检查条件组是否为空
func (fg *FilterGroup) IsEmpty() bool {
	return len(fg.Filters) == 0 && len(fg.Groups) == 0
}

// Count 返回条件总数（包括嵌套组）
func (fg *FilterGroup) Count() int {
	count := len(fg.Filters)
	for _, group := range fg.Groups {
		count += group.Count()
	}
	return count
}

// Order 排序条件
type Order struct {
	Field     string
	Direction string // "ASC" or "DESC"
}

// Query 查询条件
type Query struct {
	Filters     []*Filter    // 简单过滤条件
	FilterGroup *FilterGroup // 复合过滤条件组
	Orders      []Order      // 排序条件
	Pagination  *Pagination  // 分页信息
	LimitValue  *int         // 限制数量
	OffsetValue *int         // 偏移量
	Distinct    bool         // 是否去重
	GroupBy     []string     // 分组字段
	Having      []*Filter    // HAVING 条件
}

// NewEqFilter 创建等于过滤条件
func NewEqFilter(field string, value interface{}) *Filter {
	return &Filter{Field: field, Operator: constants.OP_EQ, Value: value}
}

// NewGtFilter 创建大于过滤条件
func NewGtFilter(field string, value interface{}) *Filter {
	return &Filter{Field: field, Operator: constants.OP_GT, Value: value}
}

// NewLtFilter 创建小于过滤条件
func NewLtFilter(field string, value interface{}) *Filter {
	return &Filter{Field: field, Operator: constants.OP_LT, Value: value}
}

// NewGteFilter 创建大于等于过滤条件
func NewGteFilter(field string, value interface{}) *Filter {
	return &Filter{Field: field, Operator: constants.OP_GTE, Value: value}
}

// NewLteFilter 创建小于等于过滤条件
func NewLteFilter(field string, value interface{}) *Filter {
	return &Filter{Field: field, Operator: constants.OP_LTE, Value: value}
}

// NewInFilter 创建 IN 过滤条件
func NewInFilter(field string, values ...interface{}) *Filter {
	if values == nil {
		values = make([]interface{}, 0)
	}
	return &Filter{Field: field, Operator: constants.OP_IN, Value: values}
}

// NewLikeFilter 创建 LIKE 过滤条件
func NewLikeFilter(field string, value string) *Filter {
	return &Filter{Field: field, Operator: constants.OP_LIKE, Value: "%" + value + "%"}
}

// NewNeqFilter 创建不等于过滤条件
func NewNeqFilter(field string, value interface{}) *Filter {
	return &Filter{Field: field, Operator: constants.OP_NEQ, Value: value}
}

// NewBetweenFilter 创建 BETWEEN 过滤条件
func NewBetweenFilter(field string, min, max interface{}) *Filter {
	return &Filter{Field: field, Operator: constants.OP_BETWEEN, Value: []interface{}{min, max}}
}

// NewNotInFilter 创建 NOT IN 过滤条件
func NewNotInFilter(field string, values ...interface{}) *Filter {
	if values == nil {
		values = make([]interface{}, 0)
	}
	return &Filter{Field: field, Operator: constants.OP_NOT_IN, Value: values}
}

// NewIsNullFilter 创建 IS NULL 过滤条件
func NewIsNullFilter(field string) *Filter {
	return &Filter{Field: field, Operator: constants.OP_IS_NULL, Value: nil}
}

// NewIsNotNullFilter 创建 IS NOT NULL 过滤条件
func NewIsNotNullFilter(field string) *Filter {
	return &Filter{Field: field, Operator: constants.OP_IS_NOT_NULL, Value: nil}
}

// NewStartsWithFilter 创建前缀匹配过滤条件
func NewStartsWithFilter(field string, value string) *Filter {
	return &Filter{Field: field, Operator: constants.OP_LIKE, Value: value + "%"}
}

// NewEndsWithFilter 创建后缀匹配过滤条件
func NewEndsWithFilter(field string, value string) *Filter {
	return &Filter{Field: field, Operator: constants.OP_LIKE, Value: "%" + value}
}

// NewContainsFilter 创建包含匹配过滤条件（与 NewLikeFilter 相同）
func NewContainsFilter(field string, value string) *Filter {
	return NewLikeFilter(field, value)
}

// NewNotLikeFilter 创建 NOT LIKE 过滤条件
func NewNotLikeFilter(field string, value string) *Filter {
	return &Filter{Field: field, Operator: constants.OP_NOT_LIKE, Value: "%" + value + "%"}
}

// NewFindInSetFilter 创建 FIND_IN_SET 过滤条件（MySQL特定）
func NewFindInSetFilter(field string, value interface{}) *Filter {
	return &Filter{Field: field, Operator: constants.OP_FIND_IN_SET, Value: value}
}
