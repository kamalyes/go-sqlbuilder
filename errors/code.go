/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-11 15:31:12
 * @FilePath: \go-sqlbuilder\errors\code.go
 * @Description: 统一的错误定义和管理
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package errors

// ErrorCode 定义错误代码类型
type ErrorCode int

// 定义错误代码常量
const (
	ErrCodeOK ErrorCode = iota

	// 构建器错误 (1000-1999)
	ErrCodeBuilderNotInitialized ErrorCode = 1001
	ErrCodeInvalidTableName      ErrorCode = 1002
	ErrCodeInvalidFieldName      ErrorCode = 1003
	ErrCodeInvalidSQLQuery       ErrorCode = 1004
	ErrCodeAdapterNotSupported   ErrorCode = 1005

	// 缓存错误 (2000-2999)
	ErrCodeCacheStoreNotFound      ErrorCode = 2001
	ErrCodeCacheKeyNotFound        ErrorCode = 2002
	ErrCodeCacheExpired            ErrorCode = 2003
	ErrCodeCacheStoreNotConfigured ErrorCode = 2004
	ErrCodeCacheInvalidData        ErrorCode = 2005

	// 查询错误 (3000-3999)
	ErrCodeInvalidOperator    ErrorCode = 3001
	ErrCodeInvalidFilterValue ErrorCode = 3002
	ErrCodePageNumberInvalid  ErrorCode = 3003
	ErrCodePageSizeInvalid    ErrorCode = 3004
	ErrCodeTimeRangeInvalid   ErrorCode = 3005
	ErrCodeEmptyFilterParam   ErrorCode = 3006

	// Redis 错误 (4000-4999)
	ErrCodeRedisConnFailed      ErrorCode = 4001
	ErrCodeRedisOperationFailed ErrorCode = 4002
	ErrCodeRedisKeyNotFound     ErrorCode = 4003
	ErrCodeRedisAdapterNotImpl  ErrorCode = 4004

	// 通用错误 (5000-5999)
	ErrCodeUnknown         ErrorCode = 5000
	ErrCodeInternal        ErrorCode = 5001
	ErrCodeInvalidParam    ErrorCode = 5002
	ErrCodeOperationFailed ErrorCode = 5003
)
