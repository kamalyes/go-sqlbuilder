/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-13 11:04:19
 * @FilePath: \go-sqlbuilder\constant\operator.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package constant

// Operator 操作符常量
type Operator string

const (
	// 比较操作符
	OP_EQ  Operator = "="
	OP_NEQ Operator = "!="
	OP_GT  Operator = ">"
	OP_GTE Operator = ">="
	OP_LT  Operator = "<"
	OP_LTE Operator = "<="

	// 字符串操作符
	OP_LIKE     Operator = "LIKE"
	OP_NOT_LIKE Operator = "NOT LIKE"

	// 集合操作符
	OP_IN      Operator = "IN"
	OP_NOT_IN  Operator = "NOT IN"
	OP_BETWEEN Operator = "BETWEEN"

	// 空值操作符
	OP_IS_NULL     Operator = "IS NULL"
	OP_IS_NOT_NULL Operator = "IS NOT NULL"

	// MySQL 特定
	OP_FIND_IN_SET Operator = "FIND_IN_SET"

	// 逻辑操作符
	LOGIC_AND Operator = "AND"
	LOGIC_OR  Operator = "OR"
)

// 兼容操作符映射
const (
	// go-core 兼容
	OP_COMPAT_LT  = "lt"
	OP_COMPAT_GT  = "gt"
	OP_COMPAT_LTE = "lte"
	OP_COMPAT_GTE = "gte"
	OP_COMPAT_EQ  = "eq"
	OP_COMPAT_NEQ = "neq"
	OP_COMPAT_LK  = "lk"
)

// OperatorMap 操作符映射表
var OperatorMap = map[string]Operator{
	"=":           OP_EQ,
	"!=":          OP_NEQ,
	">":           OP_GT,
	">=":          OP_GTE,
	"<":           OP_LT,
	"<=":          OP_LTE,
	"LIKE":        OP_LIKE,
	"NOT LIKE":    OP_NOT_LIKE,
	"IN":          OP_IN,
	"NOT IN":      OP_NOT_IN,
	"BETWEEN":     OP_BETWEEN,
	"IS NULL":     OP_IS_NULL,
	"IS NOT NULL": OP_IS_NOT_NULL,
	"FIND_IN_SET": OP_FIND_IN_SET,
}

// CompatOperatorMap 兼容操作符映射
var CompatOperatorMap = map[string]Operator{
	OP_COMPAT_EQ:  OP_EQ,
	OP_COMPAT_NEQ: OP_NEQ,
	OP_COMPAT_GT:  OP_GT,
	OP_COMPAT_GTE: OP_GTE,
	OP_COMPAT_LT:  OP_LT,
	OP_COMPAT_LTE: OP_LTE,
	OP_COMPAT_LK:  OP_LIKE,
}
