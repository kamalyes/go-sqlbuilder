/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-11 14:49:28
 * @FilePath: \go-sqlbuilder\query\param.go
 * @Description: 高级查询参数 - 支持复杂过滤、排序、分页
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package query

import (
	"fmt"
	"strings"

	logger "github.com/kamalyes/go-logger"
	"github.com/kamalyes/go-sqlbuilder/constant"
)

// Param 高级查询参数 - 支持复杂过滤、排序、分页
type Param struct {
	Filters       []*Filter                 // 过滤条件
	FilterGroups  []*FilterGroup            // 过滤组（支持 AND/OR 混合）
	TimeRanges    map[string][2]interface{} // 时间范围 field -> [startTime, endTime]
	FindInSets    map[string][]string       // FIND_IN_SET 条件
	Orders        []*OrderBy                // 排序条件
	Page          int                       // 页码（从 1 开始）
	PageSize      int                       // 每页数量
	Offset        int                       // 偏移量（与 Page 互斥）
	Limit         int                       // 限制数量
	Distinct      bool                      // 是否去重
	SelectFields  []string                  // 指定查询字段
	HavingClauses []string                  // HAVING 子句
}

// NewParam 创建高级查询参数
func NewParam() *Param {
	return &Param{
		Filters:       make([]*Filter, 0),
		FilterGroups:  make([]*FilterGroup, 0),
		TimeRanges:    make(map[string][2]interface{}),
		FindInSets:    make(map[string][]string),
		Orders:        make([]*OrderBy, 0),
		Page:          1,
		PageSize:      10,
		SelectFields:  make([]string, 0),
		HavingClauses: make([]string, 0),
	}
}

// ==================== 基础过滤方法 ====================

// AddFilter 添加单个过滤条件
func (p *Param) AddFilter(field string, operator Operator, value interface{}) *Param {
	p.Filters = append(p.Filters, &Filter{
		Field:    field,
		Operator: operator,
		Value:    value,
		Logic:    string(constant.LOGIC_AND),
	})
	return p
}

// AddOrFilter 添加 OR 过滤条件
func (p *Param) AddOrFilter(field string, operator Operator, value interface{}) *Param {
	if len(p.Filters) > 0 {
		p.Filters[len(p.Filters)-1].Logic = string(constant.LOGIC_OR)
	}
	return p.AddFilter(field, operator, value)
}

// ==================== 便捷方法 ====================

// AddEQ 添加等于过滤
func (p *Param) AddEQ(field string, value interface{}) *Param {
	return p.AddFilter(field, OP_EQ, value)
}

// AddGT 添加大于过滤
func (p *Param) AddGT(field string, value interface{}) *Param {
	return p.AddFilter(field, OP_GT, value)
}

// AddGTE 添加大于等于过滤
func (p *Param) AddGTE(field string, value interface{}) *Param {
	return p.AddFilter(field, OP_GTE, value)
}

// AddLT 添加小于过滤
func (p *Param) AddLT(field string, value interface{}) *Param {
	return p.AddFilter(field, OP_LT, value)
}

// AddLTE 添加小于等于过滤
func (p *Param) AddLTE(field string, value interface{}) *Param {
	return p.AddFilter(field, OP_LTE, value)
}

// AddNEQ 添加不等于过滤
func (p *Param) AddNEQ(field string, value interface{}) *Param {
	return p.AddFilter(field, OP_NEQ, value)
}

// AddLike 添加全模糊 LIKE 过滤
func (p *Param) AddLike(field string, value string) *Param {
	return p.AddFilter(field, OP_LIKE, "%"+value+"%")
}

// AddStartsWith 添加前缀匹配 LIKE 过滤
func (p *Param) AddStartsWith(field string, value string) *Param {
	return p.AddFilter(field, OP_LIKE, value+"%")
}

// AddEndsWith 添加后缀匹配 LIKE 过滤
func (p *Param) AddEndsWith(field string, value string) *Param {
	return p.AddFilter(field, OP_LIKE, "%"+value)
}

// AddIn 添加 IN 过滤
func (p *Param) AddIn(field string, values ...interface{}) *Param {
	return p.AddFilter(field, OP_IN, values)
}

// ==================== OR 条件便捷方法 ====================

// AddOrEQ 添加 OR 等于过滤
func (p *Param) AddOrEQ(field string, value interface{}) *Param {
	return p.AddOrFilter(field, OP_EQ, value)
}

