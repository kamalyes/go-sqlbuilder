/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-14 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-14 00:00:00
 * @FilePath: \go-sqlbuilder\constant\field.go
 * @Description: 数据库字段名常量
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package constant

// 审计字段名常量
const (
	// 时间戳字段
	FieldCreatedAt = "created_at"
	FieldUpdatedAt = "updated_at"
	FieldDeletedAt = "deleted_at"

	// 用户信息字段
	FieldCreatedBy = "created_by"
	FieldUpdatedBy = "updated_by"
	FieldDeletedBy = "deleted_by"

	// 版本控制字段
	FieldVersion  = "version"
	FieldRevision = "revision"

	// 通用字段
	FieldID          = "id"
	FieldName        = "name"
	FieldStatus      = "status"
	FieldIsDeleted   = "is_deleted"
	FieldIsActive    = "is_active"
	FieldRemark      = "remark"
	FieldDescription = "description"
)

// 默认字段值常量
const (
	DefaultDeletedValue    = 1
	DefaultNotDeletedValue = 0
	DefaultActiveValue     = 1
	DefaultInactiveValue   = 0
)
