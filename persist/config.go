/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 21:40:15
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-11 21:45:41
 * @FilePath: \go-sqlbuilder\persist\config.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package persist

// DBConfig 统一数据库配置结构，支持 MySQL、SQLite
type DBConfig struct {
	Type         string // "mysql" 或 "sqlite"，可扩展 "pg"
	Host         string // MySQL/PG
	Port         int    // MySQL/PG
	DBName       string // MySQL/PG
	Args         string // 连接参数
	Username     string // MySQL/PG
	Password     string // MySQL/PG
	Path         string // SQLite 路径
	ShowSQL      bool   // 是否显示 SQL 日志
	MaxOpenConns int    // 最大连接数
	MaxIdleConns int    // 最大空闲连接数
}