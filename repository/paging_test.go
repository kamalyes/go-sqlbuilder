/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-23 23:56:05
 * @FilePath: \go-sqlbuilder\paging_test.go
 * @Description: 分页工具 - Pagination分页元数据和辅助方法
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package repository

import (
	"github.com/kamalyes/go-sqlbuilder/constants"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// TestPagination_GetOffset 测试 Pagination GetOffset 方法
func TestPagination_GetOffset(t *testing.T) {
	// 测试正常分页
	page1 := &Pagination{Page: 1, PageSize: 10}
	assert.Equal(t, 0, page1.GetOffset())

	page2 := &Pagination{Page: 2, PageSize: 10}
	assert.Equal(t, 10, page2.GetOffset())

	page5 := &Pagination{Page: 5, PageSize: 20}
	assert.Equal(t, 80, page5.GetOffset())

	// 测试边界情况
	pageZero := &Pagination{Page: 0, PageSize: 10}
	assert.Equal(t, 0, pageZero.GetOffset()) // Page<=0 会被设为 1

	pageNegative := &Pagination{Page: -1, PageSize: 10}
	assert.Equal(t, 0, pageNegative.GetOffset()) // Page<=0 会被设为 1
}

// TestPagination_GetLimit 测试 Pagination GetLimit 方法
func TestPagination_GetLimit(t *testing.T) {
	// 测试正常限制
	page1 := &Pagination{PageSize: 10}
	assert.Equal(t, 10, page1.GetLimit())

	page2 := &Pagination{PageSize: 50}
	assert.Equal(t, 50, page2.GetLimit())

	// 测试边界情况
	pageZero := &Pagination{PageSize: 0}
	assert.Equal(t, 20, pageZero.GetLimit()) // PageSize<=0 会被设为 DefaultPageSize=20

	pageNegative := &Pagination{PageSize: -5}
	assert.Equal(t, 20, pageNegative.GetLimit()) // PageSize<=0 会被设为 DefaultPageSize=20
}

// TestPagination_GetTotalPages 测试 Pagination GetTotalPages 方法
func TestPagination_GetTotalPages(t *testing.T) {
	// 测试正常分页
	page1 := &Pagination{Total: 100, PageSize: 10}
	assert.Equal(t, int64(10), page1.GetTotalPages())

	page2 := &Pagination{Total: 95, PageSize: 10}
	assert.Equal(t, int64(10), page2.GetTotalPages()) // 向上取整

	page3 := &Pagination{Total: 91, PageSize: 10}
	assert.Equal(t, int64(10), page3.GetTotalPages())

	page4 := &Pagination{Total: 0, PageSize: 10}
	assert.Equal(t, int64(0), page4.GetTotalPages())

	// 测试边界情况
	pageZero := &Pagination{Total: 100, PageSize: 0}
	assert.Equal(t, int64(0), pageZero.GetTotalPages()) // PageSize<=0 返回 0

	pageNegative := &Pagination{Total: 100, PageSize: -10}
	assert.Equal(t, int64(0), pageNegative.GetTotalPages()) // PageSize<=0 返回 0
}

func TestGetOffset(t *testing.T) {
	p := &Pagination{Page: 1, PageSize: 10}
	assert.Equal(t, 0, p.GetOffset())

	p.Page = 2
	assert.Equal(t, 10, p.GetOffset())

	p.Page = -1
	assert.Equal(t, 0, p.GetOffset()) // 默认 Page 应为 1

	p.Page = 0
	assert.Equal(t, 0, p.GetOffset()) // 默认 Page 应为 1
}

func TestGetLimit(t *testing.T) {
	p := &Pagination{PageSize: 0}
	assert.Equal(t, constants.DefaultPageSize, p.GetLimit())

	p.PageSize = 5
	assert.Equal(t, 5, p.GetLimit())

	// MinPageSize = 1，所以测试 PageSize = 0 应用默认值
	p.PageSize = 0
	assert.Equal(t, constants.DefaultPageSize, p.GetLimit()) // PageSize<=0 先设为 DefaultPageSize

	p.PageSize = constants.MaxPageSize + 1
	assert.Equal(t, constants.MaxPageSize, p.GetLimit()) // 应用最大值限制
}

func TestGetTotalPages(t *testing.T) {
	p := &Pagination{Total: 100, PageSize: 10}
	assert.Equal(t, int64(10), p.GetTotalPages())

	p.PageSize = 0
	assert.Equal(t, int64(0), p.GetTotalPages()) // PageSize <= 0 应返回 0
}

func TestHasNextPage(t *testing.T) {
	p := &Pagination{Total: 100, PageSize: 10, Page: 1}
	assert.True(t, p.HasNextPage())

	p.Page = 10
	assert.False(t, p.HasNextPage())

	p.Total = 0
	assert.False(t, p.HasNextPage()) // 没有记录时
}

func TestHasPrevPage(t *testing.T) {
	p := &Pagination{Page: 2}
	assert.True(t, p.HasPrevPage())

	p.Page = 1
	assert.False(t, p.HasPrevPage())

	p.Page = 0
	assert.False(t, p.HasPrevPage()) // Page <= 1 时
}

func TestIsToday(t *testing.T) {
	today := time.Now()
	assert.True(t, IsToday(&today))

	yesterday := today.AddDate(0, 0, -1)
	assert.False(t, IsToday(&yesterday))

	// 测试 nil 参数
	assert.True(t, IsToday(nil))
}

func TestIsTodayRange(t *testing.T) {
	startTime, endTime := GetTodayRange()

	// 测试今天的范围
	assert.True(t, IsTodayRange(&startTime, &endTime))

	// 测试范围在今天之前
	beforeToday := startTime.AddDate(0, 0, -1)
	assert.False(t, IsTodayRange(&beforeToday, &beforeToday))

	// 测试范围在今天之后
	afterToday := endTime.AddDate(0, 0, 1)
	assert.False(t, IsTodayRange(&afterToday, &afterToday))

	// 测试范围跨越今天
	assert.True(t, IsTodayRange(&beforeToday, &afterToday))

	// 测试仅开始时间在今天
	assert.True(t, IsTodayRange(&startTime, nil))

	// 测试仅结束时间在今天
	assert.True(t, IsTodayRange(nil, &endTime))

	// 测试 nil 参数
	assert.True(t, IsTodayRange(nil, nil))
}

func TestGetTodayRange(t *testing.T) {
	start, end := GetTodayRange()
	now := time.Now()

	// 检查开始时间是否为今天的开始
	assert.Equal(t, start.Year(), now.Year())
	assert.Equal(t, start.Month(), now.Month())
	assert.Equal(t, start.Day(), now.Day())
	assert.Equal(t, start.Hour(), 0)
	assert.Equal(t, start.Minute(), 0)
	assert.Equal(t, start.Second(), 0)

	// 检查结束时间是否为明天的开始
	assert.Equal(t, end.Year(), now.Year())
	assert.Equal(t, end.Month(), now.Month())
	assert.Equal(t, end.Day(), now.Day()+1)
	assert.Equal(t, end.Hour(), 0)
	assert.Equal(t, end.Minute(), 0)
	assert.Equal(t, end.Second(), 0)
}
