/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-23 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-23 22:50:00
 * @FilePath: \go-sqlbuilder\errors.go
 * @Description: 错误代码定义 - 基于 go-toolbox/errorx 的错误码管理
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package sqlbuilder

import "github.com/kamalyes/go-toolbox/pkg/errorx"

// ErrorCode 错误代码类型(复用 errorx.ErrorType)
type ErrorCode = errorx.ErrorType

// 错误码规则：按模块划分区间，避免重复
// 0-999: 基础通用
// 1000-1999: 通用业务错误
// 2000-2999: 数据库错误
// 3000-3999: 缓存错误
// 4000-4999: SQL构建器错误
// 5000-5999: 权限相关错误
// 9000-9999: 系统级错误

// 基础通用错误 (0-999)
const (
	ErrorCodeSuccess ErrorCode = 0 // 成功
)

// 通用业务错误 (1000-1999)
const (
	ErrorCodeNotFound     ErrorCode = 1001 // 记录不存在
	ErrorCodeAlreadyExist ErrorCode = 1002 // 记录已存在
	ErrorCodeInvalidInput ErrorCode = 1003 // 输入参数无效
	ErrorCodeUnauthorized ErrorCode = 1004 // 未授权
	ErrorCodeForbidden    ErrorCode = 1005 // 禁止访问
	ErrorCodeConflict     ErrorCode = 1006 // 操作冲突
)

// 数据库错误 (2000-2999)
const (
	ErrorCodeDBError           ErrorCode = 2001 // 数据库操作失败
	ErrorCodeDBDuplicate       ErrorCode = 2002 // 数据库记录重复
	ErrorCodeDBConstraint      ErrorCode = 2003 // 数据库约束冲突
	ErrorCodeDBDeadlock        ErrorCode = 2004 // 数据库死锁
	ErrorCodeNoDatabaseConn    ErrorCode = 2005 // 无数据库连接
	ErrorCodeDBFailedUpdate    ErrorCode = 2006 // 数据库更新失败
	ErrorCodeDBFailedInsert    ErrorCode = 2007 // 数据库插入失败
	ErrorCodeDBFailedDelete    ErrorCode = 2008 // 数据库删除失败
	ErrorCodeNestedTransaction ErrorCode = 2009 // 嵌套事务错误
)

// 缓存错误 (3000-3999)
const (
	ErrorCodeCacheError              ErrorCode = 3001 // 缓存操作失败
	ErrorCodeCacheMiss               ErrorCode = 3002 // 缓存未命中
	ErrorCodeCacheStoreNotFound      ErrorCode = 3003 // 缓存存储未找到
	ErrorCodeCacheKeyNotFound        ErrorCode = 3004 // 缓存键不存在
	ErrorCodeCacheExpired            ErrorCode = 3005 // 缓存已过期
	ErrorCodeCacheStoreNotConfigured ErrorCode = 3006 // 缓存存储未配置
	ErrorCodeCacheInvalidData        ErrorCode = 3007 // 缓存数据无效
)

// SQL构建器错误 (4000-4999)
const (
	ErrorCodeResourceNotFound      ErrorCode = 4001 // 资源不存在
	ErrorCodeBuilderNotInitialized ErrorCode = 4002 // Builder未初始化
	ErrorCodeInvalidTableName      ErrorCode = 4003 // 表名无效
	ErrorCodeInvalidFieldName      ErrorCode = 4004 // 字段名无效
	ErrorCodeInvalidSQLQuery       ErrorCode = 4005 // SQL查询无效
	ErrorCodeAdapterNotSupported   ErrorCode = 4006 // 适配器不支持
)

// 查询错误 (复用4000区间扩展)
const (
	ErrorCodeInvalidOperator    ErrorCode = 4101 // 无效的操作符
	ErrorCodeInvalidFilterValue ErrorCode = 4102 // 过滤值无效
	ErrorCodePageNumberInvalid  ErrorCode = 4103 // 页码无效
	ErrorCodePageSizeInvalid    ErrorCode = 4104 // 页大小无效
	ErrorCodeTimeRangeInvalid   ErrorCode = 4105 // 时间范围无效
	ErrorCodeEmptyFilterParam   ErrorCode = 4106 // 过滤参数为空
)

