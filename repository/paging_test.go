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
	"testing"
	"time"

	"github.com/kamalyes/go-sqlbuilder/constants"
	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, 10, page1.GetTotalPages())

	page2 := &Pagination{Total: 95, PageSize: 10}
	assert.Equal(t, 10, page2.GetTotalPages()) // 向上取整

	page3 := &Pagination{Total: 91, PageSize: 10}
	assert.Equal(t, 10, page3.GetTotalPages())

	page4 := &Pagination{Total: 0, PageSize: 10}
	assert.Equal(t, 0, page4.GetTotalPages())

	// 测试边界情况
	pageZero := &Pagination{Total: 100, PageSize: 0}
	assert.Equal(t, 0, pageZero.GetTotalPages()) // PageSize<=0 返回 0

	pageNegative := &Pagination{Total: 100, PageSize: -10}
	assert.Equal(t, 0, pageNegative.GetTotalPages()) // PageSize<=0 返回 0
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
	assert.Equal(t, 10, p.GetTotalPages())

	p.PageSize = 0
	assert.Equal(t, 0, p.GetTotalPages()) // PageSize <= 0 应返回 0
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

// TestPaginationT_ToInt 测试转换为 int 类型分页
func TestPaginationT_ToInt(t *testing.T) {
	// 测试 int32 转 int
	p32 := &Pagination32{Page: 5, PageSize: 20, Total: 100}
	p := p32.ToInt()
	assert.Equal(t, 5, p.Page)
	assert.Equal(t, 20, p.PageSize)
	assert.Equal(t, int64(100), p.Total)

	// 测试 int64 转 int
	p64 := &Pagination64{Page: 10, PageSize: 50, Total: 500}
	p = p64.ToInt()
	assert.Equal(t, 10, p.Page)
	assert.Equal(t, 50, p.PageSize)
	assert.Equal(t, int64(500), p.Total)

	// 测试 int 转 int（自身转换）
	pInt := &Pagination{Page: 3, PageSize: 15, Total: 75}
	p = pInt.ToInt()
	assert.Equal(t, 3, p.Page)
	assert.Equal(t, 15, p.PageSize)
	assert.Equal(t, int64(75), p.Total)
}

// TestPaginationT_ToInt32 测试转换为 int32 类型分页
func TestPaginationT_ToInt32(t *testing.T) {
	// 测试 int 转 int32
	p := &Pagination{Page: 7, PageSize: 25, Total: 200}
	p32 := p.ToInt32()
	assert.Equal(t, int32(7), p32.Page)
	assert.Equal(t, int32(25), p32.PageSize)
	assert.Equal(t, int64(200), p32.Total)

	// 测试 int64 转 int32
	p64 := &Pagination64{Page: 8, PageSize: 30, Total: 300}
	p32 = p64.ToInt32()
	assert.Equal(t, int32(8), p32.Page)
	assert.Equal(t, int32(30), p32.PageSize)
	assert.Equal(t, int64(300), p32.Total)

	// 测试 int32 转 int32（自身转换）
	p32Orig := &Pagination32{Page: 2, PageSize: 10, Total: 50}
	p32 = p32Orig.ToInt32()
	assert.Equal(t, int32(2), p32.Page)
	assert.Equal(t, int32(10), p32.PageSize)
	assert.Equal(t, int64(50), p32.Total)
}

// TestPaginationT_ToInt64 测试转换为 int64 类型分页
func TestPaginationT_ToInt64(t *testing.T) {
	// 测试 int 转 int64
	p := &Pagination{Page: 12, PageSize: 40, Total: 600}
	p64 := p.ToInt64()
	assert.Equal(t, int64(12), p64.Page)
	assert.Equal(t, int64(40), p64.PageSize)
	assert.Equal(t, int64(600), p64.Total)

	// 测试 int32 转 int64
	p32 := &Pagination32{Page: 15, PageSize: 45, Total: 700}
	p64 = p32.ToInt64()
	assert.Equal(t, int64(15), p64.Page)
	assert.Equal(t, int64(45), p64.PageSize)
	assert.Equal(t, int64(700), p64.Total)

	// 测试 int64 转 int64（自身转换）
	p64Orig := &Pagination64{Page: 20, PageSize: 60, Total: 1000}
	p64 = p64Orig.ToInt64()
	assert.Equal(t, int64(20), p64.Page)
	assert.Equal(t, int64(60), p64.PageSize)
	assert.Equal(t, int64(1000), p64.Total)
}

// TestPaginationT_Conversion_Chain 测试链式转换
func TestPaginationT_Conversion_Chain(t *testing.T) {
	// 测试 int -> int32 -> int64 -> int
	p := &Pagination{Page: 5, PageSize: 20, Total: 100}

	p32 := p.ToInt32()
	assert.Equal(t, int32(5), p32.Page)

	p64 := p32.ToInt64()
	assert.Equal(t, int64(5), p64.Page)

	pFinal := p64.ToInt()
	assert.Equal(t, 5, pFinal.Page)
	assert.Equal(t, 20, pFinal.PageSize)
	assert.Equal(t, int64(100), pFinal.Total)
}

// TestPaginationT_Conversion_WithZeroValues 测试零值转换
func TestPaginationT_Conversion_WithZeroValues(t *testing.T) {
	// 测试零值转换
	p := &Pagination{Page: 0, PageSize: 0, Total: 0}

	p32 := p.ToInt32()
	assert.Equal(t, int32(0), p32.Page)
	assert.Equal(t, int32(0), p32.PageSize)
	assert.Equal(t, int64(0), p32.Total)

	p64 := p.ToInt64()
	assert.Equal(t, int64(0), p64.Page)
	assert.Equal(t, int64(0), p64.PageSize)
	assert.Equal(t, int64(0), p64.Total)
}

// TestPaginationT_Conversion_PreserveMethods 测试转换后方法仍然有效
func TestPaginationT_Conversion_PreserveMethods(t *testing.T) {
	// 创建一个 int32 分页
	p32 := &Pagination32{Page: 2, PageSize: 10, Total: 100}

	// 转换为 int64
	p64 := p32.ToInt64()

	// 验证方法调用正常
	assert.Equal(t, int(10), p64.GetOffset())
	assert.Equal(t, 10, p64.GetLimit())
	assert.Equal(t, int64(10), p64.GetTotalPages())
	assert.True(t, p64.HasNextPage())
	assert.True(t, p64.HasPrevPage())

	// 转换为 int
	p := p64.ToInt()

	// 再次验证方法调用
	assert.Equal(t, int(10), p.GetOffset())
	assert.Equal(t, 10, p.GetLimit())
	assert.Equal(t, 10, p.GetTotalPages())
	assert.True(t, p.HasNextPage())
	assert.True(t, p.HasPrevPage())
}
