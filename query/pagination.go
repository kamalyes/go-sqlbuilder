/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-11 00:00:00
 * @FilePath: \go-sqlbuilder\query\pagination.go
 * @Description: 分页相关结构
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package query

// PageBean 分页响应 - 借鉴 go-core 设计
type PageBean struct {
	Page     int         `json:"page"`      // 当前页码
	PageSize int         `json:"page_size"` // 每页数量
	Total    int64       `json:"total"`     // 总数
	Pages    int         `json:"pages"`     // 总页数
	Rows     interface{} `json:"rows"`      // 数据行
}

// NewPageBean 创建分页响应
func NewPageBean(page int, pageSize int, total int64, rows interface{}) *PageBean {
	pages := int((total + int64(pageSize) - 1) / int64(pageSize))
	if page < 1 {
		page = 1
	}
	if pages < 1 {
		pages = 1
	}
	return &PageBean{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
		Pages:    pages,
		Rows:     rows,
	}
}

// OrderBy 排序条件
type OrderBy struct {
	Field string // 字段名
	Order string // ASC / DESC
}

// NewOrderBy 创建排序条件
func NewOrderBy(field string, order string) *OrderBy {
	if order == "" {
		order = "ASC"
	}
	return &OrderBy{
		Field: field,
		Order: order,
	}
}
