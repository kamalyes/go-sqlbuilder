/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-13 23:19:22
 * @FilePath: \go-sqlbuilder\meta\paging.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package meta

// Paging 分页元数据
type Paging struct {
	Page     int32 `json:"page"`      // 当前页码（从1开始）
	PageSize int32 `json:"page_size"` // 每页记录数
	Offset   int32 `json:"offset"`    // 数据库偏移量
	Limit    int32 `json:"limit"`     // 查询限制数
	Total    int64 `json:"total"`     // 总记录数
}

// GetOffset 计算数据库偏移量
func (p *Paging) GetOffset() int {
	if p.Page <= 0 {
		p.Page = 1
	}
	return int((p.Page - 1) * p.PageSize)
}

// GetLimit 获取查询限制数
func (p *Paging) GetLimit() int {
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
	return int(p.PageSize)
}

// GetTotalPages 计算总页数
func (p *Paging) GetTotalPages() int64 {
	if p.PageSize <= 0 {
		return 0
	}
	return (p.Total + int64(p.PageSize) - 1) / int64(p.PageSize)
}

// HasNextPage 是否有下一页
func (p *Paging) HasNextPage() bool {
	return p.Page < int32(p.GetTotalPages())
}

// HasPrevPage 是否有上一页
func (p *Paging) HasPrevPage() bool {
	return p.Page > 1
}
