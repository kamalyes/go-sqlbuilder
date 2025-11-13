/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-13 11:03:27
 * @FilePath: \go-sqlbuilder\constant\adapter.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package constant

const (
	AdapterTypeGORM     = "gorm"
	AdapterTypeSQLX     = "sqlx"
	AdapterTypeDatabase = "database"
	AdapterTypeCustom   = "custom"
)

const (
	DialectMySQL     = "mysql"
	DialectPostgres  = "postgres"
	DialectSQLite    = "sqlite"
	DialectSQLServer = "sqlserver"
)

const (
	ParameterPlaceholder = "?"
	ParameterStyle       = "question"
)