// AddOrGT 添加 OR 大于过滤
func (p *Param) AddOrGT(field string, value interface{}) *Param {
	return p.AddOrFilter(field, OP_GT, value)
}

// AddOrGTE 添加 OR 大于等于过滤
func (p *Param) AddOrGTE(field string, value interface{}) *Param {
	return p.AddOrFilter(field, OP_GTE, value)
}

// AddOrLT 添加 OR 小于过滤
func (p *Param) AddOrLT(field string, value interface{}) *Param {
	return p.AddOrFilter(field, OP_LT, value)
}

// AddOrLTE 添加 OR 小于等于过滤
func (p *Param) AddOrLTE(field string, value interface{}) *Param {
	return p.AddOrFilter(field, OP_LTE, value)
}

// AddOrNEQ 添加 OR 不等于过滤
func (p *Param) AddOrNEQ(field string, value interface{}) *Param {
	return p.AddOrFilter(field, OP_NEQ, value)
}

// AddOrLike 添加 OR 全模糊 LIKE 过滤
func (p *Param) AddOrLike(field string, value string) *Param {
	return p.AddOrFilter(field, OP_LIKE, "%"+value+"%")
}

// AddOrIn 添加 OR IN 过滤
func (p *Param) AddOrIn(field string, values ...interface{}) *Param {
	return p.AddOrFilter(field, OP_IN, values)
}

// ==================== 时间范围和特殊操作 ====================

// AddTimeRange 添加时间范围过滤
func (p *Param) AddTimeRange(field string, startTime, endTime interface{}) *Param {
	p.TimeRanges[field] = [2]interface{}{startTime, endTime}
	return p
}

// AddFindInSet 添加 FIND_IN_SET 过滤（MySQL 特定）
func (p *Param) AddFindInSet(field string, values ...string) *Param {
	p.FindInSets[field] = values
	return p
}

// ==================== 排序方法 ====================

// AddOrder 添加排序
func (p *Param) AddOrder(field string, order string) *Param {
	if order == "" {
		order = "ASC"
	}
	p.Orders = append(p.Orders, &OrderBy{Field: field, Order: strings.ToUpper(order)})
	return p
}

// AddOrderDesc 添加降序排序
func (p *Param) AddOrderDesc(field string) *Param {
	return p.AddOrder(field, "DESC")
}

// AddOrderAsc 添加升序排序
func (p *Param) AddOrderAsc(field string) *Param {
	return p.AddOrder(field, "ASC")
}

// ==================== 分页和其他选项 ====================

// SetPage 设置分页
func (p *Param) SetPage(page int, pageSize int) *Param {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	p.Page = page
	p.PageSize = pageSize
	p.Offset = (page - 1) * pageSize
	return p
}

// SetDistinct 设置去重
func (p *Param) SetDistinct(distinct bool) *Param {
	p.Distinct = distinct
	return p
}

// SetSelectFields 设置查询字段
func (p *Param) SetSelectFields(fields ...string) *Param {
	p.SelectFields = fields
	return p
}

// AddHaving 添加 HAVING 子句
func (p *Param) AddHaving(clause string) *Param {
	p.HavingClauses = append(p.HavingClauses, clause)
	return p
}

// ==================== WHERE 子句构建 ====================

// BuildWhereClause 构建 WHERE 子句 - 返回 WHERE SQL 片段和参数
func (p *Param) BuildWhereClause() (string, []interface{}) {
	var whereClauses []string
	var args []interface{}

	// 处理普通过滤
	for _, filter := range p.Filters {
		sql, filterArgs := p.buildFilterSQL(filter)
		whereClauses = append(whereClauses, sql)
		args = append(args, filterArgs...)
	}

	// 处理时间范围
	for field, timeRange := range p.TimeRanges {
		whereClauses = append(whereClauses, fmt.Sprintf("%s BETWEEN ? AND ?", field))
		args = append(args, timeRange[0], timeRange[1])
	}

	// 处理 FIND_IN_SET
	for field, values := range p.FindInSets {
		for _, value := range values {
			whereClauses = append(whereClauses, fmt.Sprintf("FIND_IN_SET(?, %s) > 0", field))
			args = append(args, value)
		}
	}

	if len(whereClauses) == 0 {
		logger.Debug("query.Param: no where clauses built")
		return "", args
	}

	return "WHERE " + strings.Join(whereClauses, " AND "), args
}

// buildFilterSQL 构建单个过滤 SQL
func (p *Param) buildFilterSQL(filter *Filter) (string, []interface{}) {
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
