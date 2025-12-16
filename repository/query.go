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
	"github.com/kamalyes/go-toolbox/pkg/stringx"
	"github.com/kamalyes/go-toolbox/pkg/validator"
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

// SetPagination 设置分页条件 - WithPaging 的别名，提供更直观的API
func (q *Query) SetPagination(page, pageSize int) *Query {
	return q.WithPaging(page, pageSize)
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
	case int, int32, int64, uint, uint32, uint64:
		q.AddFilter(NewEqFilter(field, v))
	case bool:
		q.AddFilter(NewEqFilter(field, v))
	default:
		// 处理其他切片类型（如枚举切片）或单个值
		slice := convert.AnySliceToInterfaceSlice(v)
		if len(slice) > 0 {
			// 非空切片，添加 IN 过滤器
			q.AddFilter(NewInFilterSlice(field, slice))
		} else {
			// 检查是否是切片类型
			rv := reflect.ValueOf(v)
			if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
				// 空切片，不添加任何过滤器
				return q
			}
			// 非切片类型的其他值，添加等于过滤器
			q.AddFilter(NewEqFilter(field, v))
		}
	}
	return q
}

// AddLikeFilterIfNotEmpty 添加 LIKE 过滤条件（仅当关键词不为空时）
func (q *Query) AddLikeFilterIfNotEmpty(field, keyword string) *Query {
	if keyword != "" {
		q.AddFilter(NewLikeFilter(field, keyword))
	}
	return q
}

