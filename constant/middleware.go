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
package constant

const (
	MiddlewareTypeLogging        = "logging"
	MiddlewareTypeMetrics        = "metrics"
	MiddlewareTypeRetry          = "retry"
	MiddlewareTypeTimeout        = "timeout"
	MiddlewareTypeValidation     = "validation"
	MiddlewareTypeCircuitBreaker = "circuit_breaker"
	MiddlewareTypeRateLimit      = "rate_limit"
)

const (
	HookTypeBeforeExecution = "before_execution"
	HookTypeAfterExecution  = "after_execution"
	HookTypeOnError         = "on_error"
)

const (
	DefaultRetryCount              = 3
	DefaultRetryDelay              = 100
	DefaultCircuitBreakerThreshold = 5
	DefaultRateLimit               = 1000
)
