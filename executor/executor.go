/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-13 11:05:04
 * @FilePath: \go-sqlbuilder\executor\executor.go
 * @Description:
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package executor

import (
	"context"
	"database/sql"
	"time"

	"github.com/kamalyes/go-sqlbuilder/core"
	"github.com/kamalyes/go-sqlbuilder/logger"
	"gorm.io/gorm"
)

type Executor interface {
	Query(ctx context.Context, sql string, args ...interface{}) (*sql.Rows, error)
	Exec(ctx context.Context, sql string, args ...interface{}) (sql.Result, error)
	WithDB(db *gorm.DB) Executor
	WithLogger(logger logger.Logger) Executor
	RegisterHook(hookType string, hook Hook) Executor
	Execute(ctx context.Context, condition *core.QueryCondition) (*QueryResult, error)
}

type executor struct {
	db          *gorm.DB
	logger      logger.Logger
	hooks       *HookRegistry
	metrics     *QueryMetrics
	retryPolicy *RetryPolicy
}

type QueryResult struct {
	Rows          *sql.Rows
	RowsAffected  int64
	LastInsertID  int64
	Duration      time.Duration
	Error         error
	CacheHit      bool
	ExecutedCount int
}

type RetryPolicy struct {
	MaxRetries int
	Delay      time.Duration
	BackOff    bool
}

func NewExecutor(db *gorm.DB) Executor {
	return &executor{
		db:      db,
		logger:  logger.NewNoOpLogger(),
		hooks:   NewHookRegistry(),
		metrics: NewQueryMetrics(),
		retryPolicy: &RetryPolicy{
			MaxRetries: 3,
			Delay:      100 * time.Millisecond,
			BackOff:    true,
		},
	}
}

func (e *executor) WithDB(db *gorm.DB) Executor {
	e.db = db
	return e
}

func (e *executor) WithLogger(log logger.Logger) Executor {
	if log != nil {
		e.logger = log
	}
	return e
}

func (e *executor) RegisterHook(hookType string, hook Hook) Executor {
	if e.hooks != nil {
		e.hooks.Register(hookType, hook)
	}
	return e
}

func (e *executor) Query(ctx context.Context, sql string, args ...interface{}) (*sql.Rows, error) {
	e.logger.Info("Executing query: %s", sql)

	execCtx := &ExecutionContext{
		Context:   ctx,
		SQL:       sql,
		Args:      args,
		StartTime: time.Now(),
		Metadata:  make(map[string]interface{}),
	}

	// Before execution hook
	if e.hooks != nil {
		if err := e.hooks.Execute(HookTypeBeforeExecution, execCtx); err != nil {
			e.logger.Error("Hook error: %v", err)
			return nil, err
		}
	}

	rows, err := e.db.WithContext(ctx).Raw(sql, args...).Rows()
	execCtx.EndTime = time.Now()
	execCtx.Duration = execCtx.EndTime.Sub(execCtx.StartTime)
	execCtx.Error = err

	if err != nil {
		e.logger.Error("Query failed: %v (duration: %v)", err, execCtx.Duration)

		// Error hook
		if e.hooks != nil {
			e.hooks.Execute(HookTypeOnError, execCtx)
		}
	} else {
		e.logger.Debug("Query succeeded (duration: %v)", execCtx.Duration)

		// After execution hook
		if e.hooks != nil {
			e.hooks.Execute(HookTypeAfterExecution, execCtx)
		}
	}

	// Update metrics
	if e.metrics != nil {
		e.metrics.RecordQuery(execCtx.Duration)
	}

	return rows, err
}

