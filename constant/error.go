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
package constant

// Error message constants
const (
	// Operator errors
	ErrUnknownOperator = "unknown operator: %s"

	// Filter errors
	ErrInvalidFilterField = "invalid filter field: %s"
	ErrInvalidFilterValue = "invalid filter value for operator %s"

	// Query errors
	ErrEmptyTable        = "table name cannot be empty"
	ErrInvalidPagination = "invalid pagination parameters"

	// Database errors
	ErrNoDatabaseConnection = "no database connection available"
	ErrTransactionFailed    = "transaction failed: %w"

	// Cache errors
	ErrCacheKeyNotFound = "cache key not found"
	ErrCacheFailed      = "cache operation failed: %w"
)
