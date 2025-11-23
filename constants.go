/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-23 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-23 22:50:00
 * @FilePath: \go-sqlbuilder\constants.go
 * @Description: 常量定义 - 操作符、分页、批处理等配置
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package sqlbuilder

// Operator 操作符类型（统一定义，避免重复）
type Operator string

// ==================== 比较操作符 ====================

const (
	OpEqual              Operator = "="  // 等于
	OpNotEqual           Operator = "!=" // 不等于
	OpGreaterThan        Operator = ">"  // 大于
	OpGreaterThanOrEqual Operator = ">=" // 大于等于
	OpLessThan           Operator = "<"  // 小于
	OpLessThanOrEqual    Operator = "<=" // 小于等于
)

// ==================== 字符串操作符 ====================

const (
	OpLike    Operator = "LIKE"     // 模糊匹配
	OpNotLike Operator = "NOT LIKE" // 不匹配
)

// ==================== 集合操作符 ====================

const (
	OpIn      Operator = "IN"      // 包含
	OpNotIn   Operator = "NOT IN"  // 不包含
	OpBetween Operator = "BETWEEN" // 范围
)

// ==================== 空值操作符 ====================

const (
	OpIsNull    Operator = "IS NULL"     // 为空
	OpIsNotNull Operator = "IS NOT NULL" // 不为空
)

// ==================== 数据库特定操作符 ====================

const (
	OpFindInSet Operator = "FIND_IN_SET" // MySQL FIND_IN_SET
	OpRegex     Operator = "REGEXP"      // 正则匹配
	OpNotRegex  Operator = "NOT REGEXP"  // 正则不匹配
)

// ==================== 逻辑操作符 ====================

const (
	LogicAnd Operator = "AND" // 逻辑与
	LogicOr  Operator = "OR"  // 逻辑或
	LogicNot Operator = "NOT" // 逻辑非
)

// ==================== 向后兼容别名 ====================

const (
	OP_EQ          = OpEqual
	OP_NEQ         = OpNotEqual
	OP_GT          = OpGreaterThan
	OP_GTE         = OpGreaterThanOrEqual
	OP_LT          = OpLessThan
	OP_LTE         = OpLessThanOrEqual
	OP_LIKE        = OpLike
	OP_NOT_LIKE    = OpNotLike
	OP_IN          = OpIn
	OP_NOT_IN      = OpNotIn
	OP_BETWEEN     = OpBetween
	OP_IS_NULL     = OpIsNull
	OP_IS_NOT_NULL = OpIsNotNull
	OP_FIND_IN_SET = OpFindInSet
	LOGIC_AND      = LogicAnd
	LOGIC_OR       = LogicOr
)

// ==================== 分页默认值 ====================

const (
	// DefaultPage 默认页码
	DefaultPage = 1

	// DefaultPageSize 默认每页大小
	DefaultPageSize = 10
)

// ==================== 批处理默认值 ====================

const (
	// DefaultBatchSize 默认批处理大小
	DefaultBatchSize = 100
)

// ==================== 超时默认值 ====================

const (
	// DefaultQueryTimeout 默认查询超时时间（秒）
	DefaultQueryTimeout = 30
)

// ==================== SQL 模板常量 ====================

const (
	// SQL_EQUAL 等于条件模板
	SQL_EQUAL = "%s = ?"

	// SQL_NOT_EQUAL 不等于条件模板
	SQL_NOT_EQUAL = "%s != ?"

	// SQL_GREATER 大于条件模板
	SQL_GREATER = "%s > ?"

	// SQL_GREATER_EQUAL 大于等于条件模板
	SQL_GREATER_EQUAL = "%s >= ?"

	// SQL_LESS 小于条件模板
	SQL_LESS = "%s < ?"

	// SQL_LESS_EQUAL 小于等于条件模板
	SQL_LESS_EQUAL = "%s <= ?"

	// SQL_IN IN条件模板
	SQL_IN = "%s IN ?"

	// SQL_NOT_IN NOT IN条件模板
	SQL_NOT_IN = "%s NOT IN ?"

	// SQL_LIKE LIKE条件模板
	SQL_LIKE = "%s LIKE ?"

	// SQL_NOT_LIKE NOT LIKE条件模板
	SQL_NOT_LIKE = "%s NOT LIKE ?"

	// SQL_BETWEEN BETWEEN条件模板
	SQL_BETWEEN = "%s BETWEEN ? AND ?"

	// SQL_IS_NULL IS NULL条件模板
	SQL_IS_NULL = "%s IS NULL"

	// SQL_IS_NOT_NULL IS NOT NULL条件模板
	SQL_IS_NOT_NULL = "%s IS NOT NULL"

	// SQL_FIND_IN_SET FIND_IN_SET条件模板
	SQL_FIND_IN_SET = "FIND_IN_SET(?, %s)"

	// SQL_ORDER_BY 排序模板
	SQL_ORDER_BY = "%s %s"

	// SQL_INCREMENT 字段自增模板
	SQL_INCREMENT = "%s + ?"

	// SQL_DECREMENT 字段自减模板
	SQL_DECREMENT = "%s - ?"
)

// ==================== 排序方向常量 ====================

const (
	// Asc 升序
	Asc = "ASC"

	// Desc 降序
	Desc = "DESC"

	// DefaultOrder 默认排序方向
	DefaultOrder = Asc
)
