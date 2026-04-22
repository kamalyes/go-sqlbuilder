/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-30 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-30 00:00:00
 * @FilePath: \go-sqlbuilder\constants\database.go
 * @Description: 数据库类型常量定义
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package constants

// ==================== 数据库类型常量 ====================

const (
	// DatabaseTypeMySQL MySQL 数据库类型
	DatabaseTypeMySQL = "mysql"

	// DatabaseTypePostgreSQL PostgreSQL 数据库类型
	DatabaseTypePostgreSQL = "postgres"

	// DatabaseTypeCockroachDB CockroachDB 数据库类型
	DatabaseTypeCockroachDB = "cockroachdb"
)

// ==================== 数据库方言常量 ====================

const (
	// DialectorMySQL MySQL 方言
	DialectorMySQL = DatabaseTypeMySQL

	// DialectorPostgreSQL PostgreSQL 方言
	DialectorPostgreSQL = DatabaseTypePostgreSQL

	// DialectorCockroachDB CockroachDB 方言
	DialectorCockroachDB = DatabaseTypeCockroachDB
)

// ==================== 数据库方言组 ====================

// PostgreSQLDialectors PostgreSQL 方言组（PostgreSQL 和 CockroachDB）
var PostgreSQLDialectors = []string{
	DialectorPostgreSQL,
	DialectorCockroachDB,
}

// SupportedDialectors 支持的数据库方言列表
var SupportedDialectors = []string{
	DialectorMySQL,
	DialectorPostgreSQL,
	DialectorCockroachDB,
}

// ==================== 数据库方言检查函数 ====================

// IsMySQLDialector 检查是否为 MySQL 方言
func IsMySQLDialector(dialector string) bool {
	return dialector == DialectorMySQL
}

// IsPostgreSQLDialector 检查是否为 PostgreSQL 方言
func IsPostgreSQLDialector(dialector string) bool {
	return dialector == DialectorPostgreSQL
}

// IsCockroachDBDialector 检查是否为 CockroachDB 方言
func IsCockroachDBDialector(dialector string) bool {
	return dialector == DialectorCockroachDB
}

// IsPostgreSQLFamilyDialector 检查是否为 PostgreSQL 家族方言（PostgreSQL 或 CockroachDB）
func IsPostgreSQLFamilyDialector(dialector string) bool {
	return IsPostgreSQLDialector(dialector) || IsCockroachDBDialector(dialector)
}

// IsSupportedDialector 检查是否为支持的数据库方言
func IsSupportedDialector(dialector string) bool {
	for _, supported := range SupportedDialectors {
		if dialector == supported {
			return true
		}
	}
	return false
}