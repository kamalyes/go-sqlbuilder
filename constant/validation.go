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

// Filter error messages
const (
	ErrFilterFieldEmpty    = "filter field cannot be empty"
	ErrFilterValueNil      = "filter value cannot be nil"
	ErrFilterOperatorEmpty = "filter operator cannot be empty"
)

// Order error messages
const (
	ErrOrderFieldEmpty   = "order field cannot be empty"
	ErrOrderValueInvalid = "invalid order value"
)

// Pagination error messages
const (
	ErrPaginationLimitInvalid  = "pagination limit must be positive"
	ErrPaginationOffsetInvalid = "pagination offset must be non-negative"
	ErrPaginationPageInvalid   = "invalid page number (must be positive)"
	ErrPaginationSizeInvalid   = "invalid page size (must be positive)"
)

// Query error messages
const (
	ErrQueryNil          = "query cannot be nil"
	ErrQueryFieldEmpty   = "query field cannot be empty"
	ErrQueryConditionNil = "query condition cannot be nil"
)

// Operator error messages
const (
	ErrOperatorEmpty        = "operator cannot be empty"
	ErrOperatorNotSupported = "operator not supported: %s"
)

// Value validation error messages
const (
	ErrValueNil          = "%s value cannot be nil"
	ErrValueEmpty        = "%s value cannot be empty"
	ErrValueInvalid      = "%s value is invalid"
	ErrValueTypeMismatch = "value type mismatch: expected %s, got %s"
)

// General error messages
const (
	ErrInputInvalid  = "invalid input parameter"
	ErrParamRequired = "parameter %s is required"
	ErrParamInvalid  = "parameter %s is invalid"
)
