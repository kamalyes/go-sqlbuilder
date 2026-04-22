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

// formatTimeGroup 获取格式化字符串
func formatTimeGroup(groupType TimeGroupType, formatMap map[TimeGroupType]string) string {
	format, ok := formatMap[groupType]
	if !ok {
		format = formatMap[GroupByDay] // 默认使用 GroupByDay 格式
	}
	return format
}

// MySQLDialect MySQL 方言
type MySQLDialect struct{}

func (d *MySQLDialect) FormatTimeGroup(field string, groupType TimeGroupType) string {
	formatMap := map[TimeGroupType]string{
		GroupByHour:  "%Y-%m-%d %H:00:00",
		GroupByDay:   "%Y-%m-%d",
		GroupByWeek:  "%Y-%u",
		GroupByMonth: "%Y-%m",
		GroupByYear:  "%Y",
	}
	format := formatTimeGroup(groupType, formatMap)
	return fmt.Sprintf("DATE_FORMAT(%s, '%s')", field, format)
}

// SQLiteDialect SQLite 方言
type SQLiteDialect struct{}

func (d *SQLiteDialect) FormatTimeGroup(field string, groupType TimeGroupType) string {
	formatMap := map[TimeGroupType]string{
		GroupByHour:  "%Y-%m-%d %H:00:00",
		GroupByDay:   "%Y-%m-%d",
		GroupByWeek:  "%Y-%W",
		GroupByMonth: "%Y-%m",
		GroupByYear:  "%Y",
	}
	format := formatTimeGroup(groupType, formatMap)
	return fmt.Sprintf("strftime('%s', %s)", format, field)
}

// PostgreSQLDialect PostgreSQL 方言
type PostgreSQLDialect struct{}

func (d *PostgreSQLDialect) FormatTimeGroup(field string, groupType TimeGroupType) string {
	formatMap := map[TimeGroupType]string{
		GroupByHour:  "YYYY-MM-DD HH24:00:00",
		GroupByDay:   "YYYY-MM-DD",
		GroupByWeek:  "IYYY-IW",
		GroupByMonth: "YYYY-MM",
		GroupByYear:  "YYYY",
	}
	format := formatTimeGroup(groupType, formatMap)
	return fmt.Sprintf("TO_CHAR(%s, '%s')", field, format)
}

// CockroachDBDialect CockroachDB 方言（兼容PostgreSQL语法）
type CockroachDBDialect struct{}

func (d *CockroachDBDialect) FormatTimeGroup(field string, groupType TimeGroupType) string {
	formatMap := map[TimeGroupType]string{
		GroupByHour:  "YYYY-MM-DD HH24:00:00",
		GroupByDay:   "YYYY-MM-DD",
		GroupByWeek:  "IYYY-IW",
		GroupByMonth: "YYYY-MM",
		GroupByYear:  "YYYY",
	}
	format := formatTimeGroup(groupType, formatMap)
	return fmt.Sprintf("TO_CHAR(%s, '%s')", field, format)
}

// ClickHouseDialect ClickHouse 方言
type ClickHouseDialect struct{}

func (d *ClickHouseDialect) FormatTimeGroup(field string, groupType TimeGroupType) string {
	formatMap := map[TimeGroupType]string{
		GroupByHour:  "%Y-%m-%d %H:00:00",
		GroupByDay:   "%Y-%m-%d",
		GroupByWeek:  "%Y-%U",
		GroupByMonth: "%Y-%m",
		GroupByYear:  "%Y",
	}
	format := formatTimeGroup(groupType, formatMap)
	return fmt.Sprintf("formatDateTime(%s, '%s')", field, format)
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
	case "cockroachdb", "cockroach":
		return &CockroachDBDialect{}
	case "clickhouse":
		return &ClickHouseDialect{}
	default:
		return &MySQLDialect{}
	}
}
