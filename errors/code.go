/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-13 23:52:51
 * @FilePath: \go-sqlbuilder\errors\code.go
 * @Description: 统一的错误定义和管理
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package errors

// ErrorCode 定义错误代码类型
type ErrorCode int

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
	ErrorCodeDBError             ErrorCode = 2001 // 数据库操作失败
	ErrorCodeDBDuplicate         ErrorCode = 2002 // 数据库记录重复
	ErrorCodeDBConstraint        ErrorCode = 2003 // 数据库约束冲突
	ErrorCodeDBDeadlock          ErrorCode = 2004 // 数据库死锁
	ErrorCodeNoDatabaseConn      ErrorCode = 2005 // 无数据库连接
	ErrorCodeDBFailedUpdate      ErrorCode = 2006 // 数据库更新失败
	ErrorCodeDBFailedInsert      ErrorCode = 2007 // 数据库插入失败
	ErrorCodeDBFailedDelete      ErrorCode = 2008 // 数据库删除失败
	ErrorCodeNestedTransaction   ErrorCode = 2009 // 嵌套事务错误
)

// 缓存错误 (3000-3999)
const (
	ErrorCodeCacheError            ErrorCode = 3001 // 缓存操作失败
	ErrorCodeCacheMiss             ErrorCode = 3002 // 缓存未命中
	ErrorCodeCacheStoreNotFound    ErrorCode = 3003 // 缓存存储未找到
	ErrorCodeCacheKeyNotFound      ErrorCode = 3004 // 缓存键不存在
	ErrorCodeCacheExpired          ErrorCode = 3005 // 缓存已过期
	ErrorCodeCacheStoreNotConfigured ErrorCode = 3006 // 缓存存储未配置
	ErrorCodeCacheInvalidData      ErrorCode = 3007 // 缓存数据无效
)

// SQL构建器错误 (4000-4999)
const (
	ErrorCodeResourceNotFound    ErrorCode = 4001 // 资源不存在
	ErrorCodeBuilderNotInitialized      ErrorCode = 4002 // Builder未初始化
	ErrorCodeInvalidTableName    ErrorCode = 4003 // 表名无效
	ErrorCodeInvalidFieldName    ErrorCode = 4004 // 字段名无效
	ErrorCodeInvalidSQLQuery     ErrorCode = 4005 // SQL查询无效
	ErrorCodeAdapterNotSupported ErrorCode = 4006 // 适配器不支持
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

// Redis错误 (复用3000区间扩展)
const (
	ErrorCodeRedisConnFailed      ErrorCode = 3101 // Redis连接失败
	ErrorCodeRedisOperationFailed ErrorCode = 3102 // Redis操作失败
	ErrorCodeRedisKeyNotFound     ErrorCode = 3103 // Redis键不存在
	ErrorCodeRedisAdapterNotImpl  ErrorCode = 3104 // Redis适配器未实现
)

// 权限和访问控制错误 (5000-5999)
const (
	ErrorCodeAccessDenied             ErrorCode = 5001 // 访问被拒绝
	ErrorCodeUserMismatch             ErrorCode = 5002 // 用户不匹配
	ErrorCodeResourceNotOwnedByUser   ErrorCode = 5003 // 资源不属于该用户
)

// 系统级错误 (9000-9999)
const (
	ErrorCodeInternal         ErrorCode = 9001 // 内部服务错误
	ErrorCodeTimeout          ErrorCode = 9002 // 请求超时
	ErrorCodeTooManyRequests  ErrorCode = 9003 // 请求过于频繁（限流）
	ErrorCodeUnsupported      ErrorCode = 9004 // 不支持的操作
	ErrorCodeUnknown          ErrorCode = 9005 // 未知错误
	ErrorCodeOperationFailed  ErrorCode = 9006 // 操作失败
)

// ==================== 错误消息常量 ====================

const (
	// 数据库相关消息
	MsgNoDatabaseConnection       = "no database connection available"
	MsgCannotBeginNestedTransaction = "cannot begin transaction within transaction"
	MsgDatabaseOperationFailed    = "database operation failed"
	MsgFailedToExecuteUpdate      = "failed to execute update"
	MsgFailedToExecuteDelete      = "failed to execute delete"
	MsgFailedToExecuteInsert      = "failed to execute insert"
	MsgQueryTimeout               = "query timeout"

	// 缓存相关消息
	MsgKeyNotFound                = "key not found"
	MsgCacheOperationFailed       = "cache operation failed"
	MsgCacheStoreNotInitialized   = "cache store not initialized"
	MsgFailedToGetCache           = "failed to get cache"
	MsgFailedToSetCache           = "failed to set cache"
	MsgFailedToDeleteCache        = "failed to delete cache"

	// Builder相关消息
	MsgBuilderNotInitialized      = "builder not initialized"
	MsgInvalidTableName           = "invalid table name"
	MsgNoColumnsSelected          = "no columns selected"
	MsgInvalidWhereCondition      = "invalid where condition"

	// 资源相关消息
	MsgResourceNotFound           = "%s not found"
	MsgResourceAlreadyExists      = "%s already exists"
	MsgResourceConflict           = "%s conflict"

	// 权限相关消息
	MsgAccessDenied               = "access denied"
	MsgUserMismatch               = "user mismatch"
	MsgResourceNotOwnedByUser     = "%s is not owned by the current user"

	// 系统相关消息
	MsgInternalError              = "internal server error"
	MsgNotSupported               = "operation not supported"
	MsgInvalidArgument            = "invalid argument"

	// 类型相关消息
	MsgInvalidType                = "invalid type: %s"
	MsgTypeConversionFailed       = "type conversion failed: %s"
	MsgUnexpectedType             = "unexpected type %s, expected %s"

	// 适配器相关消息
	MsgAdapterNotSupported        = "adapter not supported"
	MsgGormNotSupportPrepare      = "gorm does not support prepared statements directly"
	MsgGormNotSupportLastInsertId = "gorm does not support LastInsertId, use returning clause"
	MsgPostgresNotSupportLastInsertId = "postgres does not support LastInsertId, use RETURNING clause"
	MsgUnknownAdapter             = "unknown adapter: %s"
	MsgUnsupportedDatabaseInstance = "unsupported database instance type"

	// 通用消息
	MsgEntityCannotBeNil          = "entity cannot be nil"
	MsgFilterCannotBeNil          = "filter cannot be nil"
	MsgAtLeastOneFilterRequired   = "at least one filter is required"
	MsgCacheValueCannotBeNil      = "cache value cannot be nil"
	MsgFailedToMarshalCache       = "failed to marshal cache value"
	MsgDeletePatternNotImplemented = "DeletePattern not implemented for current cache backend"
	MsgTTLQueryNotImplemented     = "TTL query not implemented for current cache backend"
	MsgBatchSetEncounteredErrors  = "batch set encountered errors"
	MsgUnsupportedDBType          = "unsupported db type: %s"
	MsgDestMustBePointer          = "dest must be a pointer"
)