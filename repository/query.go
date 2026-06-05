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
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/kamalyes/go-sqlbuilder/constants"
	"github.com/kamalyes/go-toolbox/pkg/convert"
	"github.com/kamalyes/go-toolbox/pkg/mathx"
	"github.com/kamalyes/go-toolbox/pkg/stringx"
	"github.com/kamalyes/go-toolbox/pkg/syncx"
	"github.com/kamalyes/go-argus"
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

// NewQuery 创建查询条件，初始化所有切片字段
func NewQuery() *Query {
	return &Query{
		Filters:      make([]*Filter, 0),
		Orders:       make([]Order, 0),
		GroupBy:      make([]string, 0),
		Having:       make([]*Filter, 0),
		SelectFields: make([]string, 0),
		OmitFields:   make([]string, 0),
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
func (q *Query) WithPaging(page, pageSize interface{}) *Query {
	p, _ := convert.MustIntT[int](page, nil)
	ps, _ := convert.MustIntT[int](pageSize, nil)
	q.Pagination = &Pagination{
		Page:     mathx.IF(p <= 0, constants.DefaultPage, p),
		PageSize: mathx.IF(ps <= 0, constants.DefaultPageSize, ps),
	}
	return q
}

// SetPage 设置当前页码
func (q *Query) SetPage(page interface{}) *Query {
	if q.Pagination == nil {
		q.Pagination = &Pagination{}
	}
	q.Pagination.Page, _ = convert.MustIntT[int](page, nil)
	return q
}

// SetPageSize 设置每页记录数
func (q *Query) SetPageSize(pageSize interface{}) *Query {
	if q.Pagination == nil {
		q.Pagination = &Pagination{}
	}
	q.Pagination.PageSize, _ = convert.MustIntT[int](pageSize, nil)
	return q
}

// SetPagination 设置分页参数（页码和每页大小）
func (q *Query) SetPagination(page, pageSize interface{}) *Query {
	return q.WithPaging(page, pageSize)
}

// GetPagination 获取分页对象（如果不存在则创建默认分页）
func (q *Query) GetPagination() *Pagination {
	if q.Pagination == nil {
		q.Pagination = &Pagination{
			Page:     constants.DefaultPage,
			PageSize: constants.DefaultPageSize,
		}
	}
	return q.Pagination
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
func (q *Query) WithDistinct() *Query {
	q.Distinct = true
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
// 支持泛型自动处理不同类型的值，包括指针自动解引用
// 规则: nil/nil指针 -> 跳过, 布尔/数字(含零值) -> 有效, 空字符串/空切片 -> 跳过, 非空切片 -> IN
func (q *Query) AddFilterIfNotEmpty(field string, value interface{}) *Query {
	deref, empty := validator.NormalizeFilterValueIfNotEmpty(value)
	if empty {
		return q
	}
	rv := reflect.ValueOf(deref)

	if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
		slice := convert.AnySliceToInterfaceSlice(deref)
		if len(slice) > 0 {
			q.AddFilter(NewInFilterSlice(field, slice))
		}
		return q
	}

	q.AddFilter(NewEqFilter(field, deref))
	return q
}

// AddRawFilter 添加原始 SQL 过滤条件（用于复杂条件，如子查询、函数等）
// 注意：此方法直接将条件添加到 SQL 中，使用时需确保 SQL 注入安全
// 示例: query.AddRawFilter("to_agent_id IS NOT NULL AND to_agent_id != ”")
func (q *Query) AddRawFilter(rawCondition string) *Query {
	if rawCondition != "" {
		// 创建一个特殊的 Filter，使用 RAW 操作符
		q.AddFilter(&Filter{
			Field:    rawCondition,
			Operator: constants.OP_RAW,
			Value:    nil,
		})
	}
	return q
}

// AddLikeFilterIfNotEmpty 添加 LIKE 过滤条件（仅当关键词不为空时）
// keyword 支持 string 和 *string 类型
func (q *Query) AddLikeFilterIfNotEmpty(field string, keyword interface{}) *Query {
	deref, empty := validator.NormalizeFilterValueIfNotEmpty(keyword)
	if empty {
		return q
	}
	q.AddFilter(NewLikeFilter(field, fmt.Sprintf("%v", deref)))
	return q
}

// AddJsonbLikeFilterIfNotEmpty 添加 jsonb 字段文本搜索过滤条件（仅当关键词不为空时）
// 用于对 PostgreSQL jsonb 类型字段进行模糊搜索，自动将字段转为 text 后匹配
// keyword 支持 string 和 *wrapperspb.StringValue 等类型
func (q *Query) AddJsonbLikeFilterIfNotEmpty(field string, keyword interface{}) *Query {
	deref, empty := validator.NormalizeFilterValueIfNotEmpty(keyword)
	if empty {
		return q
	}
	q.AddFilter(NewJsonbLikeFilter(field, fmt.Sprintf("%v", deref)))
	return q
}

// AddTimeRangeFilter 添加时间范围过滤条件
// 自动过滤掉nil和零值时间，避免生成无效的SQL条件
func (q *Query) AddTimeRangeFilter(field string, startTime, endTime interface{}) *Query {
	if validator.IsTimeValid(startTime) {
		q.AddFilter(NewGteFilter(field, startTime))
	}
	if validator.IsTimeValid(endTime) {
		q.AddFilter(NewLteFilter(field, endTime))
	}
	return q
}

// AddInFilterIfNotEmpty 添加 IN 过滤条件(仅当切片不为空时)
// values 支持任意切片类型（[]string、[]int 等）以及对应指针类型
func (q *Query) AddInFilterIfNotEmpty(field string, values interface{}) *Query {
	deref, empty := validator.NormalizeFilterValueIfNotEmpty(values)
	if empty {
		return q
	}
	if slice := convert.AnySliceToInterfaceSlice(deref); len(slice) > 0 {
		q.AddFilter(NewInFilterSlice(field, slice))
	}
	return q
}

// AddNeqFilterIfNotEmpty 添加不等于过滤条件（仅当值不为空时）
func (q *Query) AddNeqFilterIfNotEmpty(field string, value interface{}) *Query {
	deref, empty := validator.NormalizeFilterValueIfNotEmpty(value)
	if empty {
		return q
	}
	q.AddFilter(NewNeqFilter(field, deref))
	return q
}

// AddGtFilterIfNotEmpty 添加大于过滤条件（仅当值不为空时）
func (q *Query) AddGtFilterIfNotEmpty(field string, value interface{}) *Query {
	deref, empty := validator.NormalizeFilterValueIfNotEmpty(value)
	if empty {
		return q
	}
	q.AddFilter(NewGtFilter(field, deref))
	return q
}

// AddGteFilterIfNotEmpty 添加大于等于过滤条件（仅当值不为空时）
func (q *Query) AddGteFilterIfNotEmpty(field string, value interface{}) *Query {
	deref, empty := validator.NormalizeFilterValueIfNotEmpty(value)
	if empty {
		return q
	}
	q.AddFilter(NewGteFilter(field, deref))
	return q
}

// AddLtFilterIfNotEmpty 添加小于过滤条件（仅当值不为空时）
func (q *Query) AddLtFilterIfNotEmpty(field string, value interface{}) *Query {
	deref, empty := validator.NormalizeFilterValueIfNotEmpty(value)
	if empty {
		return q
	}
	q.AddFilter(NewLtFilter(field, deref))
	return q
}

// AddLteFilterIfNotEmpty 添加小于等于过滤条件（仅当值不为空时）
func (q *Query) AddLteFilterIfNotEmpty(field string, value interface{}) *Query {
	deref, empty := validator.NormalizeFilterValueIfNotEmpty(value)
	if empty {
		return q
	}
	q.AddFilter(NewLteFilter(field, deref))
	return q
}

// AddCursorFilter 游标分页方向过滤（仅当 cursor 不为空时生效）
// isPrev=true 时使用 <（向前翻页），否则使用 >（向后翻页）
func (q *Query) AddCursorFilter(field string, cursor interface{}, isPrev bool) *Query {
	deref, empty := validator.NormalizeFilterValueIfNotEmpty(cursor)
	if empty {
		return q
	}
	if isPrev {
		return q.AddFilter(NewLtFilter(field, deref))
	}
	return q.AddFilter(NewGtFilter(field, deref))
}

// AddEqOrInFilter 单值用 =，多值自动转 IN（仅当切片不为空时生效）
// 适用于"按单个或多个 ID 过滤"的常见场景，避免手动判断长度
// values 支持任意切片类型（[]string、[]int 等）以及对应指针类型
// 示例: query.AddEqOrInFilter("session_id", sessionIDs)
func (q *Query) AddEqOrInFilter(field string, values interface{}) *Query {
	deref, empty := validator.NormalizeFilterValueIfNotEmpty(values)
	if empty {
		return q
	}
	slice := convert.AnySliceToInterfaceSlice(deref)
	switch len(slice) {
	case 0:
		return q
	case 1:
		return q.AddFilter(NewEqFilter(field, slice[0]))
	default:
		return q.AddFilter(NewInFilterSlice(field, slice))
	}
}

// AddNotInFilterIfNotEmpty 添加 NOT IN 过滤条件（仅当切片不为空时）
// values 支持任意切片类型（[]string、[]int 等）以及对应指针类型
func (q *Query) AddNotInFilterIfNotEmpty(field string, values interface{}) *Query {
	deref, empty := validator.NormalizeFilterValueIfNotEmpty(values)
	if empty {
		return q
	}
	if slice := convert.AnySliceToInterfaceSlice(deref); len(slice) > 0 {
		q.AddFilter(NewNotInFilterSlice(field, slice))
	}
	return q
}

// AddBetweenFilterIfNotEmpty 添加 BETWEEN 过滤条件（仅当最小值和最大值都不为空时）
func (q *Query) AddBetweenFilterIfNotEmpty(field string, min, max interface{}) *Query {
	minDeref, minEmpty := validator.NormalizeFilterValueIfNotEmpty(min)
	maxDeref, maxEmpty := validator.NormalizeFilterValueIfNotEmpty(max)
	if minEmpty || maxEmpty {
		return q
	}
	q.AddFilter(NewBetweenFilter(field, minDeref, maxDeref))
	return q
}

// AddStartsWithFilterIfNotEmpty 添加前缀匹配过滤条件（仅当值不为空时）
// value 支持 string 和 *string 类型
func (q *Query) AddStartsWithFilterIfNotEmpty(field string, value interface{}) *Query {
	deref, empty := validator.NormalizeFilterValueIfNotEmpty(value)
	if empty {
		return q
	}
	q.AddFilter(NewStartsWithFilter(field, fmt.Sprintf("%v", deref)))
	return q
}

// AddEndsWithFilterIfNotEmpty 添加后缀匹配过滤条件（仅当值不为空时）
// value 支持 string 和 *string 类型
func (q *Query) AddEndsWithFilterIfNotEmpty(field string, value interface{}) *Query {
	deref, empty := validator.NormalizeFilterValueIfNotEmpty(value)
	if empty {
		return q
	}
	q.AddFilter(NewEndsWithFilter(field, fmt.Sprintf("%v", deref)))
	return q
}

// AddNotLikeFilterIfNotEmpty 添加 NOT LIKE 过滤条件（仅当值不为空时）
// value 支持 string 和 *string 类型
func (q *Query) AddNotLikeFilterIfNotEmpty(field string, value interface{}) *Query {
	deref, empty := validator.NormalizeFilterValueIfNotEmpty(value)
	if empty {
		return q
	}
	q.AddFilter(NewNotLikeFilter(field, fmt.Sprintf("%v", deref)))
	return q
}

// AddFindInSetFilterIfNotEmpty 添加 FIND_IN_SET 过滤条件（仅当值不为空时，MySQL特定）
func (q *Query) AddFindInSetFilterIfNotEmpty(field string, value interface{}) *Query {
	deref, empty := validator.NormalizeFilterValueIfNotEmpty(value)
	if empty {
		return q
	}
	q.AddFilter(NewFindInSetFilter(field, deref))
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
//	query.AddSafeOrder(filter.SortBy, filter.SortOrder, "created_at", "DESC", []string{"created_at", "updated_at", "id"})
func (q *Query) AddSafeOrder(sortBy, sortOrder, defaultField, defaultDirection string, allowedFields ...[]string) *Query {
	field := defaultField
	direction := defaultDirection

	if sortBy != "" && validator.IsAllowedField(sortBy, allowedFields...) {
		field = sortBy
		direction = stringx.NormalizeSQLDirection(sortOrder, defaultDirection)
	}

	// 字段无效，使用默认字段和默认方向
	q.AddOrder(field, direction)
	return q
}

// AddEqual 添加等于条件
func (q *Query) AddEqual(field string, value interface{}) *Query {
	return q.AddFilter(NewEqFilter(field, value))
}

// AddNotEqual 添加不等于条件
func (q *Query) AddNotEqual(field string, value interface{}) *Query {
	return q.AddFilter(NewNeqFilter(field, value))
}

// AddLike 添加LIKE条件（包含匹配）
func (q *Query) AddLike(field, keyword string) *Query {
	if keyword != "" {
		return q.AddFilter(NewLikeFilter(field, keyword))
	}
	return q
}

// AddStartsWith 添加前缀匹配条件
func (q *Query) AddStartsWith(field, prefix string) *Query {
	if prefix != "" {
		return q.AddFilter(NewStartsWithFilter(field, prefix))
	}
	return q
}

// AddEndsWith 添加后缀匹配条件
func (q *Query) AddEndsWith(field, suffix string) *Query {
	if suffix != "" {
		return q.AddFilter(NewEndsWithFilter(field, suffix))
	}
	return q
}

// AddIn 添加IN条件
func (q *Query) AddIn(field string, values ...interface{}) *Query {
	if len(values) > 0 {
		return q.AddFilter(NewInFilter(field, values...))
	}
	return q
}

// AddNotIn 添加NOT IN条件
func (q *Query) AddNotIn(field string, values ...interface{}) *Query {
	if len(values) > 0 {
		return q.AddFilter(NewNotInFilter(field, values...))
	}
	return q
}

// AddGreaterThan 添加大于条件
func (q *Query) AddGreaterThan(field string, value interface{}) *Query {
	return q.AddFilter(NewGtFilter(field, value))
}

// AddGreaterEqual 添加大于等于条件
func (q *Query) AddGreaterEqual(field string, value interface{}) *Query {
	return q.AddFilter(NewGteFilter(field, value))
}

// AddLessThan 添加小于条件
func (q *Query) AddLessThan(field string, value interface{}) *Query {
	return q.AddFilter(NewLtFilter(field, value))
}

// AddLessEqual 添加小于等于条件
func (q *Query) AddLessEqual(field string, value interface{}) *Query {
	return q.AddFilter(NewLteFilter(field, value))
}

// AddBetween 添加BETWEEN条件
func (q *Query) AddBetween(field string, start, end interface{}) *Query {
	return q.AddFilter(NewBetweenFilter(field, start, end))
}

// AddIsNull 添加IS NULL条件
func (q *Query) AddIsNull(field string) *Query {
	return q.AddFilter(NewIsNullFilter(field))
}

// AddIsNotNull 添加IS NOT NULL条件
func (q *Query) AddIsNotNull(field string) *Query {
	return q.AddFilter(NewIsNotNullFilter(field))
}

// AddOrderAsc 添加升序排序
func (q *Query) AddOrderAsc(field string) *Query {
	return q.AddOrder(field, constants.Asc)
}

// AddOrderDesc 添加降序排序
func (q *Query) AddOrderDesc(field string) *Query {
	return q.AddOrder(field, constants.Desc)
}

// AddRawOrder 添加原始SQL排序表达式（用于复杂排序，如多字段排序、函数排序等）
// 注意：此方法直接将排序表达式添加到SQL中，使用时需确保SQL注入安全
func (q *Query) AddRawOrder(orderExpr string) *Query {
	if orderExpr != "" {
		q.Orders = append(q.Orders, Order{
			Field:     orderExpr,
			Direction: "", // 对于原始SQL，方向包含在表达式中
		})
	}
	return q
}

// AddTimeAfter 添加时间晚于条件
func (q *Query) AddTimeAfter(field string, t time.Time) *Query {
	if !t.IsZero() {
		return q.AddFilter(NewGtFilter(field, t))
	}
	return q
}

// AddTimeBefore 添加时间早于条件
func (q *Query) AddTimeBefore(field string, t time.Time) *Query {
	if !t.IsZero() {
		return q.AddFilter(NewLtFilter(field, t))
	}
	return q
}

// AddTimeBetween 添加时间范围条件
func (q *Query) AddTimeBetween(field string, start, end time.Time) *Query {
	if !start.IsZero() {
		q.AddFilter(NewGteFilter(field, start))
	}
	if !end.IsZero() {
		q.AddFilter(NewLteFilter(field, end))
	}
	return q
}

// AddToday 添加今天的条件（日期部分匹配）
func (q *Query) AddToday(field string) *Query {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := start.Add(24 * time.Hour).Add(-time.Nanosecond)
	return q.AddBetween(field, start, end)
}

// AddThisWeek 添加本周条件
func (q *Query) AddThisWeek(field string) *Query {
	now := time.Now()
	weekday := now.Weekday()
	if weekday == 0 {
		weekday = 7
	}
	start := now.AddDate(0, 0, -int(weekday-1))
	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	end := start.AddDate(0, 0, 7).Add(-time.Nanosecond)
	return q.AddBetween(field, start, end)
}

// AddThisMonth 添加本月条件
func (q *Query) AddThisMonth(field string) *Query {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	end := start.AddDate(0, 1, 0).Add(-time.Nanosecond)
	return q.AddBetween(field, start, end)
}

// AddThisYear 添加今年条件
func (q *Query) AddThisYear(field string) *Query {
	now := time.Now()
	start := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
	end := start.AddDate(1, 0, 0).Add(-time.Nanosecond)
	return q.AddBetween(field, start, end)
}

// Select 指定要查询的字段
// 示例: query.Select("id", "name", "email")
func (q *Query) Select(fields ...string) *Query {
	if len(fields) > 0 {
		q.SelectFields = fields
	}
	return q
}

// Omit 排除不需要查询的字段
// 示例: query.Omit("password", "secret")
func (q *Query) Omit(fields ...string) *Query {
	if len(fields) > 0 {
		q.OmitFields = fields
	}
	return q
}

// OmitSensitive 排除敏感字段的便捷方法
// 默认排除: password, secret, token, api_key
func (q *Query) OmitSensitive() *Query {
	return q.Omit("password", "secret", "token", "api_key", "access_token", "refresh_token")
}

// OmitLargeFields 排除大字段的便捷方法
// 默认排除: content, description, detail, data, payload, body
func (q *Query) OmitLargeFields() *Query {
	return q.Omit("content", "description", "detail", "data", "payload", "body", "remark")
}

// BuildWhereClause 构建WHERE子句和参数
// 返回值: whereClause (SQL条件字符串), args (参数数组)
// 示例: "agent_id = ? AND work_status = ? AND start_time >= ?", []interface{}{"user123", 2, time.Now()}
func (q *Query) BuildWhereClause() (string, []interface{}) {
	if q == nil {
		return "", nil
	}

	var allConditions []string
	var allArgs []interface{}

	// 处理简单过滤条件
	q.processSimpleFilters(&allConditions, &allArgs)

	// 处理复合过滤条件组
	q.processFilterGroup(&allConditions, &allArgs)

	// 如果没有任何条件
	if len(allConditions) == 0 {
		return "", nil
	}

	// 用 AND 连接所有顶级条件
	whereClause := strings.Join(allConditions, " AND ")
	return whereClause, allArgs
}

// processSimpleFilters 处理简单过滤条件
func (q *Query) processSimpleFilters(conditions *[]string, args *[]interface{}) {
	for _, filter := range q.Filters {
		if filter != nil {
			condition, arg := q.buildFilterCondition(filter)
			if condition != "" {
				*conditions = append(*conditions, condition)
				q.appendFilterArgs(args, arg)
			}
		}
	}
}

// processFilterGroup 处理复合过滤条件组
func (q *Query) processFilterGroup(conditions *[]string, args *[]interface{}) {
	if q.FilterGroup != nil && !q.FilterGroup.IsEmpty() {
		groupCondition, groupArgs := q.buildGroupCondition(q.FilterGroup)
		if groupCondition != "" {
			// 如果组内只有一个条件，不加括号
			if q.FilterGroup.Count() == 1 {
				*conditions = append(*conditions, groupCondition)
			} else {
				*conditions = append(*conditions, "("+groupCondition+")")
			}
			*args = append(*args, groupArgs...)
		}
	}
}

// appendFilterArgs 添加过滤条件参数
func (q *Query) appendFilterArgs(args *[]interface{}, arg interface{}) {
	if arg != nil {
		// 处理BETWEEN操作的多个参数
		if values, ok := arg.([]interface{}); ok {
			*args = append(*args, values...)
		} else {
			*args = append(*args, arg)
		}
	}
}

// buildFilterCondition 构建单个过滤条件的SQL和参数
func (q *Query) buildFilterCondition(filter *Filter) (string, interface{}) {
	if filter == nil {
		return "", nil
	}

	// 处理原始 SQL 条件
	if filter.Operator == constants.OP_RAW {
		return filter.Field, nil // 直接返回 Field 作为条件，不需要参数
	}

	// 处理特殊操作符
	if result, arg := q.handleSpecialOperators(filter); result != "" {
		return result, arg
	}

	// 通用处理：使用 map 查找模板
	if template, ok := constants.OperatorTemplateMap[filter.Operator]; ok {
		return fmt.Sprintf(template, filter.Field), filter.Value
	}

	return "", nil
}

// handleSpecialOperators 处理特殊操作符
func (q *Query) handleSpecialOperators(filter *Filter) (string, interface{}) {
	switch filter.Operator {
	case constants.OP_IS_NULL, constants.OP_IS_NOT_NULL:
		return q.handleNullOperators(filter)
	case constants.OP_BETWEEN:
		return q.handleBetweenOperator(filter)
	case constants.OP_STARTS_WITH:
		return q.handleStartsWithOperator(filter)
	case constants.OP_ENDS_WITH:
		return q.handleEndsWithOperator(filter)
	case constants.OP_CONTAINS:
		return q.handleContainsOperator(filter)
	case constants.OP_FIND_IN_SET:
		return q.handleFindInSetOperator(filter)
	}
	return "", nil
}

// handleNullOperators 处理 NULL 操作符
func (q *Query) handleNullOperators(filter *Filter) (string, interface{}) {
	if template, ok := constants.OperatorTemplateMap[filter.Operator]; ok {
		return fmt.Sprintf(template, filter.Field), nil
	}
	return "", nil
}

// handleBetweenOperator 处理 BETWEEN 操作符
func (q *Query) handleBetweenOperator(filter *Filter) (string, interface{}) {
	if values, ok := filter.Value.([]interface{}); ok && len(values) == 2 {
		return fmt.Sprintf(constants.SQL_BETWEEN, filter.Field), values
	}
	return "", nil
}

// handleStartsWithOperator 处理 STARTS_WITH 操作符
func (q *Query) handleStartsWithOperator(filter *Filter) (string, interface{}) {
	if valueStr, ok := filter.Value.(string); ok {
		return fmt.Sprintf(constants.SQL_LIKE, filter.Field), valueStr + constants.SQL_WILDCARD_ANY
	}
	return "", nil
}

// handleEndsWithOperator 处理 ENDS_WITH 操作符
func (q *Query) handleEndsWithOperator(filter *Filter) (string, interface{}) {
	if valueStr, ok := filter.Value.(string); ok {
		return fmt.Sprintf(constants.SQL_LIKE, filter.Field), constants.SQL_WILDCARD_ANY + valueStr
	}
	return "", nil
}

// handleContainsOperator 处理 CONTAINS 操作符
func (q *Query) handleContainsOperator(filter *Filter) (string, interface{}) {
	if valueStr, ok := filter.Value.(string); ok {
		return fmt.Sprintf(constants.SQL_LIKE, filter.Field), constants.SQL_WILDCARD_ANY + valueStr + constants.SQL_WILDCARD_ANY
	}
	return "", nil
}

// handleFindInSetOperator 处理 FIND_IN_SET 操作符
func (q *Query) handleFindInSetOperator(filter *Filter) (string, interface{}) {
	if template, ok := constants.OperatorTemplateMap[filter.Operator]; ok {
		return fmt.Sprintf(template, filter.Field) + " > 0", filter.Value
	}
	return "", nil
}

// buildGroupCondition 递归构建过滤组的条件
func (q *Query) buildGroupCondition(group *FilterGroup) (string, []interface{}) {
	if group == nil || group.IsEmpty() {
		return "", nil
	}

	var conditions []string
	var args []interface{}

	// 处理过滤条件
	q.processGroupFilters(group, &conditions, &args)

	// 递归处理子组
	q.processSubGroups(group, &conditions, &args)

	if len(conditions) == 0 {
		return "", nil
	}

	// 根据逻辑操作符连接条件
	result := strings.Join(conditions, fmt.Sprintf(" %s ", group.LogicOp.String()))

	// 如果只有一个条件，直接返回，不加括号
	if len(conditions) == 1 {
		return result, args
	}

	return result, args
}

// processGroupFilters 处理组内的过滤条件
func (q *Query) processGroupFilters(group *FilterGroup, conditions *[]string, args *[]interface{}) {
	for _, filter := range group.Filters {
		if filter != nil {
			condition, arg := q.buildFilterCondition(filter)
			if condition != "" {
				*conditions = append(*conditions, condition)
				q.appendFilterArgs(args, arg)
			}
		}
	}
}

// processSubGroups 递归处理子组
func (q *Query) processSubGroups(group *FilterGroup, conditions *[]string, args *[]interface{}) {
	for _, subGroup := range group.Groups {
		if subGroup != nil && !subGroup.IsEmpty() {
			subCondition, subArgs := q.buildGroupCondition(subGroup)
			if subCondition != "" {
				// 如果子组只有一个条件，不加括号
				if subGroup.Count() == 1 {
					*conditions = append(*conditions, subCondition)
				} else {
					*conditions = append(*conditions, "("+subCondition+")")
				}
				*args = append(*args, subArgs...)
			}
		}
	}
}

// HasPagination 检查是否设置了分页
func (q *Query) HasPagination() bool {
	return q.Pagination != nil
}

// HasGroupBy 检查是否设置了分组
func (q *Query) HasGroupBy() bool {
	return len(q.GroupBy) > 0
}

// HasHaving 检查是否设置了HAVING条件
func (q *Query) HasHaving() bool {
	return len(q.Having) > 0
}

// HasOrders 检查是否设置了排序
func (q *Query) HasOrders() bool {
	return len(q.Orders) > 0
}

// HasSelectFields 检查是否指定了查询字段
func (q *Query) HasSelectFields() bool {
	return len(q.SelectFields) > 0
}

// HasOmitFields 检查是否指定了排除字段
func (q *Query) HasOmitFields() bool {
	return len(q.OmitFields) > 0
}

// IsLimited 检查是否设置了查询限制数量
func (q *Query) IsLimited() bool {
	return q.LimitValue != nil
}

// IsOffset 检查是否设置了查询偏移量
func (q *Query) IsOffset() bool {
	return q.OffsetValue != nil
}

// ResetFilters 重置所有过滤条件（包括简单过滤器和复合过滤组）
func (q *Query) ResetFilters() *Query {
	q.Filters = make([]*Filter, 0)
	q.FilterGroup = nil
	return q
}

// ResetOrders 重置所有排序条件
func (q *Query) ResetOrders() *Query {
	q.Orders = make([]Order, 0)
	return q
}

// ResetPagination 重置分页设置
func (q *Query) ResetPagination() *Query {
	q.Pagination = nil
	return q
}

// Clone 深拷贝查询对象，创建一个新的独立副本
func (q *Query) Clone() *Query {
	cloned := NewQuery()
	if q == nil {
		return cloned
	}

	if err := syncx.DeepCopy(&cloned, &q); err != nil {
		return NewQuery()
	}

	if cloned.Filters == nil {
		cloned.Filters = make([]*Filter, 0)
	}
	if cloned.Orders == nil {
		cloned.Orders = make([]Order, 0)
	}
	if cloned.GroupBy == nil {
		cloned.GroupBy = make([]string, 0)
	}
	if cloned.Having == nil {
		cloned.Having = make([]*Filter, 0)
	}
	if cloned.SelectFields == nil {
		cloned.SelectFields = make([]string, 0)
	}
	if cloned.OmitFields == nil {
		cloned.OmitFields = make([]string, 0)
	}

	return cloned
}