// isTimeValid 检查时间值是否有效(非nil且非零值)
func isTimeValid(timeVal interface{}) bool {
	if timeVal == nil {
		return false
	}

	// 定义 Unix 零点
	unixZero := time.Unix(0, 0)

	// 处理 *time.Time 类型
	if ptr, ok := timeVal.(*time.Time); ok {
		return ptr != nil && !ptr.IsZero() && ptr.After(unixZero)
	}

	// 处理 time.Time 类型
	if t, ok := timeVal.(time.Time); ok {
		return !t.IsZero() && t.After(unixZero)
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

	if slice := convert.AnySliceToInterfaceSlice(values); len(slice) > 0 {
		q.AddFilter(NewInFilterSlice(field, slice))
	}
	return q
}

// AddEqFilterIfNotEmpty 添加等于过滤条件（仅当值不为空时）
func (q *Query) AddEqFilterIfNotEmpty(field string, value interface{}) *Query {
	return q.AddFilterIfNotEmpty(field, value)
}

// AddNeqFilterIfNotEmpty 添加不等于过滤条件（仅当值不为空时）
func (q *Query) AddNeqFilterIfNotEmpty(field string, value interface{}) *Query {
	if !validator.IsEmptyValue(reflect.ValueOf(value)) {
		q.AddFilter(NewNeqFilter(field, value))
	}
	return q
}

// AddGtFilterIfNotEmpty 添加大于过滤条件（仅当值不为空时）
func (q *Query) AddGtFilterIfNotEmpty(field string, value interface{}) *Query {
	if !validator.IsEmptyValue(reflect.ValueOf(value)) {
		q.AddFilter(NewGtFilter(field, value))
	}
	return q
}

// AddGteFilterIfNotEmpty 添加大于等于过滤条件（仅当值不为空时）
func (q *Query) AddGteFilterIfNotEmpty(field string, value interface{}) *Query {
	if !validator.IsEmptyValue(reflect.ValueOf(value)) {
		q.AddFilter(NewGteFilter(field, value))
	}
	return q
}

// AddLtFilterIfNotEmpty 添加小于过滤条件（仅当值不为空时）
func (q *Query) AddLtFilterIfNotEmpty(field string, value interface{}) *Query {
	if !validator.IsEmptyValue(reflect.ValueOf(value)) {
		q.AddFilter(NewLtFilter(field, value))
	}
	return q
}

// AddLteFilterIfNotEmpty 添加小于等于过滤条件（仅当值不为空时）
func (q *Query) AddLteFilterIfNotEmpty(field string, value interface{}) *Query {
	if !validator.IsEmptyValue(reflect.ValueOf(value)) {
		q.AddFilter(NewLteFilter(field, value))
	}
	return q
}

// AddNotInFilterIfNotEmpty 添加 NOT IN 过滤条件（仅当切片不为空时）
func (q *Query) AddNotInFilterIfNotEmpty(field string, values interface{}) *Query {
	if values == nil {
		return q
	}
	if slice := convert.AnySliceToInterfaceSlice(values); len(slice) > 0 {
		q.AddFilter(NewNotInFilterSlice(field, slice))
	}
	return q
}

// AddBetweenFilterIfNotEmpty 添加 BETWEEN 过滤条件（仅当最小值和最大值都不为空时）
func (q *Query) AddBetweenFilterIfNotEmpty(field string, min, max interface{}) *Query {
	if !validator.IsEmptyValue(reflect.ValueOf(min)) && !validator.IsEmptyValue(reflect.ValueOf(max)) {
		q.AddFilter(NewBetweenFilter(field, min, max))
	}
	return q
}

// AddStartsWithFilterIfNotEmpty 添加前缀匹配过滤条件（仅当值不为空时）
func (q *Query) AddStartsWithFilterIfNotEmpty(field, value string) *Query {
	if !validator.IsEmptyValue(reflect.ValueOf(value)) {
		q.AddFilter(NewStartsWithFilter(field, value))
	}
	return q
}

// AddEndsWithFilterIfNotEmpty 添加后缀匹配过滤条件（仅当值不为空时）
func (q *Query) AddEndsWithFilterIfNotEmpty(field, value string) *Query {
	if !validator.IsEmptyValue(reflect.ValueOf(value)) {
		q.AddFilter(NewEndsWithFilter(field, value))
	}
	return q
}

// AddContainsFilterIfNotEmpty 添加包含匹配过滤条件（仅当值不为空时）
func (q *Query) AddContainsFilterIfNotEmpty(field, value string) *Query {
	return q.AddLikeFilterIfNotEmpty(field, value)
}

// AddNotLikeFilterIfNotEmpty 添加 NOT LIKE 过滤条件（仅当值不为空时）
func (q *Query) AddNotLikeFilterIfNotEmpty(field, value string) *Query {
	if !validator.IsEmptyValue(reflect.ValueOf(value)) {
		q.AddFilter(NewNotLikeFilter(field, value))
	}
	return q
}

// AddFindInSetFilterIfNotEmpty 添加 FIND_IN_SET 过滤条件（仅当值不为空时，MySQL特定）
func (q *Query) AddFindInSetFilterIfNotEmpty(field string, value interface{}) *Query {
	if !validator.IsEmptyValue(reflect.ValueOf(value)) {
		q.AddFilter(NewFindInSetFilter(field, value))
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
		return q.AddFilter(NewContainsFilter(field, keyword))
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

// SetDistinct 设置去重（简化版）
func (q *Query) SetDistinct() *Query {
	return q.WithDistinct(true)
}

// Page 设置分页（简化版）
func (q *Query) Page(page, pageSize int) *Query {
	return q.WithPaging(page, pageSize)
}

// Take 设置限制数量（简化版，相当于Limit）
func (q *Query) Take(limit int) *Query {
	return q.Limit(limit)
}

// Skip 设置跳过数量（简化版，相当于Offset）
func (q *Query) Skip(offset int) *Query {
	return q.Offset(offset)
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

// SelectOnly 只查询指定的单个字段（便捷方法）
func (q *Query) SelectOnly(field string) *Query {
	return q.Select(field)
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
			*conditions = append(*conditions, "("+groupCondition+")")
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
	return strings.Join(conditions, fmt.Sprintf(" %s ", group.LogicOp.String())), args
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
				*conditions = append(*conditions, "("+subCondition+")")
				*args = append(*args, subArgs...)
			}
		}
	}
}
