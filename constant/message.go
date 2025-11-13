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

// Operation messages
const (
	MsgQueryExecuted         = "query executed successfully"
	MsgTransactionStarted    = "transaction started"
	MsgTransactionCommitted  = "transaction committed"
	MsgTransactionRolledback = "transaction rolled back"
	MsgCacheMiss             = "cache miss for key: %s"
	MsgCacheHit              = "cache hit for key: %s"
)
