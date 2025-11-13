/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-13 11:03:33
 * @FilePath: \go-sqlbuilder\constant\config.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package constant

// Default configuration constants
const (
	DefaultPageSize      = 10
	DefaultTimeout       = 30 // seconds
	DefaultMaxPoolSize   = 100
	DefaultMaxIdleConns  = 10
	DefaultCacheTTL      = 3600 // seconds
	DefaultSlowQueryTime = 100  // milliseconds
)

// Configuration keys
const (
	ConfigKeyShowSQL      = "show_sql"
	ConfigKeyLogLevel     = "log_level"
	ConfigKeyLogFile      = "log_file"
	ConfigKeyMaxRetries   = "max_retries"
	ConfigKeyCacheEnabled = "cache_enabled"
)
