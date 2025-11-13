/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-13 11:04:25
 * @FilePath: \go-sqlbuilder\constant\query.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package constant

const (
	ParamPage     = "page"
	ParamPageSize = "pageSize"
	ParamOffset   = "offset"
	ParamLimit    = "limit"
	ParamSort     = "sort"
	ParamOrder    = "order"
	ParamFilter   = "filter"
	ParamSearch   = "search"
	ParamFields   = "fields"
	ParamInclude  = "include"
	ParamExclude  = "exclude"
)

const (
	ConditionKeyFilters    = "filters"
	ConditionKeyOrders     = "orders"
	ConditionKeyPagination = "pagination"
)

const (
	MinPageNumber = 1
	MinPageSize   = 1
	MaxPageSize   = 1000
	DefaultPage   = 1
)
