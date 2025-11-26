/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-23 13:29:12
 * @FilePath: \go-sqlbuilder\repository\query.go
 * @Description: 类型定义 - Filter、Pagination、QueryCondition 等
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package repository

import (
	"time"

	"github.com/kamalyes/go-sqlbuilder/constants"
)

type QueryCondition struct {
	Filters    []Filter
	Orders     []OrderBy
	Pagination *Pagination
}

type ComparedColumn struct {
	TableAlias  string
	ColumnName  string
	DBFieldName string
}

type ComparedValue struct {
	Value interface{}
	IsRaw bool
}

// NewQuery 创建查询条件
func NewQuery() *Query {
	return &Query{
		Filters: make([]*Filter, 0),
		Orders:  make([]Order, 0),
	}
}

// FindOptions 兼容旧API的查询选项结构
type FindOptions struct {
	Conditions []Condition
	Orders     []OrderBy
	Limit      int
	Offset     int
}

// Condition 查询条件结构
type Condition struct {
	Field string
	Op    constants.Operator
	Value interface{}
}

// OrderBy 排序条件结构
type OrderBy struct {
	Field     string
	Direction string
}

// AddFilter 添加过滤条件
func (q *Query) AddFilter(filter *Filter) *Query {
	if filter != nil {
		q.Filters = append(q.Filters, filter)
	}
	return q
}

// AddFilters 批量添加过滤条件
func (q *Query) AddFilters(filters ...*Filter) *Query {
	for _, f := range filters {
		q.AddFilter(f)
	}
	return q
}

// AddOrder 添加排序条件
func (q *Query) AddOrder(field, direction string) *Query {
	q.Orders = append(q.Orders, Order{Field: field, Direction: direction})
	return q
}

// WithPaging 设置分页条件
func (q *Query) WithPaging(page, pageSize int) *Query {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	q.Pagination = &Pagination{
		Page:     int32(page),
		PageSize: int32(pageSize),
	}
	return q
}

// Limit 设置查询限制数量
func (q *Query) Limit(limit int) *Query {
	q.LimitValue = &limit
	return q
}

// Offset 设置查询偏移量
func (q *Query) Offset(offset int) *Query {
	q.OffsetValue = &offset
	return q
}

// WithFilterGroup 设置复合过滤条件组
func (q *Query) WithFilterGroup(group *FilterGroup) *Query {
	q.FilterGroup = group
	return q
}

// WithDistinct 设置去重查询
func (q *Query) WithDistinct(distinct bool) *Query {
	q.Distinct = distinct
	return q
}

// AddGroupBy 添加分组字段
func (q *Query) AddGroupBy(fields ...string) *Query {
	if q.GroupBy == nil {
		q.GroupBy = make([]string, 0)
	}
	q.GroupBy = append(q.GroupBy, fields...)
	return q
}

// AddHaving 添加HAVING条件
func (q *Query) AddHaving(filter *Filter) *Query {
	if q.Having == nil {
		q.Having = make([]*Filter, 0)
	}
	if filter != nil {
		q.Having = append(q.Having, filter)
	}
	return q
}

// HasFilters 检查是否有过滤条件
func (q *Query) HasFilters() bool {
	return len(q.Filters) > 0 || (q.FilterGroup != nil && !q.FilterGroup.IsEmpty())
}

// GetAllFilters 获取所有过滤条件（扁平化）
func (q *Query) GetAllFilters() []*Filter {
	allFilters := make([]*Filter, 0)
	allFilters = append(allFilters, q.Filters...)

	if q.FilterGroup != nil {
		allFilters = append(allFilters, q.flattenFilters(q.FilterGroup)...)
	}

	return allFilters
}

// flattenFilters 递归扁平化过滤条件组
func (q *Query) flattenFilters(group *FilterGroup) []*Filter {
	filters := make([]*Filter, 0)
	filters = append(filters, group.Filters...)

	for _, subGroup := range group.Groups {
		filters = append(filters, q.flattenFilters(subGroup)...)
	}

	return filters
}

