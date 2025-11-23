/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-23 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-23 23:56:05
 * @FilePath: \go-sqlbuilder\paging_test.go
 * @Description: 分页工具 - Pagination分页元数据和辅助方法
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package sqlbuilder

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
	assert.Equal(t, 10, pageZero.GetLimit()) // PageSize<=0 会被设为 10

	pageNegative := &Pagination{PageSize: -5}
	assert.Equal(t, 10, pageNegative.GetLimit()) // PageSize<=0 会被设为 10
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
