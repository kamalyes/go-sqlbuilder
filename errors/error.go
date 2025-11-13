/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-13 00:02:15
 * @FilePath: \go-sqlbuilder\errors\error.go
 * @Description: 统一的错误定义和管理
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package errors

import "fmt"

// errorMessages 错误消息映射
var errorMessages = map[ErrorCode]string{
	// 基础通用错误 (0-999)
	ErrorCodeSuccess: "OK", // 操作成功

	// 通用业务错误 (1000-1999)
	ErrorCodeNotFound:     "Record not found",        // 记录不存在
	ErrorCodeAlreadyExist: "Record already exists",   // 记录已存在
	ErrorCodeInvalidInput: "Invalid input parameter", // 输入参数无效
	ErrorCodeUnauthorized: "Unauthorized access",     // 未授权
	ErrorCodeForbidden:    "Access forbidden",        // 禁止访问
	ErrorCodeConflict:     "Operation conflict",      // 操作冲突

	// 数据库错误 (2000-2999)
	ErrorCodeDBError:           "Database operation failed",        // 数据库操作失败
	ErrorCodeDBDuplicate:       "Database record duplicate",        // 数据库记录重复
	ErrorCodeDBConstraint:      "Database constraint violation",    // 数据库约束冲突
	ErrorCodeDBDeadlock:        "Database deadlock occurred",       // 数据库死锁
	ErrorCodeNoDatabaseConn:    "No database connection available", // 无数据库连接
	ErrorCodeDBFailedUpdate:    "Database update operation failed", // 数据库更新失败
	ErrorCodeDBFailedInsert:    "Database insert operation failed", // 数据库插入失败
	ErrorCodeDBFailedDelete:    "Database delete operation failed", // 数据库删除失败
	ErrorCodeNestedTransaction: "Nested transaction not allowed",   // 嵌套事务错误

	// 缓存错误 (3000-3999)
	ErrorCodeCacheError:              "Cache operation failed",     // 缓存操作失败
	ErrorCodeCacheMiss:               "Cache key not found (miss)", // 缓存未命中
	ErrorCodeCacheStoreNotFound:      "Cache store not found",      // 缓存存储未找到
	ErrorCodeCacheKeyNotFound:        "Cache key not found",        // 缓存键不存在
	ErrorCodeCacheExpired:            "Cache entry has expired",    // 缓存已过期
	ErrorCodeCacheStoreNotConfigured: "Cache store not configured", // 缓存存储未配置
	ErrorCodeCacheInvalidData:        "Invalid cache data format",  // 缓存数据无效

	// SQL构建器错误 (4000-4999)
	ErrorCodeResourceNotFound:      "Resource not found",             // 资源不存在
	ErrorCodeBuilderNotInitialized: "SQL builder not initialized",    // Builder未初始化
	ErrorCodeInvalidTableName:      "Invalid table name",             // 表名无效
	ErrorCodeInvalidFieldName:      "Invalid field name",             // 字段名无效
	ErrorCodeInvalidSQLQuery:       "Invalid SQL query",              // SQL查询无效
	ErrorCodeAdapterNotSupported:   "Database adapter not supported", // 适配器不支持

	// 查询错误 (4100-4199)
	ErrorCodeInvalidOperator:    "Invalid query operator",                 // 无效的操作符
	ErrorCodeInvalidFilterValue: "Invalid filter value",                   // 过滤值无效
	ErrorCodePageNumberInvalid:  "Invalid page number (must be positive)", // 页码无效
	ErrorCodePageSizeInvalid:    "Invalid page size (must be positive)",   // 页大小无效
	ErrorCodeTimeRangeInvalid:   "Invalid time range (start > end)",       // 时间范围无效
	ErrorCodeEmptyFilterParam:   "Filter parameter cannot be empty",       // 过滤参数为空

	// Redis错误 (3100-3199)
	ErrorCodeRedisConnFailed:      "Redis connection failed",       // Redis连接失败
	ErrorCodeRedisOperationFailed: "Redis operation failed",        // Redis操作失败
	ErrorCodeRedisKeyNotFound:     "Redis key not found",           // Redis键不存在
	ErrorCodeRedisAdapterNotImpl:  "Redis adapter not implemented", // Redis适配器未实现

	// 权限和访问控制错误 (5000-5999)
	ErrorCodeAccessDenied:           "Access denied",              // 访问被拒绝
	ErrorCodeUserMismatch:           "User information mismatch",  // 用户不匹配
	ErrorCodeResourceNotOwnedByUser: "Resource not owned by user", // 资源不属于该用户

	// 系统级错误 (9000-9999)
	ErrorCodeInternal:        "Internal server error",            // 内部服务错误
	ErrorCodeTimeout:         "Operation timed out",              // 请求超时
	ErrorCodeTooManyRequests: "Too many requests (rate limited)", // 请求过于频繁（限流）
	ErrorCodeUnsupported:     "Operation not supported",          // 不支持的操作
	ErrorCodeUnknown:         "Unknown error occurred",           // 未知错误
	ErrorCodeOperationFailed: "Operation failed",                 // 操作失败
}

