/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-13 11:05:16
 * @FilePath: \go-sqlbuilder\logger\interface.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package logger

import "time"

type Logger interface {
	Debug(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Fatal(msg string, fields ...interface{})

	// Printf风格的方法
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})

	With(key string, value interface{}) Logger
	WithContext(ctx interface{}) Logger

	SetLevel(level string) Logger
	GetLevel() string

	StartTimer() Timer
	LogDuration(operation string, duration interface{}, fields ...interface{})
}

type Timer interface {
	Stop()
	Duration() time.Duration
}

type LoggerFactory interface {
	CreateLogger(name string) Logger
	SetDefaultLogger(logger Logger)
	GetDefaultLogger() Logger
}

const (
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
	LevelFatal = "fatal"
)
