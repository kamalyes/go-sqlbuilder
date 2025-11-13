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
package core

import (
	"fmt"

	"github.com/kamalyes/go-sqlbuilder/constant"
	"github.com/kamalyes/go-sqlbuilder/errors"
)

type OperatorTranslator struct{}

func NewOperatorTranslator() *OperatorTranslator {
	return &OperatorTranslator{}
}

func (ot *OperatorTranslator) Translate(op interface{}) (constant.Operator, error) {
	switch v := op.(type) {
	case constant.Operator:
		return v, nil
	case string:
		return ot.translateString(v)
	default:
		return "", errors.NewError(
			errors.ErrorCodeInvalidOperator,
			fmt.Sprintf(constant.ErrUnknownOperator, fmt.Sprintf("%T", op)),
		)
	}
}

func (ot *OperatorTranslator) translateString(op string) (constant.Operator, error) {
	// 尝试标准操作符映射
	if stdOp, exists := constant.OperatorMap[op]; exists {
		return stdOp, nil
	}

	// 尝试兼容操作符映射
	if compatOp, exists := constant.CompatOperatorMap[op]; exists {
		return compatOp, nil
	}

	return "", errors.NewError(
		errors.ErrorCodeInvalidOperator,
		fmt.Sprintf(constant.ErrUnknownOperator, op),
	)
}

func (ot *OperatorTranslator) IsComparisonOperator(op constant.Operator) bool {
	switch op {
	case constant.OP_EQ, constant.OP_NEQ, constant.OP_GT, constant.OP_GTE, constant.OP_LT, constant.OP_LTE:
		return true
	}
	return false
}

func (ot *OperatorTranslator) IsStringOperator(op constant.Operator) bool {
	switch op {
	case constant.OP_LIKE, constant.OP_NOT_LIKE:
		return true
	}
	return false
}

func (ot *OperatorTranslator) IsSetOperator(op constant.Operator) bool {
	switch op {
	case constant.OP_IN, constant.OP_NOT_IN, constant.OP_BETWEEN:
		return true
	}
	return false
}

func (ot *OperatorTranslator) IsNullOperator(op constant.Operator) bool {
	switch op {
	case constant.OP_IS_NULL, constant.OP_IS_NOT_NULL:
		return true
	}
	return false
}

func (ot *OperatorTranslator) RequiresValue(op constant.Operator) bool {
	return !ot.IsNullOperator(op)
}

func (ot *OperatorTranslator) IsSpecialOperator(op constant.Operator) bool {
	return op == constant.OP_FIND_IN_SET
}

type FilterBuilder struct {
	filters []Filter
}

func NewFilterBuilder() *FilterBuilder {
	return &FilterBuilder{
		filters: make([]Filter, 0),
	}
}

func (fb *FilterBuilder) AddFilter(field string, operator constant.Operator, value interface{}) *FilterBuilder {
	fb.filters = append(fb.filters, Filter{
		Field:    field,
		Operator: operator,
		Value:    value,
	})
	return fb
}

func (fb *FilterBuilder) AddEQ(field string, value interface{}) *FilterBuilder {
	return fb.AddFilter(field, constant.OP_EQ, value)
}

func (fb *FilterBuilder) AddNEQ(field string, value interface{}) *FilterBuilder {
	return fb.AddFilter(field, constant.OP_NEQ, value)
}

func (fb *FilterBuilder) AddGT(field string, value interface{}) *FilterBuilder {
	return fb.AddFilter(field, constant.OP_GT, value)
}

func (fb *FilterBuilder) AddGTE(field string, value interface{}) *FilterBuilder {
	return fb.AddFilter(field, constant.OP_GTE, value)
}

func (fb *FilterBuilder) AddLT(field string, value interface{}) *FilterBuilder {
	return fb.AddFilter(field, constant.OP_LT, value)
}

