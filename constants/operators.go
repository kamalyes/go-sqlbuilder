/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-23 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-23 00:00:00
 * @FilePath: \go-sqlbuilder\constants\operators.go
 * @Description: 操作符常量定义
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package constants

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
