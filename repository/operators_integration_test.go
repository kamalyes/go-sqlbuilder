/*
 * @Author: kamalyes 501893067@qq.com
 * @Date: 2025-11-26 00:00:00
 * @LastEditors: kamalyes 501893067@qq.com
 * @LastEditTime: 2025-11-26 00:00:00
 * @FilePath: \go-sqlbuilder\repository\operators_integration_test.go
 * @Description: 测试新增的操作符常量集成
 *
 * Copyright (c) 2025 by kamalyes, All Rights Reserved.
 */
package repository

import (
	"testing"

	"github.com/kamalyes/go-sqlbuilder/constants"
	"github.com/stretchr/testify/assert"
)

// TestNewOperatorsIntegration 测试新增操作符的集成
func TestNewOperatorsIntegration(t *testing.T) {
	// 测试 NewStartsWithFilter
	filter := NewStartsWithFilter("name", "user")
	assert.Equal(t, "name", filter.Field)
	assert.Equal(t, constants.OP_LIKE, filter.Operator)
	assert.Equal(t, "user%", filter.Value)

	// 测试 NewEndsWithFilter
	filter = NewEndsWithFilter("email", "@example.com")
	assert.Equal(t, "email", filter.Field)
	assert.Equal(t, constants.OP_LIKE, filter.Operator)
	assert.Equal(t, "%@example.com", filter.Value)
}

// TestBuildFilterConditionWithNewOperators 测试构建过滤条件的新操作符
func TestBuildFilterConditionWithNewOperators(t *testing.T) {
	testCases := []struct {
		name              string
		filter            *Filter
		expectedCondition string
		expectedArg       interface{}
		shouldHaveArg     bool
	}{
		{
			name:              "STARTS_WITH filter",
			filter:            &Filter{Field: "name", Operator: constants.OP_STARTS_WITH, Value: "user"},
			expectedCondition: "name LIKE ?",
			expectedArg:       "user" + constants.SQL_WILDCARD_ANY,
			shouldHaveArg:     true,
		},
		{
			name:              "ENDS_WITH filter",
			filter:            &Filter{Field: "email", Operator: constants.OP_ENDS_WITH, Value: "@example.com"},
			expectedCondition: "email LIKE ?",
			expectedArg:       constants.SQL_WILDCARD_ANY + "@example.com",
			shouldHaveArg:     true,
		},
		{
			name:              "CONTAINS filter",
			filter:            &Filter{Field: "description", Operator: constants.OP_CONTAINS, Value: "keyword"},
			expectedCondition: "description LIKE ?",
			expectedArg:       constants.SQL_WILDCARD_ANY + "keyword" + constants.SQL_WILDCARD_ANY,
			shouldHaveArg:     true,
		},
		{
			name:              "STARTS_WITH with non-string value",
			filter:            &Filter{Field: "name", Operator: constants.OP_STARTS_WITH, Value: 123},
			expectedCondition: "",
			shouldHaveArg:     false,
		},
		{
			name:              "ENDS_WITH with empty string",
			filter:            &Filter{Field: "name", Operator: constants.OP_ENDS_WITH, Value: ""},
			expectedCondition: "name LIKE ?",
			expectedArg:       constants.SQL_WILDCARD_ANY,
			shouldHaveArg:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			condition, arg := buildFilterCondition(tc.filter)
			assert.Equal(t, tc.expectedCondition, condition)

			if tc.shouldHaveArg {
				assert.Equal(t, tc.expectedArg, arg)
			} else {
				assert.Nil(t, arg)
			}
		})
	}
}

// TestWildcardConstants 测试通配符常量
func TestWildcardConstants(t *testing.T) {
	assert.Equal(t, "%", constants.SQL_WILDCARD_ANY)
	assert.Equal(t, "_", constants.SQL_WILDCARD_SINGLE)
}

// TestOperatorConstantsConsistency 测试操作符常量一致性
func TestOperatorConstantsConsistency(t *testing.T) {
	// 确保新的操作符常量都有对应的别名
	assert.Equal(t, constants.OpStartsWith, constants.OP_STARTS_WITH)
	assert.Equal(t, constants.OpEndsWith, constants.OP_ENDS_WITH)
	assert.Equal(t, constants.OpContains, constants.OP_CONTAINS)

	// 确保常量值符合预期
	assert.Equal(t, "STARTS_WITH", string(constants.OP_STARTS_WITH))
	assert.Equal(t, "ENDS_WITH", string(constants.OP_ENDS_WITH))
	assert.Equal(t, "CONTAINS", string(constants.OP_CONTAINS))
}
