/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-25 08:06:51
 * @FilePath: \go-sqlbuilder\repository\paging.go
 * @Description: 分页工具 - Pagination分页元数据和辅助方法
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package repository

import (
	"time"

	"github.com/kamalyes/go-sqlbuilder/constants"
	"github.com/kamalyes/go-toolbox/pkg/mathx"
	"github.com/kamalyes/go-toolbox/pkg/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// PaginationT 泛型分页元数据
type PaginationT[T types.Integer] struct {
	Page     T     `json:"page"`      // 当前页码（从1开始）
	PageSize T     `json:"page_size"` // 每页记录数
	Total    int64 `json:"total"`     // 总记录数（固定 int64 以支持大数据量）
}

// Pagination 默认分页类型
type Pagination = PaginationT[int32]

// Pagination64 int64 版本的分页类型（用于需要大数值的场景）
type Pagination64 = PaginationT[int64]

// GetOffset 计算数据库偏移量
func (p *PaginationT[T]) GetOffset() int {
	p.Page = mathx.IF(p.Page <= 0, T(constants.DefaultPage), p.Page)
	return int((p.Page - 1) * p.PageSize)
}

// GetLimit 获取查询限制数（自动应用默认值和最大值限制）
func (p *PaginationT[T]) GetLimit() int {
	// 应用默认值和范围限制
	p.PageSize = mathx.IfDefaultAndClamp(p.PageSize, T(constants.DefaultPageSize), T(constants.MinPageSize), T(constants.MaxPageSize))
	return int(p.PageSize)
}

// GetTotalPages 计算总页数
func (p *PaginationT[T]) GetTotalPages() T {
	if p.PageSize <= 0 {
		return 0
	}
	return T((p.Total + int64(p.PageSize) - 1) / int64(p.PageSize))
}

// HasNextPage 是否有下一页
func (p *PaginationT[T]) HasNextPage() bool {
	return p.Page < p.GetTotalPages()
}

// HasPrevPage 是否有上一页
func (p *PaginationT[T]) HasPrevPage() bool {
	return p.Page > 1
}

// IsToday 判断时间是否为今天（支持 *time.Time）
func IsToday(t *time.Time) bool {
	if t == nil {
		return true // 未指定时间，默认为今天
	}

	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return t.After(todayStart) || t.Equal(todayStart)
}

// IsTodayFromProto 判断 protobuf Timestamp 是否为今天
func IsTodayFromProto(t *timestamppb.Timestamp) bool {
	if t == nil {
		return true // 未指定时间，默认为今天
	}

	tt := t.AsTime()
	return IsToday(&tt)
}

// IsTodayRange 判断时间范围是否包含今天
func IsTodayRange(startTime, endTime *time.Time) bool {
	// 如果都未指定，默认为今天
	if startTime == nil && endTime == nil {
		return true
	}

	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEnd := todayStart.Add(24 * time.Hour)

	// 如果开始时间在今天范围内
	if IsTimeInTodayRange(startTime, todayStart, todayEnd) {
		return true
	}

	// 如果结束时间在今天范围内
	if IsTimeInTodayRange(endTime, todayStart, todayEnd) {
		return true
	}

	// 如果时间范围跨越今天
	if startTime != nil && endTime != nil && startTime.Before(todayStart) && endTime.After(todayEnd) {
		return true
	}

	return false
}

// IsTimeInTodayRange 判断时间是否在今天范围内
func IsTimeInTodayRange(t *time.Time, todayStart, todayEnd time.Time) bool {
	if t == nil {
		return false
	}
	return t.Before(todayEnd) && (t.After(todayStart) || t.Equal(todayStart))
}

// GetTodayRange 获取今天的时间范围（00:00:00 - 23:59:59）
func GetTodayRange() (startTime, endTime time.Time) {
	now := time.Now()
	startTime = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endTime = startTime.Add(24 * time.Hour)
	return
}
