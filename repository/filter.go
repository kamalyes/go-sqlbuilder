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
	"fmt"

	"github.com/kamalyes/go-sqlbuilder/constants"
	"github.com/kamalyes/go-toolbox/pkg/convert"
	"github.com/kamalyes/go-toolbox/pkg/validator"
)

// MaxFilterGroupDepth 过滤条件组最大嵌套深度，防止无限嵌套导致内存溢出
const MaxFilterGroupDepth = 20

// SubQuery 子查询结构
type SubQuery struct {
	SQL  string        // 子查询 SQL
	Args []interface{} // 子查询参数
}

// NewSubQuery 创建子查询
func NewSubQuery(sql string, args ...interface{}) *SubQuery {
	return &SubQuery{
		SQL:  sql,
		Args: validator.NormalizeFilterValueSlice(args),
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

// AddGroup 向条件组添加嵌套条件组（检查嵌套深度限制）
func (fg *FilterGroup) AddGroup(group *FilterGroup) *FilterGroup {
	if group != nil && fg.getDepth() < MaxFilterGroupDepth {
		fg.Groups = append(fg.Groups, group)
	}
	return fg
}

// getDepth 计算条件组的嵌套深度
func (fg *FilterGroup) getDepth() int {
	if len(fg.Groups) == 0 {
		return 1
	}
	maxDepth := 1
	for _, g := range fg.Groups {
		if g != nil {
			if d := g.getDepth(); d+1 > maxDepth {
				maxDepth = d + 1
			}
		}
	}
	return maxDepth
}

// IsEmpty 检查条件组是否为空
func (fg *FilterGroup) IsEmpty() bool {
	return len(fg.Filters) == 0 && len(fg.Groups) == 0
}

// Count 返回条件总数（包括嵌套组）
func (fg *FilterGroup) Count() int {
	if fg == nil {
		return 0
	}
	count := len(fg.Filters)
	for _, group := range fg.Groups {
		if group != nil {
			count += group.Count()
		}
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

// AddFilterIfValueNotEmpty 当 filter 的 Value 不为空时添加过滤条件
func (fg *FilterGroup) AddFilterIfValueNotEmpty(filter *Filter) *FilterGroup {
	if filter == nil {
		return fg
	}
	deref, empty := validator.NormalizeFilterValueIfNotEmpty(filter.Value)
	if empty {
		return fg
	}
	filter.Value = deref
	fg.AddFilter(filter)
	return fg
}

// AddFilterIfNotEmpty 当值不为空时添加过滤条件
// 支持指针自动解引用，string, slice, map 等类型的空值判断
func (fg *FilterGroup) AddFilterIfNotEmpty(field string, operator constants.Operator, value interface{}) *FilterGroup {
	deref, empty := validator.NormalizeFilterValueIfNotEmpty(value)
	if empty {
		return fg
	}
	fg.AddFilter(NewFilter(field, operator, deref))
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
// value 支持 string 和 *string 类型
func (fg *FilterGroup) AddLikeFilterIfNotEmpty(field string, value interface{}) *FilterGroup {
	deref, empty := validator.NormalizeFilterValueIfNotEmpty(value)
	if empty {
		return fg
	}
	fg.AddFilter(NewLikeFilter(field, fmt.Sprintf("%v", deref)))
	return fg
}

// AddInFilterIfNotEmpty 当切片不为空时添加 IN 过滤条件
// values 支持任意切片类型（[]string、[]int 等）以及对应指针类型
func (fg *FilterGroup) AddInFilterIfNotEmpty(field string, values interface{}) *FilterGroup {
	deref, empty := validator.NormalizeFilterValueIfNotEmpty(values)
	if empty {
		return fg
	}
	if slice := convert.AnySliceToInterfaceSlice(deref); len(slice) > 0 {
		fg.AddFilter(NewInFilterSlice(field, slice))
	}
	return fg
}

// AddNotInFilterIfNotEmpty 当切片不为空时添加 NOT IN 过滤条件
// values 支持任意切片类型（[]string、[]int 等）以及对应指针类型
func (fg *FilterGroup) AddNotInFilterIfNotEmpty(field string, values interface{}) *FilterGroup {
	deref, empty := validator.NormalizeFilterValueIfNotEmpty(values)
	if empty {
		return fg
	}
	if slice := convert.AnySliceToInterfaceSlice(deref); len(slice) > 0 {
		fg.AddFilter(NewNotInFilterSlice(field, slice))
	}
	return fg
}

// AddBetweenFilterIfNotEmpty 当最小值和最大值都不为空时添加 BETWEEN 过滤条件
func (fg *FilterGroup) AddBetweenFilterIfNotEmpty(field string, min, max interface{}) *FilterGroup {
	minDeref, minEmpty := validator.NormalizeFilterValueIfNotEmpty(min)
	maxDeref, maxEmpty := validator.NormalizeFilterValueIfNotEmpty(max)
	if minEmpty || maxEmpty {
		return fg
	}
	fg.AddFilter(NewBetweenFilter(field, minDeref, maxDeref))
	return fg
}

// AddStartsWithFilterIfNotEmpty 当值不为空时添加前缀匹配过滤条件
// value 支持 string 和 *string 类型
func (fg *FilterGroup) AddStartsWithFilterIfNotEmpty(field string, value interface{}) *FilterGroup {
	deref, empty := validator.NormalizeFilterValueIfNotEmpty(value)
	if empty {
		return fg
	}
	fg.AddFilter(NewStartsWithFilter(field, fmt.Sprintf("%v", deref)))
	return fg
}

// AddEndsWithFilterIfNotEmpty 当值不为空时添加后缀匹配过滤条件
// value 支持 string 和 *string 类型
func (fg *FilterGroup) AddEndsWithFilterIfNotEmpty(field string, value interface{}) *FilterGroup {
	deref, empty := validator.NormalizeFilterValueIfNotEmpty(value)
	if empty {
		return fg
	}
	fg.AddFilter(NewEndsWithFilter(field, fmt.Sprintf("%v", deref)))
	return fg
}

// AddRegexpFilterIfNotEmpty 当值不为空时添加正则匹配过滤条件（仅MySQL/PostgreSQL）
// pattern 支持 string 和 *string 类型
func (fg *FilterGroup) AddRegexpFilterIfNotEmpty(field string, pattern interface{}) *FilterGroup {
	deref, empty := validator.NormalizeFilterValueIfNotEmpty(pattern)
	if empty {
		return fg
	}
	fg.AddFilter(NewRegexpFilter(field, fmt.Sprintf("%v", deref)))
	return fg
}

// AddNotLikeFilterIfNotEmpty 当值不为空时添加 NOT LIKE 过滤条件
// value 支持 string 和 *string 类型
func (fg *FilterGroup) AddNotLikeFilterIfNotEmpty(field string, value interface{}) *FilterGroup {
	deref, empty := validator.NormalizeFilterValueIfNotEmpty(value)
	if empty {
		return fg
	}
	fg.AddFilter(NewNotLikeFilter(field, fmt.Sprintf("%v", deref)))
	return fg
}

// AddFindInSetFilterIfNotEmpty 当值不为空时添加 FIND_IN_SET 过滤条件（MySQL特定）
func (fg *FilterGroup) AddFindInSetFilterIfNotEmpty(field string, value interface{}) *FilterGroup {
	deref, empty := validator.NormalizeFilterValueIfNotEmpty(value)
	if empty {
		return fg
	}
	fg.AddFilter(&Filter{Field: field, Operator: constants.OP_FIND_IN_SET, Value: deref})
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

// Clone 克隆条件组（深拷贝，检查嵌套深度）
func (fg *FilterGroup) Clone() *FilterGroup {
	return fg.cloneWithDepth(0)
}

// cloneWithDepth 递归克隆条件组，检查嵌套深度
func (fg *FilterGroup) cloneWithDepth(depth int) *FilterGroup {
	if depth >= MaxFilterGroupDepth {
		return NewFilterGroup(fg.LogicOp)
	}

	newGroup := NewFilterGroup(fg.LogicOp)

	// 克隆过滤条件
	for _, f := range fg.Filters {
		newFilter := NewFilter(f.Field, f.Operator, f.Value)
		newGroup.Filters = append(newGroup.Filters, newFilter)
	}

	// 克隆嵌套条件组
	for _, g := range fg.Groups {
		if g != nil {
			newGroup.Groups = append(newGroup.Groups, g.cloneWithDepth(depth+1))
		}
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
	return NewFilter(field, constants.OP_EQ, value)
}

// NewGtFilter 创建大于过滤条件
func NewGtFilter(field string, value interface{}) *Filter {
	return NewFilter(field, constants.OP_GT, value)
}

// NewLtFilter 创建小于过滤条件
func NewLtFilter(field string, value interface{}) *Filter {
	return NewFilter(field, constants.OP_LT, value)
}

// NewGteFilter 创建大于等于过滤条件
func NewGteFilter(field string, value interface{}) *Filter {
	return NewFilter(field, constants.OP_GTE, value)
}

// NewLteFilter 创建小于等于过滤条件
func NewLteFilter(field string, value interface{}) *Filter {
	return NewFilter(field, constants.OP_LTE, value)
}

// NewInFilter 创建 IN 过滤条件
func NewInFilter(field string, values ...interface{}) *Filter {
	if values == nil || (len(values) == 1 && values[0] == nil) {
		return &Filter{Field: field, Operator: constants.OP_IN, Value: nil}
	}
	return NewFilter(field, constants.OP_IN, validator.NormalizeFilterValueSlice(values))
}

// NewInFilterSlice 创建 IN 过滤条件（使用切片参数）
func NewInFilterSlice(field string, values []interface{}) *Filter {
	if values == nil {
		values = make([]interface{}, 0)
	}
	return NewFilter(field, constants.OP_IN, validator.NormalizeFilterValueSlice(values))
}

// NewLikeFilter 创建 LIKE 过滤条件
func NewLikeFilter(field string, value string) *Filter {
	return NewFilter(field, constants.OP_LIKE, "%"+value+"%")
}

// NewNeqFilter 创建不等于过滤条件
func NewNeqFilter(field string, value interface{}) *Filter {
	return NewFilter(field, constants.OP_NEQ, value)
}

// NewBetweenFilter 创建 BETWEEN 过滤条件
func NewBetweenFilter(field string, min, max interface{}) *Filter {
	return NewFilter(field, constants.OP_BETWEEN, []interface{}{min, max})
}

// NewNotInFilter 创建 NOT IN 过滤条件
func NewNotInFilter(field string, values ...interface{}) *Filter {
	if values == nil || (len(values) == 1 && values[0] == nil) {
		return &Filter{Field: field, Operator: constants.OP_NOT_IN, Value: nil}
	}
	return NewFilter(field, constants.OP_NOT_IN, validator.NormalizeFilterValueSlice(values))
}

// NewNotInFilterSlice 创建 NOT IN 过滤条件（使用切片参数）
func NewNotInFilterSlice(field string, values []interface{}) *Filter {
	if values == nil {
		values = make([]interface{}, 0)
	}
	return NewFilter(field, constants.OP_NOT_IN, validator.NormalizeFilterValueSlice(values))
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
	return NewFilter(field, constants.OP_STARTS_WITH, value)
}

// NewEndsWithFilter 创建后缀匹配过滤条件
func NewEndsWithFilter(field string, value string) *Filter {
	return NewFilter(field, constants.OP_ENDS_WITH, value)
}

// NewNotLikeFilter 创建 NOT LIKE 过滤条件
func NewNotLikeFilter(field string, value string) *Filter {
	return NewFilter(field, constants.OP_NOT_LIKE, "%"+value+"%")
}

// NewRegexpFilter 创建正则匹配过滤条件（仅MySQL/PostgreSQL）
func NewRegexpFilter(field string, pattern string) *Filter {
	return NewFilter(field, constants.OP_REGEX, pattern)
}

// NewFindInSetFilter 创建 FIND_IN_SET 过滤条件（MySQL特定）
func NewFindInSetFilter(field string, value interface{}) *Filter {
	return NewFilter(field, constants.OP_FIND_IN_SET, value)
}

// NewFilter 创建通用过滤条件(支持任意操作符)
func NewFilter(field string, operator constants.Operator, value interface{}) *Filter {
	return &Filter{Field: field, Operator: operator, Value: validator.NormalizeFilterValue(value)}
}
