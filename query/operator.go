/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-11 00:00:00
 * @FilePath: \go-sqlbuilder\query\operator.go
 * @Description: 查询操作符和常量定义
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package query

// Operator 查询操作符
type Operator string

const (
	OP_EQ          Operator = "="           // 等于
	OP_NEQ         Operator = "!="          // 不等于
	OP_GT          Operator = ">"           // 大于
	OP_GTE         Operator = ">="          // 大于等于
	OP_LT          Operator = "<"           // 小于
	OP_LTE         Operator = "<="          // 小于等于
	OP_LIKE        Operator = "LIKE"        // 模糊匹配
	OP_IN          Operator = "IN"          // 包含
	OP_BETWEEN     Operator = "BETWEEN"     // 范围
	OP_IS_NULL     Operator = "IS NULL"     // 为空
	OP_FIND_IN_SET Operator = "FIND_IN_SET" // MySQL FIND_IN_SET
)

// 比较操作符常量（兼容 go-core）
const (
	QP_LT  = "lt"  // 小于
	QP_GT  = "gt"  // 大于
	QP_LTE = "lte" // 小于等于
	QP_GTE = "gte" // 大于等于
	QP_EQ  = "eq"  // 等于
	QP_NEQ = "neq" // 不等于
	QP_LK  = "lk"  // LIKE
)

// 排序操作符常量（兼容 go-core）
const (
	QP_PD = "pd" // 降序（descending）
	QP_PA = "pa" // 升序（ascending）
)

// 或条件操作符常量（兼容 go-core）
const (
	QP_ORLK = "orlk"   // OR LIKE
	QP_ORLT = "orlt"   // OR LT
	QP_ORGT = "orgt"   // OR GT
)
