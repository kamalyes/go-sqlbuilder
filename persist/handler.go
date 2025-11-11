/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 21:14:54
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-11 21:21:09
 * @FilePath: \go-sqlbuilder\persist\handler.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package persist

import (
	"context"
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type Option func(db *gorm.DB) *gorm.DB

type DBHandler interface {
	DB() *gorm.DB
	Query(param *QueryParam) *gorm.DB
	// auto migrate
	AutoMigrate(dst ...any) error
	// tx
	Begin(opts ...*sql.TxOptions) DBHandler
	Commit() error
	Rollback() error

	Close() error
    // 兼容 go-sqlbuilder，执行原生SQL
    ExecSQL(ctx context.Context, sql string, args ...interface{}) (*gorm.DB, error)
    RawQuery(ctx context.Context, sql string, args ...interface{}) (*gorm.DB, error)
}

type CacheHandler interface {
	Set([]byte, []byte) error
	SetWithTTL([]byte, []byte, time.Duration) error
	Get([]byte) ([]byte, error)
	GetTTL([]byte) (time.Duration, error)
	Del(...[]byte) error

	Close() error
}