func (e *executor) Exec(ctx context.Context, sql string, args ...interface{}) (sql.Result, error) {
	e.logger.Info("Executing statement: %s", sql)

	execCtx := &ExecutionContext{
		Context:   ctx,
		SQL:       sql,
		Args:      args,
		StartTime: time.Now(),
		Metadata:  make(map[string]interface{}),
	}

	// Before execution hook
	if e.hooks != nil {
		if err := e.hooks.Execute(HookTypeBeforeExecution, execCtx); err != nil {
			e.logger.Error("Hook error: %v", err)
			return nil, err
		}
	}

	result := e.db.WithContext(ctx).Exec(sql, args...)
	execCtx.EndTime = time.Now()
	execCtx.Duration = execCtx.EndTime.Sub(execCtx.StartTime)
	execCtx.Error = result.Error

	if result.Error != nil {
		e.logger.Error("Statement failed: %v (duration: %v)", result.Error, execCtx.Duration)

		// Error hook
		if e.hooks != nil {
			e.hooks.Execute(HookTypeOnError, execCtx)
		}
		return nil, result.Error
	}

	e.logger.Debug("Statement succeeded (duration: %v, rows affected: %d)", execCtx.Duration, result.RowsAffected)

	// After execution hook
	if e.hooks != nil {
		e.hooks.Execute(HookTypeAfterExecution, execCtx)
	}

	// Update metrics
	if e.metrics != nil {
		e.metrics.RecordQuery(execCtx.Duration)
	}

	return &GormResult{
		rowsAffected: result.RowsAffected,
		lastInsertId: 0,
	}, nil
}

func (e *executor) Execute(ctx context.Context, condition *core.QueryCondition) (*QueryResult, error) {
	if condition == nil {
		return nil, nil
	}

	startTime := time.Now()

	// Apply query condition
	query := e.db.WithContext(ctx)

	queryApplier := core.NewQueryApplier()
	var err error
	query, err = queryApplier.ApplyCondition(query, condition)
	if err != nil {
		e.logger.Error("Failed to apply query condition: %v", err)
		return nil, err
	}

	// Execute query
	rows, err := query.Rows()
	if err != nil {
		e.logger.Error("Query execution failed: %v", err)
		return nil, err
	}

	duration := time.Since(startTime)

	result := &QueryResult{
		Rows:          rows,
		Duration:      duration,
		Error:         err,
		ExecutedCount: 1,
	}

	if e.metrics != nil {
		e.metrics.RecordQuery(duration)
	}

	return result, nil
}

type Hook func(ctx *ExecutionContext) error

type ExecutionContext struct {
	Context    context.Context
	SQL        string
	Args       []interface{}
	StartTime  time.Time
	EndTime    time.Time
	Duration   time.Duration
	Error      error
	Metadata   map[string]interface{}
	IsRetry    bool
	RetryCount int
}

const (
	HookTypeBeforeExecution = "before_execution"
	HookTypeAfterExecution  = "after_execution"
	HookTypeOnError         = "on_error"
)

type HookRegistry struct {
	hooks map[string][]Hook
}

func NewHookRegistry() *HookRegistry {
	return &HookRegistry{
		hooks: make(map[string][]Hook),
	}
}

func (hr *HookRegistry) Register(hookType string, hook Hook) {
	if hr.hooks[hookType] == nil {
		hr.hooks[hookType] = make([]Hook, 0)
	}
	hr.hooks[hookType] = append(hr.hooks[hookType], hook)
}

func (hr *HookRegistry) Execute(hookType string, ctx *ExecutionContext) error {
	hooks, ok := hr.hooks[hookType]
	if !ok {
		return nil
	}

	for _, hook := range hooks {
		if err := hook(ctx); err != nil {
			return err
		}
	}

	return nil
}

type QueryMetrics struct {
	TotalQueries    int64
	TotalDuration   time.Duration
	FastestQuery    time.Duration
	SlowestQuery    time.Duration
	AverageDuration time.Duration
}

func NewQueryMetrics() *QueryMetrics {
	return &QueryMetrics{
		FastestQuery: time.Hour,
		SlowestQuery: 0,
	}
}

func (qm *QueryMetrics) RecordQuery(duration time.Duration) {
	if qm == nil {
		return
	}

	qm.TotalQueries++
	qm.TotalDuration += duration

	if duration < qm.FastestQuery {
		qm.FastestQuery = duration
	}

	if duration > qm.SlowestQuery {
		qm.SlowestQuery = duration
	}

	qm.AverageDuration = qm.TotalDuration / time.Duration(qm.TotalQueries)
}

func (qm *QueryMetrics) Reset() {
	if qm == nil {
		return
	}

	qm.TotalQueries = 0
	qm.TotalDuration = 0
	qm.FastestQuery = time.Hour
	qm.SlowestQuery = 0
	qm.AverageDuration = 0
}

type GormResult struct {
	rowsAffected int64
	lastInsertId int64
}

func (gr *GormResult) LastInsertId() (int64, error) {
	return gr.lastInsertId, nil
}

func (gr *GormResult) RowsAffected() (int64, error) {
	return gr.rowsAffected, nil
}