// 系统级错误 (9000-9999)
const (
	ErrorCodeInternal        ErrorCode = 9001 // 内部服务错误
	ErrorCodeTimeout         ErrorCode = 9002 // 请求超时
	ErrorCodeTooManyRequests ErrorCode = 9003 // 请求过于频繁（限流）
	ErrorCodeUnsupported     ErrorCode = 9004 // 不支持的操作
	ErrorCodeUnknown         ErrorCode = 9005 // 未知错误
	ErrorCodeOperationFailed ErrorCode = 9006 // 操作失败
)

// init 注册所有错误码
func init() {
	// 基础通用错误
	errorx.RegisterError(ErrorCodeSuccess, "OK")

	// 通用业务错误
	errorx.RegisterError(ErrorCodeNotFound, "Record not found")
	errorx.RegisterError(ErrorCodeAlreadyExist, "Record already exists")
	errorx.RegisterError(ErrorCodeInvalidInput, "Invalid input parameter")
	errorx.RegisterError(ErrorCodeUnauthorized, "Unauthorized access")
	errorx.RegisterError(ErrorCodeForbidden, "Access forbidden")
	errorx.RegisterError(ErrorCodeConflict, "Operation conflict")

	// 数据库错误
	errorx.RegisterError(ErrorCodeDBError, "Database operation failed")
	errorx.RegisterError(ErrorCodeDBDuplicate, "Database record duplicate")
	errorx.RegisterError(ErrorCodeDBConstraint, "Database constraint violation")
	errorx.RegisterError(ErrorCodeDBDeadlock, "Database deadlock occurred")
	errorx.RegisterError(ErrorCodeNoDatabaseConn, "No database connection available")
	errorx.RegisterError(ErrorCodeDBFailedUpdate, "Database update operation failed")
	errorx.RegisterError(ErrorCodeDBFailedInsert, "Database insert operation failed")
	errorx.RegisterError(ErrorCodeDBFailedDelete, "Database delete operation failed")
	errorx.RegisterError(ErrorCodeNestedTransaction, "Nested transaction not allowed")

	// 缓存错误
	errorx.RegisterError(ErrorCodeCacheError, "Cache operation failed")
	errorx.RegisterError(ErrorCodeCacheMiss, "Cache key not found (miss)")
	errorx.RegisterError(ErrorCodeCacheStoreNotFound, "Cache store not found")
	errorx.RegisterError(ErrorCodeCacheKeyNotFound, "Cache key not found")
	errorx.RegisterError(ErrorCodeCacheExpired, "Cache entry has expired")
	errorx.RegisterError(ErrorCodeCacheStoreNotConfigured, "Cache store not configured")
	errorx.RegisterError(ErrorCodeCacheInvalidData, "Invalid cache data format")

	// SQL构建器错误
	errorx.RegisterError(ErrorCodeResourceNotFound, "Resource not found")
	errorx.RegisterError(ErrorCodeBuilderNotInitialized, "SQL builder not initialized")
	errorx.RegisterError(ErrorCodeInvalidTableName, "Invalid table name")
	errorx.RegisterError(ErrorCodeInvalidFieldName, "Invalid field name")
	errorx.RegisterError(ErrorCodeInvalidSQLQuery, "Invalid SQL query")
	errorx.RegisterError(ErrorCodeAdapterNotSupported, "Database adapter not supported")

	// 查询错误
	errorx.RegisterError(ErrorCodeInvalidOperator, "Invalid query operator")
	errorx.RegisterError(ErrorCodeInvalidFilterValue, "Invalid filter value")
	errorx.RegisterError(ErrorCodePageNumberInvalid, "Invalid page number (must be positive)")
	errorx.RegisterError(ErrorCodePageSizeInvalid, "Invalid page size (must be positive)")
	errorx.RegisterError(ErrorCodeTimeRangeInvalid, "Invalid time range (start > end)")
	errorx.RegisterError(ErrorCodeEmptyFilterParam, "Filter parameter cannot be empty")

	// 系统级错误
	errorx.RegisterError(ErrorCodeInternal, "Internal server error")
	errorx.RegisterError(ErrorCodeTimeout, "Operation timed out")
	errorx.RegisterError(ErrorCodeTooManyRequests, "Too many requests (rate limited)")
	errorx.RegisterError(ErrorCodeUnsupported, "Operation not supported")
	errorx.RegisterError(ErrorCodeUnknown, "Unknown error occurred")
	errorx.RegisterError(ErrorCodeOperationFailed, "Operation failed")
}
