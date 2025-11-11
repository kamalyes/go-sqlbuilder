/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-11 14:32:39
 * @FilePath: \go-sqlbuilder\errors\error.go
 * @Description: 统一的错误定义和管理
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package errors

import "fmt"

// errorMessages 错误消息映射
var errorMessages = map[ErrorCode]string{
	ErrCodeOK:                      "OK",
	ErrCodeBuilderNotInitialized:   "SQL builder not initialized",
	ErrCodeInvalidTableName:        "Invalid table name",
	ErrCodeInvalidFieldName:        "Invalid field name",
	ErrCodeInvalidSQLQuery:         "Invalid SQL query",
	ErrCodeAdapterNotSupported:     "Adapter not supported",
	ErrCodeCacheStoreNotFound:      "Cache store not found",
	ErrCodeCacheKeyNotFound:        "Cache key not found",
	ErrCodeCacheExpired:            "Cache entry expired",
	ErrCodeCacheStoreNotConfigured: "Cache store not configured",
	ErrCodeCacheInvalidData:        "Invalid cache data format",
	ErrCodeInvalidOperator:         "Invalid query operator",
	ErrCodeInvalidFilterValue:      "Invalid filter value",
	ErrCodePageNumberInvalid:       "Invalid page number",
	ErrCodePageSizeInvalid:         "Invalid page size",
	ErrCodeTimeRangeInvalid:        "Invalid time range",
	ErrCodeEmptyFilterParam:        "Empty filter parameter",
	ErrCodeRedisConnFailed:         "Redis connection failed",
	ErrCodeRedisOperationFailed:    "Redis operation failed",
	ErrCodeRedisKeyNotFound:        "Redis key not found",
	ErrCodeRedisAdapterNotImpl:     "Redis adapter not implemented",
	ErrCodeUnknown:                 "Unknown error",
	ErrCodeInternal:                "Internal error",
	ErrCodeInvalidParam:            "Invalid parameter",
	ErrCodeOperationFailed:         "Operation failed",
}

// AppError 应用错误结构
type AppError struct {
	Code    ErrorCode // 错误代码
	Message string    // 错误消息
	Details string    // 错误详情
}

// NewError 创建新错误
func NewError(code ErrorCode, message string) *AppError {
	if msg, ok := errorMessages[code]; ok {
		return &AppError{
			Code:    code,
			Message: msg,
			Details: message,
		}
	}
	return &AppError{
		Code:    code,
		Message: errorMessages[ErrCodeUnknown],
		Details: message,
	}
}

// NewErrorf 使用格式化字符串创建错误
func NewErrorf(code ErrorCode, format string, args ...interface{}) *AppError {
	return NewError(code, fmt.Sprintf(format, args...))
}

// Error 实现 error 接口
func (e *AppError) Error() string {
	if e.Details == "" {
		return fmt.Sprintf("[%d] %s", e.Code, e.Message)
	}
	return fmt.Sprintf("[%d] %s: %s", e.Code, e.Message, e.Details)
}

// String 实现 Stringer 接口，便于直接打印
func (e *AppError) String() string {
	return e.Error()
}

// GetCode 获取错误代码
func (e *AppError) GetCode() ErrorCode {
	return e.Code
}

// GetMessage 获取错误消息
func (e *AppError) GetMessage() string {
	return e.Message
}

// GetDetails 获取错误详情
func (e *AppError) GetDetails() string {
	return e.Details
}

// WithDetails 添加错误详情
func (e *AppError) WithDetails(details string) *AppError {
	e.Details = details
	return e
}

// IsErrorCode 检查错误代码是否匹配
func IsErrorCode(err error, code ErrorCode) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == code
	}
	return false
}

// GetErrorCode 从错误中提取错误代码
func GetErrorCode(err error) ErrorCode {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code
	}
	return ErrCodeUnknown
}

// ErrorCodeString 获取错误代码的字符串表示
func ErrorCodeString(code ErrorCode) string {
	if msg, ok := errorMessages[code]; ok {
		return msg
	}
	return errorMessages[ErrCodeUnknown]
}

// Predefined Errors 预定义错误变量

// CacheErrors
var (
	ErrCacheNotConfigured = NewError(ErrCodeCacheStoreNotConfigured, "")
	ErrCacheKeyMissing    = NewError(ErrCodeCacheKeyNotFound, "")
	ErrCacheDataInvalid   = NewError(ErrCodeCacheInvalidData, "")
)

// QueryErrors
var (
	ErrInvalidOp        = NewError(ErrCodeInvalidOperator, "")
	ErrEmptyFilter      = NewError(ErrCodeEmptyFilterParam, "")
	ErrInvalidPage      = NewError(ErrCodePageNumberInvalid, "")
	ErrInvalidPageNum   = NewError(ErrCodePageNumberInvalid, "")
	ErrInvalidPageSize  = NewError(ErrCodePageSizeInvalid, "")
	ErrInvalidTimeRange = NewError(ErrCodeTimeRangeInvalid, "")
)

// BuilderErrors
var (
	ErrBuilderNotInit = NewError(ErrCodeBuilderNotInitialized, "")
	ErrInvalidTable   = NewError(ErrCodeInvalidTableName, "")
	ErrInvalidField   = NewError(ErrCodeInvalidFieldName, "")
)

// RedisErrors
var (
	ErrRedisNotImpl = NewError(ErrCodeRedisAdapterNotImpl, "Redis adapter not implemented")
	ErrRedisFailed  = NewError(ErrCodeRedisOperationFailed, "")
)
