/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-23 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-23 22:50:00
 * @FilePath: \go-sqlbuilder\paging.go
 * @Description: 分页工具 - Pagination分页元数据和辅助方法
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package sqlbuilder

// Pagination 分页元数据
type Pagination struct {
	Page     int32 `json:"page"`      // 当前页码（从1开始）
	PageSize int32 `json:"page_size"` // 每页记录数
	Offset   int32 `json:"offset"`    // 数据库偏移量
	Limit    int32 `json:"limit"`     // 查询限制数
	Total    int64 `json:"total"`     // 总记录数
}

// GetOffset 计算数据库偏移量
func (p *Pagination) GetOffset() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	return int((p.Page - 1) * p.PageSize)
}

// GetLimit 获取查询限制数
func (p *Pagination) GetLimit() int {
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	return int(p.PageSize)
}

// GetTotalPages 计算总页数
func (p *Pagination) GetTotalPages() int64 {
	if p.PageSize <= 0 {
		return 0
	}
	return (p.Total + int64(p.PageSize) - 1) / int64(p.PageSize)
}

// HasNextPage 是否有下一页
func (p *Pagination) HasNextPage() bool {
	return p.Page < int32(p.GetTotalPages())
}

// HasPrevPage 是否有上一页
func (p *Pagination) HasPrevPage() bool {
	return p.Page > 1
}
