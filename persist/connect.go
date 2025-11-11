/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 21:57:26
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-11 21:57:35
 * @FilePath: \go-sqlbuilder\persist\connect.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package persist

import (
	"fmt"
	"os"
	"time"

	jsoniter "github.com/json-iterator/go"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewDBHandler 根据配置创建数据库连接（MySQL/SQLite）
func NewDBHandler(config *DBConfig) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	opts := &gorm.Config{}
	if config.ShowSQL {
		opts.Logger = logger.Default.LogMode(logger.Info)
	}
	switch config.Type {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
			config.Username,
			config.Password,
			config.Host,
			config.Port,
			config.DBName,
			config.Args,
		)
		db, err = gorm.Open(mysql.Open(dsn), opts)
	case "sqlite":
		dsn := config.Path
		if config.Args != "" {
			dsn = fmt.Sprintf("%s?%s", config.Path, config.Args)
		}
		db, err = gorm.Open(sqlite.Open(dsn), opts)
	default:
		return nil, fmt.Errorf("unsupported db type: %s", config.Type)
	}
	if err != nil {
		return nil, err
	}
	sqlDB, _ := db.DB()
	if config.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	}
	if config.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	}
	sqlDB.SetConnMaxIdleTime(time.Minute)
	return db, nil
}

func LoadDBConfigFromFile(configPath string) (*DBConfig, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	config := &DBConfig{}
	if err := jsoniter.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}
