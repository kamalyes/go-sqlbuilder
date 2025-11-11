/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-11 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-11 15:31:12
 * @FilePath: \go-sqlbuilder\errors\error_test.go
 * @Description: 错误处理测试
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */

package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewError(t *testing.T) {
	err := NewError(ErrCodeBuilderNotInitialized, "")
	assert.NotNil(t, err, "NewError should not return nil")
	assert.Equal(t, ErrCodeBuilderNotInitialized, err.Code, "Expected error code does not match")
	assert.Equal(t, "SQL builder not initialized", err.Message, "Expected error message does not match")
}

func TestNewErrorf(t *testing.T) {
	err := NewErrorf(ErrCodeInvalidTableName, "table %s not found", "users")
	assert.NotNil(t, err, "NewErrorf should not return nil")
	assert.Equal(t, ErrCodeInvalidTableName, err.Code, "Expected error code does not match")
	expectedDetails := "table users not found"
	assert.Equal(t, expectedDetails, err.Details, "Expected error details do not match")
}

func TestAppErrorError(t *testing.T) {
	err := NewError(ErrCodeInvalidFieldName, "")
	errStr := err.Error()
	assert.Equal(t, "[1003] Invalid field name", errStr, "Expected error string does not match")

	errWithDetails := NewError(ErrCodeInvalidFieldName, "field 'id' not found")
	errStr = errWithDetails.Error()
	assert.Equal(t, "[1003] Invalid field name: field 'id' not found", errStr, "Expected error string with details does not match")
}

func TestAppErrorString(t *testing.T) {
	err := NewError(ErrCodeCacheStoreNotConfigured, "")
	errStr := err.String()
	assert.Equal(t, "[2004] Cache store not configured", errStr, "Expected error string does not match")
}

func TestGetCode(t *testing.T) {
	err := NewError(ErrCodeInvalidOperator, "")
	code := err.GetCode()
	assert.Equal(t, ErrCodeInvalidOperator, code, "Expected error code does not match")
}

func TestGetMessage(t *testing.T) {
	err := NewError(ErrCodePageNumberInvalid, "")
	msg := err.GetMessage()
	assert.Equal(t, "Invalid page number", msg, "Expected error message does not match")
}

func TestGetDetails(t *testing.T) {
	expectedDetails := "page number must be greater than 0"
	err := NewError(ErrCodePageNumberInvalid, expectedDetails)
	details := err.GetDetails()
	assert.Equal(t, expectedDetails, details, "Expected error details do not match")
}

func TestWithDetails(t *testing.T) {
	err := NewError(ErrCodePageSizeInvalid, "")
	assert.Empty(t, err.Details, "Expected empty details initially")

	newErr := err.WithDetails("page size must be positive")
	assert.Equal(t, "page size must be positive", newErr.Details, "Expected details do not match")
	assert.Same(t, newErr, err, "WithDetails should return the same error object")
}

func TestIsErrorCode(t *testing.T) {
	err := NewError(ErrCodeRedisOperationFailed, "")
	assert.True(t, IsErrorCode(err, ErrCodeRedisOperationFailed), "IsErrorCode should return true for matching code")
	assert.False(t, IsErrorCode(err, ErrCodeRedisConnFailed), "IsErrorCode should return false for non-matching code")

	// Test with non-AppError
	normalErr := fmt.Errorf("some error")
	assert.False(t, IsErrorCode(normalErr, ErrCodeUnknown), "IsErrorCode should return false for non-AppError")
}

func TestGetErrorCode(t *testing.T) {
	err := NewError(ErrCodeTimeRangeInvalid, "")
	code := GetErrorCode(err)
	assert.Equal(t, ErrCodeTimeRangeInvalid, code, "Expected error code does not match")

	// Test with non-AppError
	normalErr := fmt.Errorf("some error")
	code = GetErrorCode(normalErr)
	assert.Equal(t, ErrCodeUnknown, code, "Expected code for non-AppError does not match")
}

func TestErrorCodeString(t *testing.T) {
	tests := []struct {
		code     ErrorCode
		expected string
	}{
		{ErrCodeBuilderNotInitialized, "SQL builder not initialized"},
		{ErrCodeCacheStoreNotFound, "Cache store not found"},
		{ErrCodeInvalidOperator, "Invalid query operator"},
		{ErrCodeRedisConnFailed, "Redis connection failed"},
		{ErrCodeUnknown, "Unknown error"},
	}

	for _, test := range tests {
		result := ErrorCodeString(test.code)
		assert.Equal(t, test.expected, result, "ErrorCodeString does not return the expected result")
	}
}

func TestPredefinedErrors(t *testing.T) {
	tests := []struct {
		err     *AppError
		code    ErrorCode
		message string
		name    string
	}{
		{ErrCacheNotConfigured, ErrCodeCacheStoreNotConfigured, "Cache store not configured", "ErrCacheNotConfigured"},
		{ErrCacheKeyMissing, ErrCodeCacheKeyNotFound, "Cache key not found", "ErrCacheKeyMissing"},
		{ErrCacheDataInvalid, ErrCodeCacheInvalidData, "Invalid cache data format", "ErrCacheDataInvalid"},
		{ErrInvalidOp, ErrCodeInvalidOperator, "Invalid query operator", "ErrInvalidOp"},
		{ErrEmptyFilter, ErrCodeEmptyFilterParam, "Empty filter parameter", "ErrEmptyFilter"},
		{ErrInvalidPage, ErrCodePageNumberInvalid, "Invalid page number", "ErrInvalidPage"},
		{ErrInvalidPageSize, ErrCodePageSizeInvalid, "Invalid page size", "ErrInvalidPageSize"},
		{ErrInvalidTimeRange, ErrCodeTimeRangeInvalid, "Invalid time range", "ErrInvalidTimeRange"},
		{ErrBuilderNotInit, ErrCodeBuilderNotInitialized, "SQL builder not initialized", "ErrBuilderNotInit"},
		{ErrInvalidTable, ErrCodeInvalidTableName, "Invalid table name", "ErrInvalidTable"},
		{ErrInvalidField, ErrCodeInvalidFieldName, "Invalid field name", "ErrInvalidField"},
		{ErrRedisNotImpl, ErrCodeRedisAdapterNotImpl, "Redis adapter not implemented", "ErrRedisNotImpl"},
		{ErrRedisFailed, ErrCodeRedisOperationFailed, "Redis operation failed", "ErrRedisFailed"},
	}

	for _, test := range tests {
		assert.NotNil(t, test.err, "Predefined error %s is nil", test.name)
		assert.Equal(t, test.code, test.err.Code, "%s: expected code %d, got %d", test.name, test.code, test.err.Code)
		assert.Equal(t, test.message, test.err.Message, "%s: expected message '%s', got '%s'", test.name, test.message, test.err.Message)
	}
}

func TestErrorInterface(t *testing.T) {
	err := NewError(ErrCodeInvalidSQLQuery, "")
	var _ error = err // 确保实现了 error 接口
	assert.Equal(t, "[1004] Invalid SQL query", err.Error(), "Error interface not properly implemented")
}

func TestStringerInterface(t *testing.T) {
	err := NewError(ErrCodeAdapterNotSupported, "")
	var _ fmt.Stringer = err // 确保实现了 fmt.Stringer 接口
	result := fmt.Sprintf("%s", err)
	assert.Equal(t, "[1005] Adapter not supported", result, "Stringer interface not properly implemented")
}
