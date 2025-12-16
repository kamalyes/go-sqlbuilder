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

import (
	"reflect"

	"github.com/kamalyes/go-sqlbuilder/constants"
	"github.com/kamalyes/go-toolbox/pkg/validator"
)

// SubQuery 子查询结构
type SubQuery struct {
	SQL  string        // 子查询 SQL
	Args []interface{} // 子查询参数
}

// NewSubQuery 创建子查询
func NewSubQuery(sql string, args ...interface{}) *SubQuery {
	return &SubQuery{
		SQL:  sql,
		Args: args,
	}
}

// Filter 过滤条件
type Filter struct {
	Field    string             // 字段名
	Operator constants.Operator // 操作符
	Value    interface{}        // 值（可以是普通值或 *SubQuery）
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

// AddFilterIf 当条件为真时添加过滤条件
func (fg *FilterGroup) AddFilterIf(condition bool, filter *Filter) *FilterGroup {
	if condition && filter != nil {
		fg.Filters = append(fg.Filters, filter)
	}
	return fg
}

// AddFilterIfNotEmpty 当值不为空时添加过滤条件
// 支持 string, slice, map 等类型的空值判断
func (fg *FilterGroup) AddFilterIfNotEmpty(field string, operator constants.Operator, value interface{}) *FilterGroup {
	if !validator.IsEmptyValue(reflect.ValueOf(value)) {
		fg.AddFilter(&Filter{Field: field, Operator: operator, Value: value})
	}
	return fg
}

// AddEqFilterIfNotEmpty 当值不为空时添加等于过滤条件
func (fg *FilterGroup) AddEqFilterIfNotEmpty(field string, value interface{}) *FilterGroup {
	return fg.AddFilterIfNotEmpty(field, constants.OP_EQ, value)
}

// AddNeqFilterIfNotEmpty 当值不为空时添加不等于过滤条件
func (fg *FilterGroup) AddNeqFilterIfNotEmpty(field string, value interface{}) *FilterGroup {
	return fg.AddFilterIfNotEmpty(field, constants.OP_NEQ, value)
}

// AddGtFilterIfNotEmpty 当值不为空时添加大于过滤条件
func (fg *FilterGroup) AddGtFilterIfNotEmpty(field string, value interface{}) *FilterGroup {
	return fg.AddFilterIfNotEmpty(field, constants.OP_GT, value)
}

// AddGteFilterIfNotEmpty 当值不为空时添加大于等于过滤条件
func (fg *FilterGroup) AddGteFilterIfNotEmpty(field string, value interface{}) *FilterGroup {
	return fg.AddFilterIfNotEmpty(field, constants.OP_GTE, value)
}

// AddLtFilterIfNotEmpty 当值不为空时添加小于过滤条件
func (fg *FilterGroup) AddLtFilterIfNotEmpty(field string, value interface{}) *FilterGroup {
	return fg.AddFilterIfNotEmpty(field, constants.OP_LT, value)
}

// AddLteFilterIfNotEmpty 当值不为空时添加小于等于过滤条件
func (fg *FilterGroup) AddLteFilterIfNotEmpty(field string, value interface{}) *FilterGroup {
	return fg.AddFilterIfNotEmpty(field, constants.OP_LTE, value)
}

// AddLikeFilterIfNotEmpty 当值不为空时添加 LIKE 过滤条件
func (fg *FilterGroup) AddLikeFilterIfNotEmpty(field string, value string) *FilterGroup {
	if !validator.IsEmptyValue(reflect.ValueOf(value)) {
		fg.AddFilter(&Filter{Field: field, Operator: constants.OP_LIKE, Value: "%" + value + "%"})
	}
	return fg
}

// AddInFilterIfNotEmpty 当切片不为空时添加 IN 过滤条件
func (fg *FilterGroup) AddInFilterIfNotEmpty(field string, values []interface{}) *FilterGroup {
	if !validator.IsEmptyValue(reflect.ValueOf(values)) {
		fg.AddFilter(&Filter{Field: field, Operator: constants.OP_IN, Value: values})
	}
	return fg
}

// AddNotInFilterIfNotEmpty 当切片不为空时添加 NOT IN 过滤条件
func (fg *FilterGroup) AddNotInFilterIfNotEmpty(field string, values []interface{}) *FilterGroup {
	if !validator.IsEmptyValue(reflect.ValueOf(values)) {
		fg.AddFilter(&Filter{Field: field, Operator: constants.OP_NOT_IN, Value: values})
	}
	return fg
}

// AddBetweenFilterIfNotEmpty 当最小值和最大值都不为空时添加 BETWEEN 过滤条件
func (fg *FilterGroup) AddBetweenFilterIfNotEmpty(field string, min, max interface{}) *FilterGroup {
	if !validator.IsEmptyValue(reflect.ValueOf(min)) && !validator.IsEmptyValue(reflect.ValueOf(max)) {
		fg.AddFilter(&Filter{Field: field, Operator: constants.OP_BETWEEN, Value: []interface{}{min, max}})
	}
	return fg
}

// AddStartsWithFilterIfNotEmpty 当值不为空时添加前缀匹配过滤条件
func (fg *FilterGroup) AddStartsWithFilterIfNotEmpty(field string, value string) *FilterGroup {
	if !validator.IsEmptyValue(reflect.ValueOf(value)) {
		fg.AddFilter(&Filter{Field: field, Operator: constants.OP_LIKE, Value: value + "%"})
	}
	return fg
}

// AddEndsWithFilterIfNotEmpty 当值不为空时添加后缀匹配过滤条件
func (fg *FilterGroup) AddEndsWithFilterIfNotEmpty(field string, value string) *FilterGroup {
	if !validator.IsEmptyValue(reflect.ValueOf(value)) {
		fg.AddFilter(&Filter{Field: field, Operator: constants.OP_LIKE, Value: "%" + value})
	}
	return fg
}

// AddContainsFilterIfNotEmpty 当值不为空时添加包含匹配过滤条件
func (fg *FilterGroup) AddContainsFilterIfNotEmpty(field string, value string) *FilterGroup {
	return fg.AddLikeFilterIfNotEmpty(field, value)
}

// AddNotLikeFilterIfNotEmpty 当值不为空时添加 NOT LIKE 过滤条件
func (fg *FilterGroup) AddNotLikeFilterIfNotEmpty(field string, value string) *FilterGroup {
	if !validator.IsEmptyValue(reflect.ValueOf(value)) {
		fg.AddFilter(&Filter{Field: field, Operator: constants.OP_NOT_LIKE, Value: "%" + value + "%"})
	}
	return fg
}

// AddFindInSetFilterIfNotEmpty 当值不为空时添加 FIND_IN_SET 过滤条件（MySQL特定）
func (fg *FilterGroup) AddFindInSetFilterIfNotEmpty(field string, value interface{}) *FilterGroup {
	if !validator.IsEmptyValue(reflect.ValueOf(value)) {
		fg.AddFilter(&Filter{Field: field, Operator: constants.OP_FIND_IN_SET, Value: value})
	}
	return fg
}

// AddGroupIf 当条件为真时添加嵌套条件组
func (fg *FilterGroup) AddGroupIf(condition bool, group *FilterGroup) *FilterGroup {
	if condition && group != nil && !group.IsEmpty() {
		fg.Groups = append(fg.Groups, group)
	}
	return fg
}

// AddGroupIfNotEmpty 当嵌套条件组不为空时添加
func (fg *FilterGroup) AddGroupIfNotEmpty(group *FilterGroup) *FilterGroup {
	if group != nil && !group.IsEmpty() {
		fg.Groups = append(fg.Groups, group)
	}
	return fg
}

// Clear 清空所有过滤条件和条件组
func (fg *FilterGroup) Clear() *FilterGroup {
	fg.Filters = make([]*Filter, 0)
	fg.Groups = make([]*FilterGroup, 0)
	return fg
}

// Clone 克隆条件组（深拷贝）
func (fg *FilterGroup) Clone() *FilterGroup {
	newGroup := NewFilterGroup(fg.LogicOp)

	// 克隆过滤条件
	for _, f := range fg.Filters {
		newFilter := &Filter{
			Field:    f.Field,
			Operator: f.Operator,
			Value:    f.Value,
		}
		newGroup.Filters = append(newGroup.Filters, newFilter)
	}

	// 克隆嵌套条件组
	for _, g := range fg.Groups {
		newGroup.Groups = append(newGroup.Groups, g.Clone())
	}

	return newGroup
}

// Order 排序条件
type Order struct {
	Field     string
	Direction string // "ASC" or "DESC"
}

// Query 查询条件
type Query struct {
	Filters      []*Filter    // 简单过滤条件
	FilterGroup  *FilterGroup // 复合过滤条件组
	Orders       []Order      // 排序条件
	Pagination   *Pagination  // 分页信息
	LimitValue   *int         // 限制数量
	OffsetValue  *int         // 偏移量
	Distinct     bool         // 是否去重
	GroupBy      []string     // 分组字段
	Having       []*Filter    // HAVING 条件
	SelectFields []string     // 要查询的字段列表（为空则查询所有字段）
	OmitFields   []string     // 要排除的字段列表
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

// NewInFilterSlice 创建 IN 过滤条件（使用切片参数）
func NewInFilterSlice(field string, values []interface{}) *Filter {
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

// NewNotInFilterSlice 创建 NOT IN 过滤条件（使用切片参数）
func NewNotInFilterSlice(field string, values []interface{}) *Filter {
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
