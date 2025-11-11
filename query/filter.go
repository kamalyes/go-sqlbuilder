/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-11 00:00:00
 * @FilePath: \go-sqlbuilder\query\filter.go
 * @Description: 查询过滤器结构定义
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package query

// Filter 单个过滤条件
type Filter struct {
	Field    string      // 字段名
	Operator Operator    // 操作符
	Value    interface{} // 值（可能是单个值或数组）
	Logic    string      // AND / OR
}

// FilterGroup 过滤组（支持嵌套）
type FilterGroup struct {
	Filters []*Filter
	Groups  []*FilterGroup
	Logic   string // AND / OR
}

// BaseInfoFilter 基础过滤信息 - 借鉴 go-core 设计
type BaseInfoFilter struct {
	DBField    string        // 数据库字段名
	Values     []interface{} // 过滤值列表
	ExactMatch bool          // 是否精确匹配（默认 false 表示 LIKE）
	AllRegex   bool          // 是否全部匹配正则（所有值都包含）
}

// NewFilter 创建过滤条件
func NewFilter(field string, operator Operator, value interface{}) *Filter {
	return &Filter{
		Field:    field,
		Operator: operator,
		Value:    value,
		Logic:    "AND",
	}
}

// NewBaseInfoFilter 创建基础过滤信息
func NewBaseInfoFilter(field string, values []interface{}) *BaseInfoFilter {
	return &BaseInfoFilter{
		DBField:    field,
		Values:     values,
		ExactMatch: false,
		AllRegex:   false,
	}
}
