/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-13 11:03:58
 * @FilePath: \go-sqlbuilder\constant\logger.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package constant

// Logger component names
const (
	ComponentCore       = "core"
	ComponentExecutor   = "executor"
	ComponentCompiler   = "compiler"
	ComponentAdapter    = "adapter"
	ComponentRepository = "repository"
	ComponentCache      = "cache"
	ComponentMiddleware = "middleware"
)

// Logger event names
const (
	EventQueryStart          = "query_start"
	EventQueryComplete       = "query_complete"
	EventQueryError          = "query_error"
	EventTransactionStart    = "transaction_start"
	EventTransactionCommit   = "transaction_commit"
	EventTransactionRollback = "transaction_rollback"
	EventCacheHit            = "cache_hit"
	EventCacheMiss           = "cache_miss"
)

// Logger context keys
const (
	ContextKeySQL          = "sql"
	ContextKeyDuration     = "duration"
	ContextKeyRowsAffected = "rows_affected"
	ContextKeyError        = "error"
	ContextKeyComponent    = "component"
)