// AddFilterIfNotEmpty 添加过滤条件（仅当值不为空时）
// 支持泛型自动处理不同类型的值
func (q *Query) AddFilterIfNotEmpty(field string, value interface{}) *Query {
	if value == nil {
		return q
	}

	switch v := value.(type) {
	case string:
		if v != "" {
			q.AddFilter(NewEqFilter(field, v))
		}
	case *string:
		if v != nil && *v != "" {
			q.AddFilter(NewEqFilter(field, *v))
		}
	case []string:
		if len(v) > 0 {
			values := make([]interface{}, len(v))
			for i, s := range v {
				values[i] = s
			}
			q.AddFilter(NewInFilterSlice(field, values))
		}
	case []int:
		if len(v) > 0 {
			values := make([]interface{}, len(v))
			for i, n := range v {
				values[i] = n
			}
			q.AddFilter(NewInFilterSlice(field, values))
		}
	case []int32:
		if len(v) > 0 {
			values := make([]interface{}, len(v))
			for i, n := range v {
				values[i] = n
			}
			q.AddFilter(NewInFilterSlice(field, values))
		}
	case []int64:
		if len(v) > 0 {
			values := make([]interface{}, len(v))
			for i, n := range v {
				values[i] = n
			}
			q.AddFilter(NewInFilterSlice(field, values))
		}
	case int, int32, int64, uint, uint32, uint64:
		q.AddFilter(NewEqFilter(field, v))
	case bool:
		q.AddFilter(NewEqFilter(field, v))
	default:
		// 处理枚举类型（通过反射）
		// 尝试将切片类型转换为 []interface{}
		if IsSliceType(v) {
			if slice := ConvertToInterfaceSlice(v); len(slice) > 0 {
				q.AddFilter(NewInFilterSlice(field, slice))
			}
		} else {
			// 单个枚举值或其他类型
			q.AddFilter(NewEqFilter(field, v))
		}
	}
	return q
}

// AddLikeFilterIfNotEmpty 添加 LIKE 过滤条件（仅当关键词不为空时）
func (q *Query) AddLikeFilterIfNotEmpty(field, keyword string) *Query {
	if keyword != "" {
		q.AddFilter(NewContainsFilter(field, keyword))
	}
	return q
}

// isTimeValid 检查时间值是否有效(非nil且非零值)
func isTimeValid(timeVal interface{}) bool {
	if timeVal == nil {
		return false
	}

	// 处理 *time.Time 类型
	if ptr, ok := timeVal.(*time.Time); ok {
		return ptr != nil && !ptr.IsZero()
	}

	// 处理 time.Time 类型
	if t, ok := timeVal.(time.Time); ok {
		return !t.IsZero()
	}

	// 其他类型认为有效
	return true
}

// AddTimeRangeFilter 添加时间范围过滤条件
// 自动过滤掉nil和零值时间，避免生成无效的SQL条件
func (q *Query) AddTimeRangeFilter(field string, startTime, endTime interface{}) *Query {
	if isTimeValid(startTime) {
		q.AddFilter(NewGteFilter(field, startTime))
	}
	if isTimeValid(endTime) {
		q.AddFilter(NewLteFilter(field, endTime))
	}
	return q
}

// AddInFilterIfNotEmpty 添加 IN 过滤条件(仅当切片不为空时)
func (q *Query) AddInFilterIfNotEmpty(field string, values interface{}) *Query {
	if values == nil {
		return q
	}

	if slice := ConvertToInterfaceSlice(values); len(slice) > 0 {
		q.AddFilter(NewInFilterSlice(field, slice))
	}
	return q
}

// AddSafeOrder 安全地添加排序条件
// 参数:
//   - sortBy: 排序字段(可选,为空时使用defaultField)
//   - sortOrder: 排序方向(可选,仅支持"ASC"/"DESC",为空时使用defaultDirection)
//   - defaultField: 默认排序字段
//   - defaultDirection: 默认排序方向
//   - allowedFields: 允许的字段白名单(可选,为空则不限制)
//
// 示例:
//
//	query.AddSafeOrder(filter.SortBy, filter.SortOrder, "created_at", constants.Desc, []string{"created_at", "updated_at", "id"})
func (q *Query) AddSafeOrder(sortBy, sortOrder, defaultField, defaultDirection string, allowedFields ...[]string) *Query {
	// 1. 处理排序字段
	field := defaultField
	fieldValid := false // 标记字段是否有效

	if sortBy != "" {
		// 如果提供了白名单,检查字段是否在白名单中
		if len(allowedFields) > 0 && len(allowedFields[0]) > 0 {
			for _, allowedField := range allowedFields[0] {
				if sortBy == allowedField {
					field = sortBy
					fieldValid = true
					break
				}
			}
		} else {
			// 没有白名单,但要验证字段名是否安全(仅允许字母、数字、下划线)
			if isSafeFieldName(sortBy) {
				field = sortBy
				fieldValid = true
			}
		}
	}

	// 2. 处理排序方向(标准化为大写)
	direction := defaultDirection
	// 只有当字段有效时,才处理自定义的排序方向
	if fieldValid && sortOrder != "" {
		upperOrder := ""
		for _, ch := range sortOrder {
			if ch >= 'a' && ch <= 'z' {
				upperOrder += string(ch - 32)
			} else if ch >= 'A' && ch <= 'Z' {
				upperOrder += string(ch)
			}
		}
		if upperOrder == "ASC" || upperOrder == "DESC" {
			direction = upperOrder
		}
	}

	// 3. 添加排序
	q.AddOrder(field, direction)
	return q
}

// isSafeFieldName 检查字段名是否安全(仅包含字母、数字、下划线、点号)
func isSafeFieldName(field string) bool {
	if field == "" {
		return false
	}
	for _, ch := range field {
		if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') || ch == '_' || ch == '.') {
			return false
		}
	}
	return true
}
