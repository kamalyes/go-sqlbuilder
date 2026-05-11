/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2026-05-11 02:00:15
 * @FilePath: \go-sqlbuilder\repository\scope.go
 * @Description: 类型定义 - Filter、Pagination、QueryCondition 等
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package repository

// ScopeHook 作用域钩子接口
type ScopeHook interface {
	ApplyScope(query *Query) *Query
}

// WithScope 为查询添加作用域钩子
func (q *Query) WithScope(hook ScopeHook) *Query {
	if hook != nil {
		return hook.ApplyScope(q)
	}
	return q
}
