/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-13 11:03:15
 * @FilePath: \go-sqlbuilder\constant\error.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package core

import (
	"github.com/kamalyes/go-sqlbuilder/constant"
)

type Filter struct {
	Field    string
	Operator constant.Operator
	Value    interface{}
}

type OrderBy struct {
	Field string
	Order string
}

type Pagination struct {
	Page     int
	PageSize int
	Offset   int
	Limit    int
}

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
