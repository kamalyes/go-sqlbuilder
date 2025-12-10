/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-12-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-12-11 00:00:00
 * @FilePath: \go-sqlbuilder\repository\dialect.go
 * @Description: 数据库方言适配器
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package repository

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// Dialect 数据库方言接口
type Dialect interface {
	// FormatTimeGroup 格式化时间分组表达式
	FormatTimeGroup(field string, groupType TimeGroupType) string
}

// MySQLDialect MySQL 方言
type MySQLDialect struct{}

func (d *MySQLDialect) FormatTimeGroup(field string, groupType TimeGroupType) string {
	var format string
	switch groupType {
	case GroupByHour:
		format = "%Y-%m-%d %H:00:00"
	case GroupByDay:
		format = "%Y-%m-%d"
	case GroupByWeek:
		format = "%Y-%u"
	case GroupByMonth:
		format = "%Y-%m"
	case GroupByYear:
		format = "%Y"
	default:
		format = "%Y-%m-%d"
	}
	return fmt.Sprintf("DATE_FORMAT(%s, '%s')", field, format)
}

// SQLiteDialect SQLite 方言
type SQLiteDialect struct{}

func (d *SQLiteDialect) FormatTimeGroup(field string, groupType TimeGroupType) string {
	var format string
	switch groupType {
	case GroupByHour:
		format = "%Y-%m-%d %H:00:00"
	case GroupByDay:
		format = "%Y-%m-%d"
	case GroupByWeek:
		format = "%Y-%W"
	case GroupByMonth:
		format = "%Y-%m"
	case GroupByYear:
		format = "%Y"
	default:
		format = "%Y-%m-%d"
	}
	return fmt.Sprintf("strftime('%s', %s)", format, field)
}

// PostgreSQLDialect PostgreSQL 方言
type PostgreSQLDialect struct{}

func (d *PostgreSQLDialect) FormatTimeGroup(field string, groupType TimeGroupType) string {
	var format string
	switch groupType {
	case GroupByHour:
		format = "YYYY-MM-DD HH24:00:00"
	case GroupByDay:
		format = "YYYY-MM-DD"
	case GroupByWeek:
		format = "IYYY-IW"
	case GroupByMonth:
		format = "YYYY-MM"
	case GroupByYear:
		format = "YYYY"
	default:
		format = "YYYY-MM-DD"
	}
	return fmt.Sprintf("TO_CHAR(%s, '%s')", field, format)
}

// DetectDialect 自动检测数据库方言
func DetectDialect(db *gorm.DB) Dialect {
	dialector := db.Dialector.Name()
	
	switch strings.ToLower(dialector) {
	case "mysql":
		return &MySQLDialect{}
	case "sqlite", "sqlite3":
		return &SQLiteDialect{}
	case "postgres", "postgresql":
		return &PostgreSQLDialect{}
	default:
		// 默认使用 MySQL 方言
		return &MySQLDialect{}
	}
}
