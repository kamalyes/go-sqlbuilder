/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-13 11:03:15
 * @FilePath: \go-sqlbuilder\constant\error.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package db

import "gorm.io/gorm"

// Handler 数据库处理器接口
// 所有 Repository 都基于这个接口工作
// 用户需要提供 GORM DB 实例的包装
type Handler interface {
	// DB 返回 GORM 数据库实例
	DB() *gorm.DB
}

// GormHandler 标准 GORM 处理器实现
type GormHandler struct {
	db *gorm.DB
}

// NewGormHandler 创建 GORM 处理器
func NewGormHandler(db *gorm.DB) Handler {
	if db == nil {
		panic("gorm DB instance cannot be nil")
	}
	return &GormHandler{db: db}
}

// DB 返回底层 GORM 实例
func (h *GormHandler) DB() *gorm.DB {
	return h.db
}
