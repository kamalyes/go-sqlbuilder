/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-14 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-14 00:00:00
 * @FilePath: \go-sqlbuilder\constant\operation.go
 * @Description: 数据库操作类型常量
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package constant

// 数据库操作类型常量
const (
	// 基本CRUD操作
	OperationTypeCreate = "create"
	OperationTypeRead   = "read"
	OperationTypeUpdate = "update"
	OperationTypeDelete = "delete"
	OperationTypeUpsert = "upsert"
	OperationTypeQuery  = "query"

	// 批量操作
	OperationTypeBatchInsert = "batch_insert"
	OperationTypeBatchUpdate = "batch_update"
	OperationTypeBatchDelete = "batch_delete"
	OperationTypeBatchUpsert = "batch_upsert"

	// 事务操作
	OperationTypeTransaction = "transaction"
	OperationTypeCommit      = "commit"
	OperationTypeRollback    = "rollback"

	// 其他操作
	OperationTypeCount     = "count"
	OperationTypeExists    = "exists"
	OperationTypeAggregate = "aggregate"
)
