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

import "sync"

type DefaultLoggerFactory struct {
	mu            sync.RWMutex
	defaultLogger Logger
	loggers       map[string]Logger
	loggerCreator LoggerCreator
}

type LoggerCreator func(name string) Logger

func NewLoggerFactory(creator LoggerCreator) *DefaultLoggerFactory {
	if creator == nil {
		creator = NewGoLogger
	}

	return &DefaultLoggerFactory{
		loggers:       make(map[string]Logger),
		loggerCreator: creator,
	}
}

func (f *DefaultLoggerFactory) CreateLogger(name string) Logger {
	f.mu.Lock()
	defer f.mu.Unlock()

	if logger, exists := f.loggers[name]; exists {
		return logger
	}

	logger := f.loggerCreator(name)
	f.loggers[name] = logger
	return logger
}

func (f *DefaultLoggerFactory) SetDefaultLogger(logger Logger) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.defaultLogger = logger
}

func (f *DefaultLoggerFactory) GetDefaultLogger() Logger {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if f.defaultLogger == nil {
		return NewNoOpLogger()
	}
	return f.defaultLogger
}
