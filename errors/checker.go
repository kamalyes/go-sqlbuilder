/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-13 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-13 10:26:53
 * @FilePath: \go-sqlbuilder\errors\checker.go
 * @Description: 数据库错误检测工具函数
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package errors

import (
	"errors"

	"github.com/go-sql-driver/mysql"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// IsDuplicateKeyError 检查是否为重复键错误（唯一索引冲突）
// 支持：
// - GORM ErrDuplicatedKey (GORM 1.20+)
// - MySQL Error 1062 (ER_DUP_ENTRY)
// - PostgreSQL Error 23505 (unique_violation)
func IsDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}

	// 方法1: 检查 GORM 的 ErrDuplicatedKey
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}

	// 方法2: 检查 MySQL 驱动的错误码 1062 (ER_DUP_ENTRY)
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return true
	}

	// 方法3: 检查 PostgreSQL 驱动的错误码 23505 (unique_violation)
	var pqErr *pq.Error
	if errors.As(err, &pqErr) && pqErr.Code == "23505" {
		return true
	}

	return false
}

// IsRecordNotFoundError 检查是否为记录不存在错误
func IsRecordNotFoundError(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

// IsForeignKeyViolation 检查是否为外键约束错误
func IsForeignKeyViolation(err error) bool {
	if err == nil {
		return false
	}

	// MySQL Error 1451: Cannot delete or update a parent row
	// MySQL Error 1452: Cannot add or update a child row
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1451 || mysqlErr.Number == 1452
	}

	// PostgreSQL Error 23503: foreign_key_violation
	var pqErr *pq.Error
	if errors.As(err, &pqErr) && pqErr.Code == "23503" {
		return true
	}

	return false
}