func (fb *FilterBuilder) AddLTE(field string, value interface{}) *FilterBuilder {
	return fb.AddFilter(field, constant.OP_LTE, value)
}

func (fb *FilterBuilder) AddLIKE(field string, value interface{}) *FilterBuilder {
	return fb.AddFilter(field, constant.OP_LIKE, value)
}

func (fb *FilterBuilder) AddNOT_LIKE(field string, value interface{}) *FilterBuilder {
	return fb.AddFilter(field, constant.OP_NOT_LIKE, value)
}

func (fb *FilterBuilder) AddIN(field string, value interface{}) *FilterBuilder {
	return fb.AddFilter(field, constant.OP_IN, value)
}

func (fb *FilterBuilder) AddNOT_IN(field string, value interface{}) *FilterBuilder {
	return fb.AddFilter(field, constant.OP_NOT_IN, value)
}

func (fb *FilterBuilder) AddBETWEEN(field string, value interface{}) *FilterBuilder {
	return fb.AddFilter(field, constant.OP_BETWEEN, value)
}

func (fb *FilterBuilder) AddIS_NULL(field string) *FilterBuilder {
	return fb.AddFilter(field, constant.OP_IS_NULL, nil)
}

func (fb *FilterBuilder) AddIS_NOT_NULL(field string) *FilterBuilder {
	return fb.AddFilter(field, constant.OP_IS_NOT_NULL, nil)
}

func (fb *FilterBuilder) AddFIND_IN_SET(field string, value interface{}) *FilterBuilder {
	return fb.AddFilter(field, constant.OP_FIND_IN_SET, value)
}

func (fb *FilterBuilder) Build() []Filter {
	return fb.filters
}

func (fb *FilterBuilder) Clear() *FilterBuilder {
	fb.filters = make([]Filter, 0)
	return fb
}

type OrderBuilder struct {
	orders []OrderBy
}

func NewOrderBuilder() *OrderBuilder {
	return &OrderBuilder{
		orders: make([]OrderBy, 0),
	}
}

func (ob *OrderBuilder) AddOrder(field string, order string) *OrderBuilder {
	ob.orders = append(ob.orders, OrderBy{
		Field: field,
		Order: order,
	})
	return ob
}

func (ob *OrderBuilder) Asc(field string) *OrderBuilder {
	return ob.AddOrder(field, constant.OrderASC)
}

func (ob *OrderBuilder) Desc(field string) *OrderBuilder {
	return ob.AddOrder(field, constant.OrderDESC)
}

func (ob *OrderBuilder) Build() []OrderBy {
	return ob.orders
}

func (ob *OrderBuilder) Clear() *OrderBuilder {
	ob.orders = make([]OrderBy, 0)
	return ob
}

type PaginationBuilder struct {
	page     int
	pageSize int
}

func NewPaginationBuilder() *PaginationBuilder {
	return &PaginationBuilder{
		page:     1,
		pageSize: constant.DefaultPageSize,
	}
}

func (pb *PaginationBuilder) Page(page int) *PaginationBuilder {
	if page > 0 {
		pb.page = page
	}
	return pb
}

func (pb *PaginationBuilder) PageSize(pageSize int) *PaginationBuilder {
	if pageSize > 0 {
		pb.pageSize = pageSize
	}
	return pb
}

func (pb *PaginationBuilder) Offset(offset int) *PaginationBuilder {
	if offset >= 0 {
		pb.page = offset/pb.pageSize + 1
	}
	return pb
}

func (pb *PaginationBuilder) Limit(limit int) *PaginationBuilder {
	if limit > 0 {
		pb.pageSize = limit
	}
	return pb
}

func (pb *PaginationBuilder) Build() *Pagination {
	offset := (pb.page - 1) * pb.pageSize
	return &Pagination{
		Page:     pb.page,
		PageSize: pb.pageSize,
		Offset:   offset,
		Limit:    pb.pageSize,
	}
}
