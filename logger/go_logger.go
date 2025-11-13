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
package logger

import (
	"fmt"
	"time"

	log "github.com/kamalyes/go-logger"
)

type GoLoggerAdapter struct {
	logger log.ILogger
	name   string
	level  string
}

func NewGoLogger(name string) Logger {
	return &GoLoggerAdapter{
		logger: log.GetGlobalLogger().WithField("component", name),
		name:   name,
		level:  LevelInfo,
	}
}

func (gl *GoLoggerAdapter) Debug(msg string, fields ...interface{}) {
	if gl.isLevelEnabled(LevelDebug) {
		gl.logger.Debug(fmt.Sprintf(msg, fields...))
	}
}

func (gl *GoLoggerAdapter) Info(msg string, fields ...interface{}) {
	if gl.isLevelEnabled(LevelInfo) {
		gl.logger.Info(fmt.Sprintf(msg, fields...))
	}
}

func (gl *GoLoggerAdapter) Warn(msg string, fields ...interface{}) {
	if gl.isLevelEnabled(LevelWarn) {
		gl.logger.Warn(fmt.Sprintf(msg, fields...))
	}
}

func (gl *GoLoggerAdapter) Error(msg string, fields ...interface{}) {
	if gl.isLevelEnabled(LevelError) {
		gl.logger.Error(fmt.Sprintf(msg, fields...))
	}
}

func (gl *GoLoggerAdapter) Fatal(msg string, fields ...interface{}) {
	gl.logger.Fatal(fmt.Sprintf(msg, fields...))
}

func (gl *GoLoggerAdapter) Debugf(format string, args ...interface{}) {
	if gl.isLevelEnabled(LevelDebug) {
		gl.logger.Debug(fmt.Sprintf(format, args...))
	}
}

func (gl *GoLoggerAdapter) Infof(format string, args ...interface{}) {
	if gl.isLevelEnabled(LevelInfo) {
		gl.logger.Info(fmt.Sprintf(format, args...))
	}
}

func (gl *GoLoggerAdapter) Warnf(format string, args ...interface{}) {
	if gl.isLevelEnabled(LevelWarn) {
		gl.logger.Warn(fmt.Sprintf(format, args...))
	}
}

func (gl *GoLoggerAdapter) Errorf(format string, args ...interface{}) {
	if gl.isLevelEnabled(LevelError) {
		gl.logger.Error(fmt.Sprintf(format, args...))
	}
}

func (gl *GoLoggerAdapter) Fatalf(format string, args ...interface{}) {
	gl.logger.Fatal(fmt.Sprintf(format, args...))
}

func (gl *GoLoggerAdapter) With(key string, value interface{}) Logger {
	newLogger := gl.logger.WithField(key, value)
	return &GoLoggerAdapter{
		logger: newLogger,
		name:   gl.name,
		level:  gl.level,
	}
}

func (gl *GoLoggerAdapter) WithContext(ctx interface{}) Logger {
	return gl.With("context", ctx)
}

func (gl *GoLoggerAdapter) SetLevel(level string) Logger {
	gl.level = level
	return gl
}

func (gl *GoLoggerAdapter) GetLevel() string {
	return gl.level
}

func (gl *GoLoggerAdapter) StartTimer() Timer {
	return &timer{
		start:  time.Now(),
		logger: gl,
	}
}

func (gl *GoLoggerAdapter) LogDuration(operation string, duration interface{}, fields ...interface{}) {
	msg := fmt.Sprintf("Operation %s completed in %v", operation, duration)
	gl.Info(msg, fields...)
}

func (gl *GoLoggerAdapter) isLevelEnabled(level string) bool {
	levels := map[string]int{
		LevelDebug: 0,
		LevelInfo:  1,
		LevelWarn:  2,
		LevelError: 3,
		LevelFatal: 4,
	}

	currentLevel := levels[gl.level]
	targetLevel := levels[level]

	return targetLevel >= currentLevel
}

type timer struct {
	start  time.Time
	logger Logger
}

func (t *timer) Stop() {
	duration := time.Since(t.start)
	t.logger.LogDuration("operation", duration)
}

func (t *timer) Duration() time.Duration {
	return time.Since(t.start)
}
