/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-13 11:05:59
 * @FilePath: \go-sqlbuilder\logger\noop.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package logger

import "time"

type NoOpLogger struct{}

func NewNoOpLogger() Logger {
	return &NoOpLogger{}
}

func (n *NoOpLogger) Debug(msg string, fields ...interface{})                                   {}
func (n *NoOpLogger) Info(msg string, fields ...interface{})                                    {}
func (n *NoOpLogger) Warn(msg string, fields ...interface{})                                    {}
func (n *NoOpLogger) Error(msg string, fields ...interface{})                                   {}
func (n *NoOpLogger) Fatal(msg string, fields ...interface{})                                   {}
func (n *NoOpLogger) Debugf(format string, args ...interface{})                                 {}
func (n *NoOpLogger) Infof(format string, args ...interface{})                                  {}
func (n *NoOpLogger) Warnf(format string, args ...interface{})                                  {}
func (n *NoOpLogger) Errorf(format string, args ...interface{})                                 {}
func (n *NoOpLogger) Fatalf(format string, args ...interface{})                                 {}
func (n *NoOpLogger) With(key string, value interface{}) Logger                                 { return n }
func (n *NoOpLogger) WithContext(ctx interface{}) Logger                                        { return n }
func (n *NoOpLogger) SetLevel(level string) Logger                                              { return n }
func (n *NoOpLogger) GetLevel() string                                                          { return LevelError }
func (n *NoOpLogger) StartTimer() Timer                                                         { return &noOpTimer{} }
func (n *NoOpLogger) LogDuration(operation string, duration interface{}, fields ...interface{}) {}

type noOpTimer struct{}

func (t *noOpTimer) Stop()                   {}
func (t *noOpTimer) Duration() time.Duration { return 0 }