// AppError 应用错误结构
type AppError struct {
	Code    ErrorCode // 错误代码
	Message string    // 错误消息
	Details string    // 错误详情
}

// BusinessError 业务错误（别名，保持兼容）
type BusinessError = AppError

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
		Message: errorMessages[ErrorCodeUnknown],
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
	return ErrorCodeUnknown
}

// ErrorCodeString 获取错误代码的字符串表示
func ErrorCodeString(code ErrorCode) string {
	if msg, ok := errorMessages[code]; ok {
		return msg
	}
	return errorMessages[ErrorCodeUnknown]
}

// ==================== 业务错误构造函数 ====================

// NewNotFound 创建记录不存在错误
func NewNotFound(resource string) *AppError {
	return NewError(ErrorCodeNotFound, fmt.Sprintf("%s not found", resource))
}

// NewAlreadyExists 创建记录已存在错误
func NewAlreadyExists(resource string) *AppError {
	return NewError(ErrorCodeAlreadyExist, fmt.Sprintf("%s already exists", resource))
}

// NewInvalidInput 创建参数无效错误
func NewInvalidInput(field, reason string) *AppError {
	return NewError(ErrorCodeInvalidInput, fmt.Sprintf("invalid %s: %s", field, reason))
}

// NewAccessDenied 创建访问被拒绝错误
func NewAccessDenied(reason string) *AppError {
	return NewError(ErrorCodeAccessDenied, fmt.Sprintf("access denied: %s", reason))
}

// NewConflict 创建冲突错误
func NewConflict(resource string) *AppError {
	return NewError(ErrorCodeConflict, fmt.Sprintf("%s conflict", resource))
}

// NewDatabaseError 创建数据库错误
func NewDatabaseError(operation string, err error) *AppError {
	return NewError(ErrorCodeDBError, fmt.Sprintf("database %s failed: %v", operation, err))
}

// NewInternal 创建内部错误
func NewInternal(message string) *AppError {
	return NewError(ErrorCodeInternal, "internal server error").WithDetails(message)
}

// ==================== 错误检查函数 ====================

// IsNotFound 检查是否是记录不存在错误
func IsNotFound(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == ErrorCodeNotFound
	}
	return false
}

// IsAccessDenied 检查是否是权限错误
func IsAccessDenied(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == ErrorCodeAccessDenied || appErr.Code == ErrorCodeUserMismatch ||
			appErr.Code == ErrorCodeResourceNotOwnedByUser
	}
	return false
}

// IsConflict 检查是否是冲突错误
func IsConflict(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == ErrorCodeConflict
	}
	return false
}

// IsAlreadyExists 检查是否是已存在错误
func IsAlreadyExists(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == ErrorCodeAlreadyExist
	}
	return false
}

// Predefined Errors 预定义错误变量

// CacheErrors
var (
	ErrCacheNotConfigured = NewError(ErrorCodeCacheStoreNotConfigured, "")
	ErrCacheKeyMissing    = NewError(ErrorCodeCacheKeyNotFound, "")
	ErrCacheDataInvalid   = NewError(ErrorCodeCacheInvalidData, "")
)

// QueryErrors
var (
	ErrInvalidOp        = NewError(ErrorCodeInvalidOperator, "")
	ErrEmptyFilter      = NewError(ErrorCodeEmptyFilterParam, "")
	ErrInvalidPage      = NewError(ErrorCodePageNumberInvalid, "")
	ErrInvalidPageNum   = NewError(ErrorCodePageNumberInvalid, "")
	ErrInvalidPageSize  = NewError(ErrorCodePageSizeInvalid, "")
	ErrInvalidTimeRange = NewError(ErrorCodeTimeRangeInvalid, "")
)

// BuilderErrors
var (
	ErrBuilderNotInit = NewError(ErrorCodeBuilderNotInitialized, "")
	ErrInvalidTable   = NewError(ErrorCodeInvalidTableName, "")
	ErrInvalidField   = NewError(ErrorCodeInvalidFieldName, "")
)

// Wrap 包装已有的错误并添加错误代码
func Wrap(err error, code ErrorCode) *AppError {
	if err == nil {
		return nil
	}
	return NewError(code, err.Error())
}
