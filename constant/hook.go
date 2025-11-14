/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-14 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-14 00:00:00
 * @FilePath: \go-sqlbuilder\constant\hook.go
 * @Description: 钩子事件常量
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package constant

// 钩子事件常量
const (
	// 创建操作钩子
	HookEventBeforeCreate = "before_create"
	HookEventAfterCreate  = "after_create"

	// 更新操作钩子
	HookEventBeforeUpdate = "before_update"
	HookEventAfterUpdate  = "after_update"

	// 删除操作钩子
	HookEventBeforeDelete = "before_delete"
	HookEventAfterDelete  = "after_delete"

	// 查询操作钩子
	HookEventBeforeQuery = "before_query"
	HookEventAfterQuery  = "after_query"

	// 保存操作钩子 (upsert)
	HookEventBeforeSave = "before_save"
	HookEventAfterSave  = "after_save"

	// 批量操作钩子
	HookEventBeforeBatch = "before_batch"
	HookEventAfterBatch  = "after_batch"

	// 事务钩子
	HookEventBeforeTransaction = "before_transaction"
	HookEventAfterTransaction  = "after_transaction"
	HookEventBeforeCommit      = "before_commit"
	HookEventAfterCommit       = "after_commit"
	HookEventBeforeRollback    = "before_rollback"
	HookEventAfterRollback     = "after_rollback"
)
